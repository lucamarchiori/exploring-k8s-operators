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
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
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

func (r *TetrisReconciler) EnsureTetris(cr *cachev1alpha1.Tetris, cl client.Client, scheme *runtime.Scheme) (err error) {
	fmt.Println("TetrisReconciler: Ensure Tetris")
	appName := "tetris"
	labels := map[string]string{"app": appName}
	matchLabels := map[string]string{"app": appName}

	err = ensureDeployment(cr, cl, appName, labels, matchLabels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating Deployment:", err)
		return err
	}

	err = r.EnsureClusterIp(cr, cl, appName, labels, matchLabels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating ClusterIp:", err)
		return err
	}

	if cr.Spec.EnableNodePort {
		err = ensureNodePort(cr, cl, appName, labels, matchLabels)
		if err != nil {
			fmt.Println("TetrisReconciler: Error creating or updating NodePort:", err)
			return err
		}
	} else {
		err = deleteNodePort(cr, cl, appName)
		if err != nil {
			fmt.Println("TetrisReconciler: Error deleting NodePort:", err)
			return err
		}
	}

	return nil
}

func (r *TetrisReconciler) EnsureClusterIp(cr *cachev1alpha1.Tetris, cl client.Client, appName string, labels map[string]string, matchLabels map[string]string) error {
	clusterIpName := fmt.Sprintf("%s-clusterip", appName)
	clusterIp := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: clusterIpName, Namespace: cr.Namespace}}

	_, err := ctrl.CreateOrUpdate(context.Background(), cl, clusterIp, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate NodePort")
		clusterIp.ObjectMeta.Labels = labels
		clusterIp.Spec = v1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: matchLabels,
			Ports: []v1.ServicePort{{
				Protocol:   "TCP",
				Port:       80,
				TargetPort: intstr.FromInt(80),
			},
			},
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func deleteNodePort(cr *cachev1alpha1.Tetris, cl client.Client, appName string) error {
	nodePortName := fmt.Sprintf("%s-nodeport", appName)
	nodePort := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: nodePortName, Namespace: cr.Namespace}}
	err := cl.Delete(context.Background(), nodePort)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	fmt.Println("TetrisReconciler: Successfully deleted NodePort")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TetrisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Tetris{}).
		Complete(r)
}

func ensureNodePort(cr *cachev1alpha1.Tetris, client client.Client, appName string, labels map[string]string, matchLabels map[string]string) error {
	nodePortName := fmt.Sprintf("%s-nodeport", appName)
	nodePort := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: nodePortName, Namespace: cr.Namespace}}

	_, err := ctrl.CreateOrUpdate(context.Background(), client, nodePort, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate NodePort")
		nodePort.ObjectMeta.Labels = labels
		nodePort.Spec = v1.ServiceSpec{
			Type:     "NodePort",
			Selector: matchLabels,
			Ports: []v1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(80),
				NodePort:   int32(30000),
			},
			},
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func ensureDeployment(cr *cachev1alpha1.Tetris, client client.Client, appName string, labels map[string]string, matchLabels map[string]string) error {
	deploymentName := fmt.Sprintf("%s-deployment", appName)
	deployment := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: deploymentName, Namespace: cr.Namespace}}

	_, err := ctrl.CreateOrUpdate(context.Background(), client, deployment, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate Deployment")
		deployment.ObjectMeta.Labels = labels
		deployment.Spec = tetrisDeploymentSpec(matchLabels, *cr.Spec.Replicas)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func tetrisDeploymentSpec(matchLabels map[string]string, instances int32) (deploymentSpec appsv1.DeploymentSpec) {
	fmt.Println("TetrisReconciler: Tetris Deployment Spec")
	deploymentSpec = appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
		Replicas: &instances,
		Template: v1.PodTemplateSpec{
			Spec: tetrisPodSpec(),
			ObjectMeta: metav1.ObjectMeta{
				Labels: matchLabels,
			},
		},
	}

	return deploymentSpec
}

func tetrisPodSpec() (podSpec v1.PodSpec) {
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
