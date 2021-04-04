package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ownerKey = ".metadata.controller"
	apiGVStr = serverv1alpha1.GroupVersion.String()
)

func (r *MLServerReconciler) createDeploymentIfNotExist(ctx context.Context, log logr.Logger, mlServer serverv1alpha1.MLServer, req ctrl.Request) error {
	// create deployment if not found
	deploymentSpec := appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: mlServer.Namespace, Name: mlServer.Spec.ServerName}, &deploymentSpec)
	if apierrors.IsNotFound(err) {
		log.Info(fmt.Sprintf("could not find existing deployment for %s, creating...", serverv1alpha1.KIND))

		deploymentSpec = constructDeploymentSpec(mlServer)
		if err := r.Create(ctx, &deploymentSpec); err != nil {
			log.Error(err, fmt.Sprintf("unable to create Deployment for %s", serverv1alpha1.KIND))
			return nil
		}

		r.Recorder.Eventf(&mlServer, corev1.EventTypeNormal, "Created", "created deployment %q", deploymentSpec.Name)
		log.Info("deployment created")

		return nil
	}

	if err != nil {
		log.Error(err, fmt.Sprintf("failed to get Deployment for %q", serverv1alpha1.KIND))
		return err
	}

	log.Info(fmt.Sprintf("deployment %q exist, skip", req.Name))

	return nil
}
