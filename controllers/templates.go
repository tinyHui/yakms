package controllers

import (
	"fmt"
	"github.com/tinyHui/yakms/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
)

func constructJobSpec(server v1alpha1.MLServer) batchv1.Job {
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
			Name:        "ml-server",
			Namespace:   server.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers:    []v1.Container{{Name: "ml-server", Image: "hello-world"}},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
		},
	}
}

func constructDeploymentSpec(server v1alpha1.MLServer) appsv1.Deployment {
	var (
		deploymentNameKey   = "mlservers.server.yakms.io/deployment-name"
		deploymentNameValue = fmt.Sprintf("%s-deployment", server.Spec.ServerName)

		labels = map[string]string{
			"app":             server.Spec.ServerName,
			deploymentNameKey: deploymentNameValue,
		}
	)

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            server.Spec.ServerName,
			Namespace:       server.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&server, serverv1alpha1.GroupVersion.WithKind(serverv1alpha1.KIND))},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &server.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "hello-world-nginx",
							Image: "paulbouwer/hello-kubernetes:1.9",
							Ports: []v1.ContainerPort{{
								ContainerPort: 8080,
							}},
						},
					},
				},
			},
		},
	}
}
