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

package v1alpha1

import (
	"context"
	"fmt"

	"go.openly.dev/pointy"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.

type TetrisCustomDefaulter struct {
	Replicas       int32
	Domain         string
	EnableNodePort bool
	NodePortValue  int32
}

type TetrisCustomValidator struct {
	Replicas       int32
	Domain         string
	EnableNodePort bool
	NodePortValue  int32
}

var _ webhook.CustomValidator = &TetrisCustomValidator{}
var _ webhook.CustomDefaulter = &TetrisCustomDefaulter{}

// SetupWebhookWithManager will setup the manager to manage the webhooks
// +kubebuilder:webhook:path=/mutate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=create;update,versions=v1alpha1,name=mtetris.kb.io,admissionReviewVersions=v1
// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=false,failurePolicy=fail,groups=cache.tetris-operator.secomind.com,resources=tetris,versions=v1alpha1,name=vtetris.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *Tetris) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(&TetrisCustomValidator{}).
		WithDefaulter(&TetrisCustomDefaulter{
			Replicas:       2,
			Domain:         "mytertriscustomdomain.com",
			EnableNodePort: false,
			NodePortValue:  30000,
		}).
		Complete()
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Tetris.
func (d *TetrisCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	cronjob, ok := obj.(*Tetris)

	if !ok {
		return fmt.Errorf("expected an CronJob object but got %T", obj)
	}

	// Set default values
	d.applyDefaults(cronjob)
	return nil
}

// applyDefaults applies default values to CronJob fields.
func (d *TetrisCustomDefaulter) applyDefaults(cronJob *Tetris) {
	// Set default Replicas if not specified
	if cronJob.Spec.Replicas == nil {
		cronJob.Spec.Replicas = pointy.Int32(d.Replicas)
	}

	// Set default Domain if not specified
	if cronJob.Spec.Domain == nil {
		cronJob.Spec.Domain = pointy.String(d.Domain)
	}

	if cronJob.Spec.EnableNodePort == nil {
		cronJob.Spec.EnableNodePort = &d.EnableNodePort
	}

	if cronJob.Spec.NodePortValue == nil {
		cronJob.Spec.NodePortValue = &d.NodePortValue
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
