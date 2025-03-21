package v1alpha2

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"tetris-operator.github.com/api/v1alpha1"
)

// ConvertTo converts this Tetris to the Hub version (v1alpha1).
func (src *Tetris) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha1.Tetris)
	dst.Spec.Domain = src.Spec.Domain
	dst.Spec.EnableNodePort = src.Spec.NodePort.Enabled
	dst.Spec.NodePortValue = src.Spec.NodePort.Port
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Status.NodePortEnabled = src.Status.NodePortEnabled
	dst.ObjectMeta = src.ObjectMeta
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha1) to this version.
func (dst *Tetris) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha1.Tetris)
	dst.Spec.Domain = src.Spec.Domain
	dst.Spec.NodePort.Enabled = src.Spec.EnableNodePort
	dst.Spec.NodePort.Port = src.Spec.NodePortValue
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Status.NodePortEnabled = src.Status.NodePortEnabled
	dst.ObjectMeta = src.ObjectMeta
	return nil
}
