package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	apiv1alpha1 "tetris-operator.github.com/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Go imports
// log is for logging in this package.

type TetrisCustomDefaulter struct {
	DefaultReplicas int32
	Domain          string
	EnableNodePort  bool
	NodePortValue   int32
}

// SetupWebhookWithManager will setup the manager to manage the webhooks
// +kubebuilder:webhook:path=/mutate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=create;update,versions=v1alpha1,name=mv1alpha1tetris.kb.io,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=false,failurePolicy=fail,groups=cache.tetris-operator.secomind.com,resources=tetris,versions=v1alpha1,name=vv1alpha1tetris.kb.io,sideEffects=None,admissionReviewVersions=v1

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&apiv1alpha1.Tetris{}).
		WithValidator(&TetrisCustomValidator{}).
		WithDefaulter(&TetrisCustomDefaulter{
			DefaultReplicas: 2,
			Domain:          "mytertriscustomdomain.com",
			EnableNodePort:  false,
			NodePortValue:   30000,
		}).
		Complete()
}

var _ webhook.CustomDefaulter = &TetrisCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Tetris.
func (d *TetrisCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	cronjob, ok := obj.(*apiv1alpha1.Tetris)

	if !ok {
		return fmt.Errorf("expected an CronJob object but got %T", obj)
	}

	// Set default values
	d.applyDefaults(cronjob)
	return nil
}

// applyDefaults applies default values to CronJob fields.
func (d *TetrisCustomDefaulter) applyDefaults(t *apiv1alpha1.Tetris) {
	// Set default Replicas if not specified
	if t.Spec.Replicas == nil {
		t.Spec.Replicas = &d.DefaultReplicas
	}

	// Set default Domain if not specified
	if t.Spec.Domain == nil {
		t.Spec.Domain = &d.Domain
	}

	if t.Spec.EnableNodePort == nil {
		t.Spec.EnableNodePort = &d.EnableNodePort
	}

	if t.Spec.NodePortValue == nil {
		t.Spec.NodePortValue = &d.NodePortValue
	}

}

type TetrisCustomValidator struct {
	Replicas       int32
	Domain         string
	EnableNodePort bool
	NodePortValue  int32
}

var _ webhook.CustomValidator = &TetrisCustomValidator{}

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
