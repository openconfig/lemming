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

package controllers

import (
	"context"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/google/go-cmp/cmp"
	lemmingscheme "github.com/openconfig/lemming/operator/api/clientset/scheme"
	lemmingv1alpha1 "github.com/openconfig/lemming/operator/api/lemming/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
)

func TestReconcile(t *testing.T) {
	lemmingscheme.AddToScheme(k8sscheme.Scheme)

	tests := []struct {
		desc        string
		init        []client.Object
		wantErr     bool
		wantRes     reconcile.Result
		wantService *corev1.Service
		wantPod     *corev1.Pod
	}{{
		desc: "create on empty cluster",
		init: []client.Object{&lemmingv1alpha1.Lemming{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lemming",
				Namespace: "fake",
			},
			Spec: lemmingv1alpha1.LemmingSpec{
				Image:   "lemming:latest",
				Command: "./test",
				Args:    []string{"--alsologtostderr"},
				Env: []corev1.EnvVar{{
					Name:  "ENV_TEST",
					Value: "FOO",
				}},
				ConfigPath: "/config",
				ConfigFile: "foo.yaml",
				InitImage:  "sleep:latest",
				Ports: map[string]lemmingv1alpha1.ServicePort{
					"gnmi": {
						InnerPort: 1234,
						OuterPort: 5678,
					},
				},
				TLS: lemmingv1alpha1.TLSSpec{
					SelfSigned: lemmingv1alpha1.SelfSignedSpec{
						KeySize:    2048,
						CommonName: "lemming",
					},
				},
				InterfaceCount: 1,
				InitSleep:      1,
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU: resource.MustParse("1"),
					},
				},
			},
		}},
		wantService: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-lemming",
				Namespace: "fake",
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app":  "lemming",
					"topo": "fake",
				},
				Type: corev1.ServiceTypeLoadBalancer,
				Ports: []corev1.ServicePort{{
					Name:       "gnmi",
					Protocol:   corev1.ProtocolTCP,
					Port:       5678,
					TargetPort: intstr.FromInt(1234),
				}},
			},
		},
		wantPod: &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lemming",
				Namespace: "fake",
			},
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{{
					Name:  "init",
					Image: "sleep:latest",
					Args:  []string{"1", "1"},
				}},
				Containers: []corev1.Container{{
					Name:    "lemming",
					Image:   "lemming:latest",
					Command: []string{"./test"},
					Args:    []string{"--alsologtostderr", "--target=lemming"},
					Env: []corev1.EnvVar{{
						Name:  "ENV_TEST",
						Value: "FOO",
					}},
					SecurityContext: &corev1.SecurityContext{
						Privileged: pointer.Bool(true),
					},
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("1"),
						},
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "tls",
						ReadOnly:  true,
						MountPath: "/certs",
					}},
				}},
				Volumes: []corev1.Volume{{
					Name: "tls",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "lemming-tls",
						},
					},
				}},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			c := fake.NewClientBuilder().WithScheme(k8sscheme.Scheme).WithObjects(tt.init...).Build()
			lr := &LemmingReconciler{
				Client: c,
				Scheme: k8sscheme.Scheme,
			}
			res, err := lr.Reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Name: "lemming", Namespace: "fake"}})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Reconcile() unexpected error: got %v, want %t", err, tt.wantErr)
			}
			if !reflect.DeepEqual(res, tt.wantRes) {
				t.Fatalf("Reconcile() unexpected result got: %v, want %v", res, tt.wantRes)
			}
			svc := &corev1.Service{}
			if err := c.Get(context.Background(), types.NamespacedName{Name: tt.wantService.Name, Namespace: tt.wantService.Namespace}, svc); err != nil {
				t.Fatalf("Failed to get svc: %v", err)
			}
			if diff := cmp.Diff(tt.wantService.Spec, svc.Spec); diff != "" {
				t.Errorf("Got service different from expected (-want,+got):\n %s", diff)
			}
			pod := &corev1.Pod{}
			if err := c.Get(context.Background(), types.NamespacedName{Name: tt.wantPod.Name, Namespace: tt.wantPod.Namespace}, pod); err != nil {
				t.Fatalf("Failed to get pod: %v", err)
			}
			if diff := cmp.Diff(tt.wantPod.Spec, pod.Spec); diff != "" {
				t.Errorf("Got pod different from expected (-want,+got):\n %s", diff)
			}
		})
	}
}
