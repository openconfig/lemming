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
	"fmt"
	"path/filepath"
	"sort"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/cert"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	lemmingv1alpha1 "github.com/openconfig/lemming/operator/api/lemming/v1alpha1"
)

// LemmingReconciler reconciles a Lemming object
type LemmingReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	ExtraArgs []string
}

//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods;services;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *LemmingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	lemming := &lemmingv1alpha1.Lemming{}
	if err := r.Get(ctx, req.NamespacedName, lemming); err != nil {
		log.Error(err, "unable to get lemming")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secret, err := r.reconcileSecrets(ctx, lemming)
	if err != nil {
		log.Error(err, "unable to get reconcile secret")
		return ctrl.Result{}, err
	}
	var secretName string
	if secret != nil {
		secretName = secret.GetName()
	}

	pod, err := r.reconcilePod(ctx, lemming, secretName)
	if err != nil {
		log.Error(err, "unable to get reconcile pod")
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, lemming); err != nil {
		log.Error(err, "unable to get reconcile service: %v")
		return ctrl.Result{}, err
	}

	switch pod.Status.Phase {
	case corev1.PodRunning:
		lemming.Status.Phase = lemmingv1alpha1.Running
	case corev1.PodFailed:
		lemming.Status.Phase = lemmingv1alpha1.Failed
	default:
		lemming.Status.Phase = lemmingv1alpha1.Unknown
	}
	lemming.Status.Message = fmt.Sprintf("Pod Details: %s", pod.Status.Message)
	if err := r.Status().Update(ctx, lemming); err != nil {
		log.Error(err, "unable to update lemming status")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *LemmingReconciler) reconcileSecrets(ctx context.Context, lemming *lemmingv1alpha1.Lemming) (*corev1.Secret, error) {
	log := log.FromContext(ctx)
	secretName := fmt.Sprintf("%s-tls", lemming.Name)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: lemming.Namespace,
		},
	}
	err := r.Get(ctx, client.ObjectKeyFromObject(secret), secret)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	if !apierrors.IsNotFound(err) {
		if lemming.Spec.TLS == nil {
			log.Info("no tls config and secret exists, deleting it.")
			return nil, r.Delete(ctx, secret)
		}
		return secret, nil
	}
	if lemming.Spec.TLS == nil {
		log.Info("no tls config and secret doesn't exist, doing nothing.")
		return nil, nil
	}

	if lemming.Spec.TLS.SelfSigned != nil {
		if err := ctrl.SetControllerReference(lemming, secret, r.Scheme); err != nil {
			return nil, err
		}
		cert, key, err := cert.GenerateSelfSignedCertKey(lemming.Spec.TLS.SelfSigned.CommonName, nil, nil)
		if err != nil {
			return nil, err
		}
		secret.Data = map[string][]byte{
			"tls.crt": cert,
			"tls.key": key,
		}
		secret.Type = corev1.SecretTypeTLS
		log.Info("tls config not empty and secret doesn't exist, creating it.")
		return secret, r.Create(ctx, secret)
	}
	log.Info("no self signed config and secret doesn't exist, doing nothing.")
	return nil, nil
}

const (
	secretMountPath = "/certs"
)

var requiredArgs = map[string]struct{}{
	"--enable_dataplane": {},
	"--alsologtostderr":  {},
}

