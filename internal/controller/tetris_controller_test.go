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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.openly.dev/pointy"
	v1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cachev1alpha1 "tetris-operator.github.com/api/v1alpha1"
)

var _ = Describe("Tetris Controller", func() {

	// Define tetris constants to use as test
	const (
		ResourceName      = "test-tetris"
		ResourceNamespace = "default"
		EnableNodePort    = true
		NodePortValue     = 30003
		Replicas          = 2
		Domain            = "customtetrisdomain.com"
		APIVersion        = "v1alpha1"
		Kind              = "Tetris"
	)

	Context("When reconciling a resource", func() {
		ctx := context.Background()
		var tetris *cachev1alpha1.Tetris

		BeforeEach(func() {
			By("creating the custom resource for the Kind Tetris")

			tetris = &cachev1alpha1.Tetris{
				TypeMeta: metav1.TypeMeta{
					APIVersion: APIVersion,
					Kind:       Kind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: ResourceNamespace,
				},
				Spec: cachev1alpha1.TetrisSpec{
					EnableNodePort: pointy.Bool(EnableNodePort),
					NodePortValue:  pointy.Int32(NodePortValue),
					Replicas:       pointy.Int32(Replicas),
					Domain:         pointy.String(Domain),
				}}

			Expect(k8sClient.Create(ctx, tetris)).To(Succeed())
		})

		It("should create CR", func() {
			crKey := types.NamespacedName{Name: tetris.Name, Namespace: tetris.Namespace}
			createdCR := &cachev1alpha1.Tetris{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, crKey, createdCR)
				return err == nil
			}, 10*time.Second, 250*time.Millisecond).Should(BeTrue())
		})

		It("should create deployment", func() {
			deployKey := types.NamespacedName{
				Name:      tetris.Name + "-deployment", // Add "-deployment" suffix
				Namespace: tetris.Namespace,
			}
			createdDeploy := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deployKey, createdDeploy)
				return err == nil
			}, time.Second*10).Should(BeTrue())
		})

		It("should create NodePort", func() {
			deployKey := types.NamespacedName{
				Name:      tetris.Name + "-nodeport", // Add "-deployment" suffix
				Namespace: tetris.Namespace,
			}
			np := &v1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deployKey, np)
				return err == nil
			}, time.Second*10).Should(BeTrue())
		})

		AfterEach(func() {
			resource := &cachev1alpha1.Tetris{}
			deployKey := types.NamespacedName{Name: tetris.Name, Namespace: tetris.Namespace}
			err := k8sClient.Get(ctx, deployKey, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Tetris")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
	})
})
