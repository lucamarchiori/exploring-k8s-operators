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
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var tetrislog = logf.Log.WithName("tetris-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Tetris) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=create;update,versions=v1alpha1,name=mtetris.kb.io,admissionReviewVersions=v1

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Tetris) Default() {
	tetrislog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-cache-tetris-operator-secomind-com-v1alpha1-tetris,mutating=false,failurePolicy=fail,sideEffects=None,groups=cache.tetris-operator.secomind.com,resources=tetris,verbs=create;update,versions=v1alpha1,name=vtetris.kb.io,admissionReviewVersions=v1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Tetris) ValidateCreate() (admission.Warnings, error) {
	tetrislog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.

	err := checkReplicasNumber(int(*r.Spec.Replicas))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Tetris) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	tetrislog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.

	err := checkReplicasNumber(int(*r.Spec.Replicas))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Tetris) ValidateDelete() (admission.Warnings, error) {
	tetrislog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// Ensure that the number of replicas requested is between 1 and 5
func checkReplicasNumber(replicas int) (err error) {
	if replicas >= 1 && replicas <= 5 {
		return nil
	}

	return errors.New("The number of replicas must be between 1 and 5")

}
