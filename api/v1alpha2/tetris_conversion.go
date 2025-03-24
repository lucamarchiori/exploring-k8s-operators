package v1alpha2

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"tetris-operator.github.com/api/v1alpha1"
)

// ConvertTo converts this Tetris to the Hub version (v1alpha1).
func (src *Tetris) ConvertTo(dstRaw conversion.Hub) error {
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("Converting from v1alpha2 to v1alpha1")

	dst := dstRaw.(*v1alpha1.Tetris)

	// Copy metadata with deep copy
	src.ObjectMeta.DeepCopyInto(&dst.ObjectMeta) // Remove redundant assignment

	// Copy scalar fields
	dst.Spec.Domain = src.Spec.Domain
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Status.NodePortEnabled = src.Status.NodePortEnabled

	// Handle NodePort conversion
	if src.Spec.NodePort != nil {
		dst.Spec.EnableNodePort = src.Spec.NodePort.Enabled
		dst.Spec.NodePortValue = src.Spec.NodePort.Port
	} else {
		dst.Spec.EnableNodePort = nil
		dst.Spec.NodePortValue = nil
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha1) to this version.
func (dst *Tetris) ConvertFrom(srcRaw conversion.Hub) error {
	setupLog := ctrl.Log.WithName("setup")
	setupLog.Info("Converting from v1alpha1 to v1alpha2")
	src := srcRaw.(*v1alpha1.Tetris)

	// Copy metadata with deep copy
	src.ObjectMeta.DeepCopyInto(&dst.ObjectMeta)

	// Copy scalar fields
	dst.Spec.Domain = src.Spec.Domain
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Status.NodePortEnabled = src.Status.NodePortEnabled

	// Initialize NodePort if needed
	if dst.Spec.NodePort == nil {
		dst.Spec.NodePort = &NodePort{}
	}

	// Handle potential nil values in source
	if src.Spec.EnableNodePort != nil {
		dst.Spec.NodePort.Enabled = src.Spec.EnableNodePort
	}
	if src.Spec.NodePortValue != nil {
		dst.Spec.NodePort.Port = src.Spec.NodePortValue
	}

	return nil
}
