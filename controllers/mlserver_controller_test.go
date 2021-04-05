package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
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

			stopCh chan struct{}
		)

		BeforeEach(func() {
			stopCh = make(chan struct{})

			*ns = corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: "testns-" + randStringRunes(5)},
			}

			err := k8sClient.Create(ctx, ns)
			Expect(err).NotTo(HaveOccurred(), "failed to create test namespace")

			lookupKey = &types.NamespacedName{Name: dummyName, Namespace: ns.Name}

			mgr, err := ctrl.NewManager(cfg, ctrl.Options{})
			Expect(err).NotTo(HaveOccurred(), "failed to create manager")

			controller := &MLServerReconciler{
				Client:   mgr.GetClient(),
				Log:      logf.Log,
				Context:  ctx,
				Recorder: mgr.GetEventRecorderFor("test-controller"),
			}

			err = controller.SetupWithManager(mgr)
			Expect(err).NotTo(HaveOccurred(), "failed to setup controller")

			go func() {
				err := mgr.Start(stopCh)
				Expect(err).NotTo(HaveOccurred(), "failed to start manager")
			}()
		})

		AfterEach(func() {
			close(stopCh)

			err := k8sClient.Delete(ctx, ns)
			Expect(err).NotTo(HaveOccurred(), "failed to delete test namespace")
		})

		It("should create new MLServer creates Service and Deployment", func() {
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
					Replicas:   int32(1),
				},
			}
			Expect(k8sClient.Create(ctx, mlserver)).NotTo(HaveOccurred(), "failed to create MLServer resource")

			// check service created
			service := &corev1.Service{}
			// use need to return until getting newly created resources
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: dummyServerName, Namespace: ns.Name}, service)
			}, timeout, interval).Should(BeNil())
			Expect(len(service.Spec.Ports)).Should(Equal(1))
			Expect(service.Spec.Ports[0].Port).Should(Equal(port))
			Expect(service.Spec.Selector).Should(HaveKeyWithValue("app", dummyServerName))

			// check deployment created
			deployment := &appsv1.Deployment{}
			// use need to return until getting newly created resources
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: dummyServerName, Namespace: ns.Name}, deployment)
			}, timeout, interval).Should(BeNil())
			Expect(len(deployment.Spec.Template.Spec.Containers)).Should(Equal(1))
			Expect(deployment.Spec.Selector.MatchLabels).Should(HaveKeyWithValue("app", dummyServerName))
			Expect(deployment.Spec.Template.ObjectMeta.Labels).Should(HaveKeyWithValue("app", dummyServerName))
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

			// wait the resources starts
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: dummyServerName, Namespace: ns.Name}, deployment)
			}, timeout, interval).Should(BeNil())
			Expect(deployment.Spec).ShouldNot(BeNil())

			mlServer := &serverv1alpha1.MLServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, *lookupKey, mlServer)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(mlServer.Spec.ServerName).Should(Equal(dummyServerName))
			Expect(mlServer.Spec.Port).Should(Equal(port))
			Expect(mlServer.Spec.Replicas).Should(Equal(int32(3)))
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

			// wait the resources starts
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: dummyServerName, Namespace: ns.Name}, deployment)
			}, timeout, interval).Should(BeNil())
			Expect(deployment.Spec).ShouldNot(BeNil())

			// check MLServer
			mlServer := &serverv1alpha1.MLServer{}
			Eventually(func() error {
				return k8sClient.Get(ctx, *lookupKey, mlServer)
			}, timeout, interval).Should(BeNil())
			Expect(mlServer.Spec.ServerName).Should(Equal(dummyServerName))
			Expect(mlServer.Spec.Port).Should(Equal(port))
			Expect(mlServer.Spec.Replicas).Should(Equal(int32(1)))
		})

	})
})
