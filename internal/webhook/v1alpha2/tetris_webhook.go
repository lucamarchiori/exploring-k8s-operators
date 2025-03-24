package v1alpha2

import (
	"context"
	"fmt"

	"go.openly.dev/pointy"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	apiv1alpha2 "tetris-operator.github.com/api/v1alpha2"
)

// nolint:unused

type TetrisCustomDefaulter struct {
	Replicas int32
	Domain   string
	Nodeport apiv1alpha2.NodePort
}

type TetrisCustomValidator struct {
	Replicas int32
	Domain   string
	Nodeport apiv1alpha2.NodePort
}

var _ webhook.CustomValidator = &TetrisCustomValidator{}
var _ webhook.CustomDefaulter = &TetrisCustomDefaulter{}

// SetupWebhookWithManager will setup the manager to manage the webhooks
// +kubebuilder:webhook:path=/mutate-cache-tetris-operator-secomind-com-v1alpha2-tetris,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=create;update,versions=v1alpha2,name=mtetris.kb.io,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-cache-tetris-operator-secomind-com-v1alpha2-tetris,mutating=false,failurePolicy=fail,groups=cache.tetris-operator.secomind.com,resources=tetris,versions=v1alpha2,name=vtetris.kb.io,sideEffects=None,admissionReviewVersions=v1

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&apiv1alpha2.Tetris{}).
		WithValidator(&TetrisCustomValidator{}).
		WithDefaulter(&TetrisCustomDefaulter{
			Replicas: 2,
			Domain:   "mytertriscustomdomain.com",
			Nodeport: apiv1alpha2.NodePort{
				Enabled: pointy.Bool(false),
				Port:    pointy.Int32(30001),
			},
		}).
		Complete()
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Tetris.
func (d *TetrisCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	t, ok := obj.(*apiv1alpha2.Tetris)

	if !ok {
		return fmt.Errorf("expected an CronJob object but got %T", obj)
	}

	// Set default values
	d.applyDefaults(t)
	return nil
}

func (d *TetrisCustomDefaulter) applyDefaults(cronJob *apiv1alpha2.Tetris) {
	// Set default Replicas if not specified
	if cronJob.Spec.Replicas == nil {
		cronJob.Spec.Replicas = pointy.Int32(d.Replicas)
	}

	// Set default Domain if not specified
	if cronJob.Spec.Domain == nil {
		cronJob.Spec.Domain = pointy.String(d.Domain)
	}

	// Initialize Service if nil
	if cronJob.Spec.NodePort == nil {
		cronJob.Spec.NodePort = &apiv1alpha2.NodePort{
			Enabled: pointy.Bool(*d.Nodeport.Enabled),
			Port:    pointy.Int32(*d.Nodeport.Port),
		}
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TetrisCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TetrisCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TetrisCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
