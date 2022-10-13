/*
Copyright 2022 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lemmingv1alpha1 "github.com/openconfig/lemming/operator/api/lemming/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// LemmingReconciler reconciles a Lemming object
type LemmingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

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

	if err := r.reconcilePod(ctx, lemming); err != nil {
		log.Error(err, "unable to get reconcile pod")
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, lemming); err != nil {
		log.Error(err, "unable to get reconcile service: %v")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *LemmingReconciler) reconcilePod(ctx context.Context, lemming *lemmingv1alpha1.Lemming) error {
	log := log.FromContext(ctx)
	pod := &corev1.Pod{}
	err := r.Get(ctx, types.NamespacedName{Name: lemming.Name, Namespace: lemming.Namespace}, pod)
	var newPod bool

	if apierrors.IsNotFound(err) {
		log.Info("new pod, creating initial spec")
		if err := r.setupInitialPod(pod, lemming); err != nil {
			return fmt.Errorf("failed to setup initial deployment: %v", err)
		}
		newPod = true
	} else if err != nil {
		return err
	}

	oldPodSpec := pod.Spec.DeepCopy()
	pod.Spec.Containers[0].Image = lemming.Spec.Image
	pod.Spec.InitContainers[0].Image = lemming.Spec.InitImage
	pod.Spec.Containers[0].Command = []string{lemming.Spec.Command}
	pod.Spec.Containers[0].Args = lemming.Spec.Args
	pod.Spec.Containers[0].Env = lemming.Spec.Env
	pod.Spec.Containers[0].Resources = lemming.Spec.Resources

	if newPod {
		return r.Create(ctx, pod)
	}

	if reflect.DeepEqual(oldPodSpec, &pod.Spec) {
		log.Info("pod unchanged, doing nothing")
		return nil
	}
	// Pods are mostly immutable, so recreate it if the spec changed.
	if err := r.Delete(ctx, pod, client.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
		return err
	}
	return r.Create(ctx, pod)
}

// setupInitialPod creates the initial pod configuration for fields that don't change.
func (r *LemmingReconciler) setupInitialPod(pod *corev1.Pod, lemming *lemmingv1alpha1.Lemming) error {
	pod.ObjectMeta = metav1.ObjectMeta{
		Name:      lemming.Name,
		Namespace: lemming.Namespace,
		Labels: map[string]string{
			"app":  lemming.Name,
			"topo": lemming.Namespace,
		},
	}
	if err := ctrl.SetControllerReference(lemming, pod, r.Scheme); err != nil {
		return err
	}
	pod.Spec = corev1.PodSpec{
		InitContainers: []corev1.Container{{
			Name: "init",
			Args: []string{
				fmt.Sprintf("%d", lemming.Spec.InterfaceCount+1),
				fmt.Sprintf("%d", lemming.Spec.InitSleep),
			},
		}},
		Containers: []corev1.Container{{
			Name: "lemming",
		}},
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
		Complete(r)
}