func (r *LemmingReconciler) reconcilePod(ctx context.Context, lemming *lemmingv1alpha1.Lemming, secretName string) (*corev1.Pod, error) {
	log := log.FromContext(ctx)
	pod := &corev1.Pod{}
	err := r.Get(ctx, types.NamespacedName{Name: lemming.Name, Namespace: lemming.Namespace}, pod)
	var newPod bool

	if apierrors.IsNotFound(err) {
		log.Info("new pod, creating initial spec")
		if err := r.setupInitialPod(pod, lemming); err != nil {
			return nil, fmt.Errorf("failed to setup initial pod: %v", err)
		}
		newPod = true
	} else if err != nil {
		return nil, err
	}

	oldPodSpec := pod.Spec.DeepCopy()

	pod.Spec.Containers[0].Image = lemming.Spec.Image
	pod.Spec.InitContainers[0].Image = lemming.Spec.InitImage
	pod.Spec.InitContainers[0].Args = []string{fmt.Sprintf("%d", lemming.Spec.InterfaceCount), fmt.Sprintf("%d", lemming.Spec.InitSleep)}
	pod.Spec.Containers[0].Command = []string{lemming.Spec.Command}
	pod.Spec.Containers[0].Args = append(lemming.Spec.Args, fmt.Sprintf("--target=%s", lemming.Name))
	pod.Spec.Containers[0].Env = lemming.Spec.Env
	pod.Spec.Containers[0].Resources = lemming.Spec.Resources

	for _, arg := range pod.Spec.Containers[0].Args {
		if _, ok := requiredArgs[arg]; ok {
			delete(requiredArgs, arg)
		}
	}
	sortedArgs := make([]string, 0, len(requiredArgs))
	for arg := range requiredArgs {
		sortedArgs = append(sortedArgs, arg)
	}
	sortedArgs = append(sortedArgs, r.ExtraArgs...)
	sort.Strings(sortedArgs)
	pod.Spec.Containers[0].Args = append(pod.Spec.Containers[0].Args, sortedArgs...)

	mounts := map[string]corev1.VolumeMount{}
	volumes := map[string]corev1.Volume{}

	for _, vol := range pod.Spec.Volumes {
		volumes[vol.Name] = vol
	}
	for _, mount := range pod.Spec.Containers[0].VolumeMounts {
		mounts[mount.Name] = mount
	}

	var changedMounts bool
	if _, ok := volumes["tls"]; secretName != "" && !ok {
		log.Info("adding tls secret to pod spec")
		volumes["tls"] = corev1.Volume{
			Name: "tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		}
		mounts["tls"] = corev1.VolumeMount{
			Name:      "tls",
			ReadOnly:  true,
			MountPath: secretMountPath,
		}
		changedMounts = true
	} else if secretName == "" && ok {
		delete(mounts, "tls")
		delete(volumes, "tls")
		changedMounts = true
	}
	if secretName != "" {
		pod.Spec.Containers[0].Args = append(pod.Spec.Containers[0].Args, "--tls_key_file", filepath.Join(secretMountPath, "tls.key"), "--tls_cert_file", filepath.Join(secretMountPath, "tls.crt"))
	}

	if changedMounts {
		pod.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{}
		pod.Spec.Volumes = []corev1.Volume{}
		for _, mount := range mounts {
			pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, mount)
		}
		for _, volume := range volumes {
			pod.Spec.Volumes = append(pod.Spec.Volumes, volume)
		}
	}

	if newPod {
		return pod, r.Create(ctx, pod)
	}

	if equality.Semantic.DeepEqual(oldPodSpec, &pod.Spec) {
		log.Info("pod unchanged, doing nothing")
		return pod, nil
	}
	log.Info("pod changed, recreating", "diff", cmp.Diff(*oldPodSpec, pod.Spec))
	// Pods are mostly immutable, so recreate it if the spec changed.
	if err := r.Delete(ctx, pod, client.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
		return nil, err
	}
	return pod, r.Create(ctx, pod)
}

// setupInitialPod creates the initial pod configuration for fields that don't change.
func (r *LemmingReconciler) setupInitialPod(pod *corev1.Pod, lemming *lemmingv1alpha1.Lemming) error {
	pod.ObjectMeta = metav1.ObjectMeta{
		Name:      lemming.Name,
		Namespace: lemming.Namespace,
		Labels:    lemming.Labels,
	}
	pod.Spec.InitContainers = []corev1.Container{{
		Name: "init",
	}}
	pod.Spec.Containers = []corev1.Container{{
		Name: "lemming",
		SecurityContext: &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
		},
	}}

	if err := ctrl.SetControllerReference(lemming, pod, r.Scheme); err != nil {
		return err
	}
	return nil
}

func (r *LemmingReconciler) reconcileService(ctx context.Context, lemming *lemmingv1alpha1.Lemming) error {
	var service corev1.Service
	svcName := fmt.Sprintf("service-%s", lemming.Name)

	err := r.Get(ctx, types.NamespacedName{Name: svcName, Namespace: lemming.Namespace}, &service)
	var newService bool
	if apierrors.IsNotFound(err) {
		service.ObjectMeta = metav1.ObjectMeta{
			Name:      svcName,
			Namespace: lemming.Namespace,
			Labels: map[string]string{
				"name": lemming.Name,
			},
		}
		service.Spec = corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app":  lemming.Name,
				"topo": lemming.Namespace,
			},
		}
		if err := ctrl.SetControllerReference(lemming, &service, r.Scheme); err != nil {
			return err
		}
		newService = true
	} else if err != nil {
		return err
	}
	service.Spec.Ports = []corev1.ServicePort{}
	for name, p := range lemming.Spec.Ports {
		service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
			Name:       name,
			Port:       p.OuterPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(int(p.InnerPort)),
		})
	}
	if len(lemming.Spec.Ports) == 0 && newService {
		return nil
	}
	if len(lemming.Spec.Ports) == 0 {
		return r.Delete(ctx, &service)
	}
	if newService {
		return r.Create(ctx, &service)
	}

	return r.Update(ctx, &service)
}

// SetupWithManager sets up the controller with the Manager.
func (r *LemmingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lemmingv1alpha1.Lemming{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
