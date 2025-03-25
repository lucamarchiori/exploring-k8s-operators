package v1alpha2

import (
	. "github.com/onsi/gomega"
	"go.openly.dev/pointy"
	"tetris-operator.github.com/api/v1alpha1"
	"tetris-operator.github.com/api/v1alpha2"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/api/equality"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Tetris Webhook", func() {

	Context("v1alpha1 <-> v1alpha2 conversions", func() {
		When("converting from v1alpha1 to v1alpha2 and back", func() {
			It("should preserve all fields", func() {
				// Step 1: Create a v1alpha1 Tetris object
				v1alpha1Obj := &v1alpha1.Tetris{
					ObjectMeta: v1.ObjectMeta{
						Name: "test-tetris",
					},
					Spec: v1alpha1.TetrisSpec{
						Replicas:       pointy.Int32(3),
						Domain:         pointy.String("example.com"),
						EnableNodePort: pointy.Bool(true),
						NodePortValue:  pointy.Int32(31000),
					},
				}

				// Step 2: Convert to v1alpha2
				v1alpha2Obj := &v1alpha2.Tetris{}
				Expect(v1alpha2Obj.ConvertFrom(v1alpha1Obj)).To(Succeed())

				// Step 3: Convert back to v1alpha1
				convertedBack := &v1alpha1.Tetris{}
				Expect(v1alpha2Obj.ConvertTo(convertedBack)).To(Succeed())

				// Step 4: Compare original and converted-back objects
				Expect(equality.Semantic.DeepEqual(v1alpha1Obj, convertedBack)).To(BeTrue())

				// Manual check
				Expect(pointy.PointersValueEqual(v1alpha1Obj.Spec.NodePortValue, v1alpha2Obj.Spec.NodePort.Port)).To(BeTrue())
			})
		})
	})
})
