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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TetrisSpec defines the desired state of Tetris
type TetrisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas,omitempty"`
	// +kubebuilder:validation:Required
	Domain   *string   `json:"domain,omitempty"`
	NodePort *NodePort `json:"nodePort,omitempty"`
}

// TetrisStatus defines the observed state of Tetris
type TetrisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NodePortEnabled bool `json:"nodePortEnabled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Tetris is the Schema for the tetris API
type Tetris struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TetrisSpec   `json:"spec,omitempty"`
	Status TetrisStatus `json:"status,omitempty"`
}

type NodePort struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Enabled *bool `json:"enabled,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum:=30000
	// +kubebuilder:validation:Maximum:=32767
	Port *int32 `json:"port,omitempty"`
}

// +kubebuilder:object:root=true

// TetrisList contains a list of Tetris
type TetrisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tetris `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tetris{}, &TetrisList{})
}
