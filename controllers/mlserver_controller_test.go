package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("MLServer controller", func() {
	const (
		dummyName       = "dummy-server"
		dummyServerName = "dummy-ml"

		port = int32(8080)

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating MLServer Status", func() {
		var (
			ctx       = context.TODO()
			ns        = &corev1.Namespace{}
			lookupKey = &types.NamespacedName{}
		)

		BeforeEach(func() {
			*ns = corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: "testns-" + randStringRunes(5)},
			}

			err := k8sClient.Create(ctx, ns)
			Expect(err).NotTo(HaveOccurred(), "failed to create test namespace")

			lookupKey = &types.NamespacedName{Name: dummyName, Namespace: ns.Name}
		})

		AfterEach(func() {
			err := k8sClient.Delete(ctx, ns)
			Expect(err).NotTo(HaveOccurred(), "failed to delete test namespace")
		})

		It("should create new MLServer with spec and 3 replica if 3 is provided", func() {
			mlserver := serverv1alpha1.MLServer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: serverv1alpha1.GroupVersion.String(),
					Kind:       serverv1alpha1.KIND,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      dummyName,
					Namespace: ns.Name,
				},
				Spec: serverv1alpha1.MLServerSpec{
					ServerName: dummyServerName,
					Port:       port,
					Replicas:   int32(3),
				},
			}
			Expect(k8sClient.Create(ctx, &mlserver)).NotTo(HaveOccurred(), "failed to create MLServer resource")

			created := &serverv1alpha1.MLServer{}
			// use need to return until getting newly created server
			Eventually(func() bool {
				err := k8sClient.Get(ctx, *lookupKey, created)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(created.Spec.ServerName).Should(Equal(dummyServerName))
			Expect(created.Spec.Port).Should(Equal(port))
			Expect(created.Spec.Replicas).Should(Equal(int32(3)))
		})

		It("should create new MLServer with spec and 1 replica if none is provided", func() {
			mlserver := &serverv1alpha1.MLServer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: serverv1alpha1.GroupVersion.String(),
					Kind:       serverv1alpha1.KIND,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      dummyName,
					Namespace: ns.Name,
				},
				Spec: serverv1alpha1.MLServerSpec{
					ServerName: dummyServerName,
					Port:       port,
				},
			}
			Expect(k8sClient.Create(ctx, mlserver)).NotTo(HaveOccurred(), "failed to create MLServer resource")

			created := &serverv1alpha1.MLServer{}

			// use need to return until getting newly created server
			Eventually(func() bool {
				err := k8sClient.Get(ctx, *lookupKey, created)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(created.Spec.ServerName).Should(Equal(dummyServerName))
			Expect(created.Spec.Port).Should(Equal(port))
			Expect(created.Spec.Replicas).Should(Equal(int32(1)))
		})
	})
})
