/*
Copyright 2022.

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
	"sort"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lemmingv1alpha1 "github.com/openconfig/lemming/operator/api/lemming/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LemmingReconciler reconciles a Lemming object
type LemmingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lemming.openconfig.net,resources=lemmings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Lemming object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
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
	deploy := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deploy)
	var newDeployment bool

	if apierrors.IsNotFound(err) {
		if r.setupInitialDeployment(deploy, lemming); err != nil {
			return ctrl.Result{}, err
		}
		newDeployment = true
	} else if err != nil {
		log.Error(err, "unable to get deployment")
		return ctrl.Result{}, err
	}
	deploy.Spec.Template.Spec.Containers[0].Image = lemming.Spec.Image
	deploy.Spec.Template.Spec.Containers[0].Command = []string{lemming.Spec.Command}
	deploy.Spec.Template.Spec.Containers[0].Args = lemming.Spec.Args
	envs := make([]corev1.EnvVar, 0, len(lemming.Spec.Env))
	for k, v := range lemming.Spec.Env {
		envs = append(envs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	sort.Slice(envs, func(i, j int) bool {
		return envs[i].Name < envs[j].Name
	})
	deploy.Spec.Template.Spec.Containers[0].Env = envs

	if newDeployment {
		err = r.Create(ctx, deploy)
	} else {
		err = r.Update(ctx, deploy)
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *LemmingReconciler) setupInitialDeployment(deploy *appsv1.Deployment, lemming *lemmingv1alpha1.Lemming) error {
	deploy.ObjectMeta = metav1.ObjectMeta{
		Name:      lemming.Name,
		Namespace: lemming.Namespace,
	}
	if err := ctrl.SetControllerReference(lemming, deploy, r.Scheme); err != nil {
		return err
	}
	deploy.Spec = appsv1.DeploymentSpec{
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		},
		Replicas: pointer.Int32(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"name": lemming.Name,
				"type": "lemming",
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name": lemming.Name,
					"type": "lemming",
				},
			},
		},
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LemmingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lemmingv1alpha1.Lemming{}).
		Complete(r)
}
