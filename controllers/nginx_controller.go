/*
Copyright 2022 LuciferInLove.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	nginxv1beta1 "github.com/LuciferInLove/kubernetes-nginx-operator/api/v1beta1"
)

// NginxReconciler reconciles an Nginx object
type NginxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nginx.custom-nginx.org,resources=nginxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nginx.custom-nginx.org,resources=nginxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nginx.custom-nginx.org,resources=nginxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *NginxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logger := log.Log.WithName("controller").WithName("nginx")

	nginx := nginxv1beta1.Nginx{}

	if err := r.Client.Get(ctx, req.NamespacedName, &nginx); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Nginx resource is not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Nginx resource.")
		return ctrl.Result{}, err
	}

	// Add finalizer to delete external resources
	secretFinalizer := "custom-nginx.org/finalizer"

	// Check if Nginx resource is not under deletion and add finalizer
	if nginx.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&nginx, secretFinalizer) {
			controllerutil.AddFinalizer(&nginx, secretFinalizer)
			if err := r.Update(ctx, &nginx); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// If resource is under deletion, delete secret and then delete the finalizer
		if controllerutil.ContainsFinalizer(&nginx, secretFinalizer) {
			// Remove secret
			if err := r.deleteSecret(ctx, nginx); err != nil {
				return ctrl.Result{}, err
			}

			// Remove the finalizer
			controllerutil.RemoveFinalizer(&nginx, secretFinalizer)
			if err := r.Update(ctx, &nginx); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	err := r.deploymentReconcile(ctx, nginx)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.serviceReconcile(ctx, nginx)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.ingressReconcile(ctx, nginx)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.certificateReconcile(ctx, nginx)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nginxv1beta1.Nginx{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				return false
			},
		}).
		Complete(r)
}
