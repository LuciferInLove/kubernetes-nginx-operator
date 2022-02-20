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
	networking "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NginxReconciler) ingressReconcile(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("ingress")

	// Check if the Ingress already exists and create one if it doesn't
	ingress := networking.Ingress{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name}, &ingress)
	if apierrors.IsNotFound(err) {
		logger.Info("Could not find existing Ingress for " + nginx.Name + ", creating...")

		ingress = *createIngress(nginx)
		if err := r.Client.Create(ctx, &ingress); err != nil {
			logger.Error(err, "Failed to create Ingress.")
			return err
		}

		logger.Info("Created Ingress " + ingress.Name)
		return nil
	}

	if err != nil {
		logger.Error(err, "Failed to get Ingress "+nginx.Name)
		return err
	}

	return nil
}

func createIngress(nginx nginxv1beta1.Nginx) *networking.Ingress {
	var pathImplementationSpecific networking.PathType = networking.PathTypeImplementationSpecific

	service := networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:            nginx.Name,
			Namespace:       nginx.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&nginx, nginxv1beta1.GroupVersion.WithKind("Nginx"))},
		},
		Spec: networking.IngressSpec{
			Rules: []networking.IngressRule{
				{
					Host: nginx.Spec.Host,
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathImplementationSpecific,
									Backend: networking.IngressBackend{
										Service: &networking.IngressServiceBackend{
											Name: nginx.Name,
											Port: networking.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []networking.IngressTLS{
				{
					Hosts:      []string{nginx.Spec.Host},
					SecretName: nginx.Name + "-tls",
				},
			},
		},
	}
	return &service
}
