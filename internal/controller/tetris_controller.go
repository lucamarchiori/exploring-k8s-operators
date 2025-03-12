/*
Copyright 2025.

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

package controller

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	cachev1alpha1 "tetris-operator.github.com/api/v1alpha1"
)

// TetrisReconciler reconciles a Tetris object
type TetrisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.tetris-operator.secomind.com,resources=tetris/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.tetris-operator.secomind.com,resources=tetris/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *TetrisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the Tetris instance
	instance := &cachev1alpha1.Tetris{}

	fmt.Println("TetrisReconciler: Get Tetris instance")
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	r.EnsureTetris(instance, r.Client, r.Scheme)

	return ctrl.Result{}, nil
}

func (r *TetrisReconciler) EnsureTetris(cr *cachev1alpha1.Tetris, client client.Client, scheme *runtime.Scheme) {
	fmt.Println("TetrisReconciler: Ensure Tetris")
	deploymentName := "tetris-deployment"
	labels := map[string]string{
		"app": deploymentName,
	}
	matchLabels := map[string]string{"app": deploymentName}

	deployment := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: deploymentName, Namespace: cr.Namespace}}

	_, err := ctrl.CreateOrUpdate(context.Background(), client, deployment, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate Deployment")
		deployment.ObjectMeta.Labels = labels
		deployment.Spec = tetrisDeploymentSpec(cr, matchLabels, *cr.Spec.Replicas)
		by, _ := json.Marshal(deployment)
		fmt.Println(string(by))
		return nil
	})

	if err != nil {
		fmt.Println("TetrisReconciler: Error CreateOrUpdate Deployment", err)
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *TetrisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Tetris{}).
		Complete(r)
}

func tetrisDeploymentSpec(cr *cachev1alpha1.Tetris, matchLabels map[string]string, instances int32) (deploymentSpec appsv1.DeploymentSpec) {
	fmt.Println("TetrisReconciler: Tetris Deployment Spec")
	deploymentSpec = appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
		Replicas: &instances,
		Template: v1.PodTemplateSpec{
			Spec: tetrisPodSpec(cr),
			ObjectMeta: metav1.ObjectMeta{
				Labels: matchLabels,
			},
		},
	}

	return deploymentSpec
}

func tetrisPodSpec(cr *cachev1alpha1.Tetris) (podSpec v1.PodSpec) {
	fmt.Println("TetrisReconciler: Tetris Pod Spec")
	podSpec = v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  "tetris",
				Image: "annopaolo/tetris",
			},
		},
	}

	return podSpec

}
