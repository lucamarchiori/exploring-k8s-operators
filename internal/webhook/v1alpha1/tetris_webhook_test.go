package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.openly.dev/pointy"

	"tetris-operator.github.com/api/v1alpha1"
)

var _ = Describe("Tetris Webhook", func() {
	var (
		obj       *v1alpha1.Tetris
		validator TetrisCustomValidator
		defaulter TetrisCustomDefaulter
	)

	BeforeEach(func() {
		obj = &v1alpha1.Tetris{}
		validator = TetrisCustomValidator{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		defaulter = TetrisCustomDefaulter{
			DefaultReplicas: 2, // Set the correct default
			Domain:          "mytertriscustomdomain.com",
			EnableNodePort:  false,
			NodePortValue:   30000,
		}
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
	})

	Context("When creating Tetris under Defaulting Webhook", func() {
		It("Should fill in the default value if a required field is empty", func() {
			By("simulating a scenario where defaults should be applied")
			obj.Spec.Replicas = nil
			obj.Spec.Domain = nil
			obj.Spec.EnableNodePort = nil
			obj.Spec.NodePortValue = nil
			By("calling the Default method to apply defaults")
			defaulter.Default(ctx, obj)
			By("checking that the default values are set")
			Expect(obj.Spec.Replicas).To(Equal(pointy.Int32(defaulter.DefaultReplicas)))
			Expect(obj.Spec.Domain).To(Equal(pointy.String("mytertriscustomdomain.com")))
			Expect(obj.Spec.EnableNodePort).To(Equal(pointy.Bool(false)))
			Expect(obj.Spec.NodePortValue).To(Equal(pointy.Int32(30000)))
		})
	})

	Context("When creating Tetris under Validating Webhook", func() {
		It("Should deny if a required field is empty", func() {

			// TODO(user): Add your logic here

		})

		It("Should admit if all required fields are provided", func() {

			// TODO(user): Add your logic here

		})
	})
	Context("When creating Tetris under Conversion Webhook", func() {

	})

})
