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

	"go.openly.dev/pointy"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete

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

	// Reconcile every 10 seconds
	//return ctrl.Result{RequeueAfter: time.Second * 10}, err

}

func (r *TetrisReconciler) EnsureTetris(cr *cachev1alpha1.Tetris, cl client.Client, scheme *runtime.Scheme) (err error) {
	fmt.Println("TetrisReconciler: Ensure Tetris")
	appName := cr.Name
	labels := map[string]string{"app": appName}
	matchLabels := map[string]string{"app": appName}

	err = ensureDeployment(cr, cl, appName, labels, matchLabels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating Deployment:", err)
		return err
	}

	err = r.ensureClusterIp(cr, cl, appName, labels, matchLabels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating ClusterIp:", err)
		return err
	}

	err = r.ensureIngress(cr, cl, appName, labels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating Ingress:", err)
		return err
	}

	err = ensureNodePort(cr, cl, appName, labels, matchLabels)
	if err != nil {
		fmt.Println("TetrisReconciler: Error creating or updating NodePort:", err)
		return err
	}

	return nil
}

func (r *TetrisReconciler) ensureIngress(cr *cachev1alpha1.Tetris, cl client.Client, appName string, labels map[string]string) error {
	ingressName := fmt.Sprintf("%s-ingress", appName)
	ingress := &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{Name: ingressName, Namespace: cr.Namespace}}

	if cr.Spec.Domain == nil {
		return fmt.Errorf("Domain name not defined")
	}

	domain := *cr.Spec.Domain

	var className string = "ngnix"
	pathType := networkingv1.PathTypePrefix

	_, err := ctrl.CreateOrUpdate(context.Background(), cl, ingress, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate Ingress")
		ingress.ObjectMeta.Labels = labels
		ingress.Spec = networkingv1.IngressSpec{
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{domain},
					SecretName: "tetris-secret",
				},
			},
			IngressClassName: &className,
			Rules: []networkingv1.IngressRule{
				{
					Host: domain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "service-tetris",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
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

func (r *TetrisReconciler) ensureClusterIp(cr *cachev1alpha1.Tetris, cl client.Client, appName string, labels map[string]string, matchLabels map[string]string) error {
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

func ensureNodePort(cr *cachev1alpha1.Tetris, c client.Client, appName string, labels map[string]string, matchLabels map[string]string) (err error) {

	nodePortName := fmt.Sprintf("%s-nodeport", appName)
	nodePort := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: nodePortName, Namespace: cr.Namespace}}

	// NodePort disabled by default if not specified in the CRD
	if !pointy.BoolValue(cr.Spec.EnableNodePort, false) {
		// Check if NodePort exixts and delete it only if it does
		np := &v1.Service{}
		if err = c.Get(context.TODO(), types.NamespacedName{Name: nodePortName, Namespace: cr.Namespace}, np); err == nil {
			fmt.Println("TetrisReconciler: Unwanted NodePort found")
			err = deleteNodePort(cr, c, appName)
			if err != nil {
				fmt.Println("TetrisReconciler: Error deleting NodePort:", err)
				return err
			}
		}

		cr.Status.NodePortEnabled = false
		return nil
	}

	// Set default port if not specified, otherwise assign the custom one
	nodePortValue := int32(30000)
	if cr.Spec.NodePortValue != nil {
		nodePortValue = int32(*cr.Spec.NodePortValue)
	}

	_, err = ctrl.CreateOrUpdate(context.Background(), c, nodePort, func() error {
		fmt.Println("TetrisReconciler: CreateOrUpdate NodePort")
		nodePort.ObjectMeta.Labels = labels
		nodePort.Spec = v1.ServiceSpec{
			Type:     "NodePort",
			Selector: matchLabels,
			Ports: []v1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(80),
				NodePort:   nodePortValue,
			},
			},
		}

		cr.Status.NodePortEnabled = true
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
