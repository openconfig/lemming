// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package binding wraps knebind to allow running integration tests without the need to supply kne/ondatra flags for configuration files.
package binding

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/openconfig/ondatra/binding"
	"github.com/openconfig/ondatra/proto"
	"google.golang.org/protobuf/encoding/prototext"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	tpb "github.com/openconfig/kne/proto/topo"
	kinit "github.com/openconfig/ondatra/knebind/init"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const config = `
topology: {{  .Topology }}
kubecfg: {{ .Kubeconfig }}
cli: {{ .CLI }}
username: {{ .Username }}
password: {{ .Password }}
`

type configFields struct {
	Topology   string
	Kubeconfig string
	CLI        string
	Username   string
	Password   string
}

var (
	keep = flag.Bool("keep", false, "Keep topology deployed after test")
)

// Get returns the custom lemming binding. The topoDir is the relative path to a
// directory containing the Ondatra testbed and KNE topology pb.txt files.
func Get(topoDir string) func() (binding.Binding, error) {
	dir, _ := filepath.Abs(topoDir)
	testbedFile := filepath.Join(dir, "testbed.pb.txt")
	topoFile := filepath.Join(dir, "topology.pb.txt")

	flag.Set("testbed", testbedFile)

	return func() (binding.Binding, error) {
		fmt.Println("Checking for KNE CLI installation")
		knePath, err := getKNECLIPath()
		if err != nil {
			return nil, err
		}

		fmt.Println("Creating knebind config file")
		if err := writeKNEBindCfg(knePath, topoFile); err != nil {
			return nil, err
		}
		b, err := kinit.Init()
		if err != nil {
			return nil, err
		}
		return &LemmingBind{
			Binding:  b,
			topoFile: topoFile,
			kubecfg:  os.ExpandEnv("$HOME/.kube/config"),
		}, nil
	}
}

// LemmingBind wraps the Ondatra knebind and adds the ability to set the testbed and topology from inside the test.
// TODO: Add unit tests and upstream some of this to Ondatra.
type LemmingBind struct {
	binding.Binding
	topoFile string
	kubecfg  string
	created  bool
	id       string
	cancel   func()
}

// Release runs knebind release then deletes the topology if it was created by this binding.
func (lb *LemmingBind) Release(ctx context.Context) error {
	defer lb.cancel()
	if err := lb.Binding.Release(ctx); err != nil {
		return err
	}
	if !lb.created || *keep {
		return nil
	}
	if out, err := exec.Command("kne", "delete", lb.topoFile).CombinedOutput(); err != nil {
		return fmt.Errorf("failed delete topology: %v output:\n%s", err, string(out))
	}
	return nil
}

// Reserve deploys the topology if it doesn't exists, then runs knebind Reserve.
func (lb *LemmingBind) Reserve(ctx context.Context, tb *proto.Testbed, runTime time.Duration, waitTime time.Duration, partial map[string]string) (*binding.Reservation, error) {
	topoBytes, err := os.ReadFile(lb.topoFile)
	if err != nil {
		return nil, err
	}
	topo := &tpb.Topology{}
	if err := prototext.Unmarshal(topoBytes, topo); err != nil {
		return nil, err
	}

	// Get kubernetes API client.
	cfg, err := clientcmd.BuildConfigFromFlags("", lb.kubecfg)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Acquire a lock on this topology, prevent another test from using this topology while this test is running.
	lb.id = uuid.New().String()
	lock, err := resourcelock.New(
		resourcelock.ConfigMapsLeasesResourceLock,
		"default",
		topo.GetName(),
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{Identity: lb.id},
	)
	if err != nil {
		return nil, err
	}

	ready := make(chan bool)
	elect, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		ReleaseOnCancel: true,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				ready <- true
			},
			OnStoppedLeading: func() {},
			OnNewLeader:      func(identity string) {},
		},
	})
	if err != nil {
		return nil, err
	}
	electCtx, cancel := context.WithCancel(ctx)
	lb.cancel = cancel

	go elect.Run(electCtx)

	fmt.Println("Waiting for topology lock")
	<-ready

	// Check if topology already exists, if not deploy it.
	if _, err := client.CoreV1().Namespaces().Get(ctx, topo.GetName(), metav1.GetOptions{}); apierrors.IsNotFound(err) {
		fmt.Println("Deploying KNE topology")
		if out, err := exec.Command("kne", "create", lb.topoFile).CombinedOutput(); err != nil {
			return nil, fmt.Errorf("failed create topology: %v output:\n%s", err, string(out))
		}
		lb.created = true
		// TODO: Wait for all pods to be ready.
		time.Sleep(5 * time.Second)
	} else if err != nil {
		return nil, err
	}

	fmt.Println("Reserving KNE topology")
	return lb.Binding.Reserve(ctx, tb, runTime, waitTime, partial)
}

func getKNECLIPath() (string, error) {
	path, err := exec.LookPath("kne")
	if err != nil {
		return exec.LookPath("kne_cli")
	}
	return path, nil
}

func writeKNEBindCfg(knePath, topoFile string) error {
	cf := &configFields{
		Kubeconfig: os.ExpandEnv("$HOME/.kube/config"),
		CLI:        knePath,
		Username:   "fake",
		Password:   "fake",
		Topology:   topoFile,
	}
	tmpl, err := template.New("config").Parse(config)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.CreateTemp("", "config")
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, cf); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	flag.Set("config", f.Name())
	return nil
}
