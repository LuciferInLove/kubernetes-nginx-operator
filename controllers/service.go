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

	nginxv1beta1 "github.com/LuciferInLove/kubernetes-nginx-operator/api/v1beta1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NginxReconciler) serviceReconcile(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("service")

	// Check if the Service already exists and create one if it doesn't
	service := core.Service{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name}, &service)
	if apierrors.IsNotFound(err) {
		logger.Info("Could not find existing Service for " + nginx.Name + ", creating...")

		service = *createService(nginx)
		if err := r.Client.Create(ctx, &service); err != nil {
			logger.Error(err, "Failed to create Service.")
			return err
		}

		logger.Info("Created Service " + service.Name)
		return nil
	}

	if err != nil {
		logger.Error(err, "Failed to get Service "+nginx.Name)
		return err
	}

	return nil
}

func createService(nginx nginxv1beta1.Nginx) *core.Service {
	service := core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            nginx.Name,
			Namespace:       nginx.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&nginx, nginxv1beta1.GroupVersion.WithKind("Nginx"))},
		},
		Spec: core.ServiceSpec{
			ClusterIP: "None",
			Type:      "ClusterIP",
			Selector: map[string]string{
				"custom-nginx.org/deployment-name": nginx.Name,
			},
			Ports: []core.ServicePort{
				{
					Port: 80,
					Name: nginx.Name,
				},
			},
		},
	}
	return &service
}
