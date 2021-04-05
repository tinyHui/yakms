package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *MLServerReconciler) deleteEventFilter(ctx context.Context, log logr.Logger) func(e event.DeleteEvent) bool {
	return func(e event.DeleteEvent) bool {
		ownerReferences := e.Meta.GetOwnerReferences()

		for _, ownerReference := range ownerReferences {
			if ownerReference.Kind == serverv1alpha1.KIND {
				log.Info(fmt.Sprintf("Delete %q invoked", e.Meta.GetName()))
				break
			}
		}

		// we prevent this event goes into reconcile
		return false
	}
}
