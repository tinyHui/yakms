/*


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

package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	serverv1alpha1 "github.com/tinyHui/yakms/api/v1alpha1"
)

// MLServerReconciler reconciles a MLServer object
type MLServerReconciler struct {
	client.Client
	Log logr.Logger

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=server.yakms.io,resources=mlservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=server.yakms.io,resources=mlservers/status,verbs=get;update;patch

func (r *MLServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("mlserver", req.NamespacedName)
	c := r.Client

	var mlServer serverv1alpha1.MLServer
	if err := c.Get(ctx, req.NamespacedName, &mlServer); err != nil {
		log.Error(err, "unable to fetch MLServers")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log = log.WithValues("deployment_name", mlServer.Spec.ServerName)

	_, err := r.createServiceIfNotExist(ctx, log, mlServer, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.createDeploymentIfNotExist(ctx, log, mlServer, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MLServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&appsv1.Deployment{}, ownerKey, func(rawObj runtime.Object) []string {
		job := rawObj.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != serverv1alpha1.KIND {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&serverv1alpha1.MLServer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
