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
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NginxReconciler) deploymentReconcile(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("deployment")

	// Check if the Deployment already exists and create one if it doesn't
	deployment := apps.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name}, &deployment)
	if apierrors.IsNotFound(err) {
		logger.Info("Could not find existing Deployment for " + nginx.Name + ", creating...")

		deployment = *createDeployment(nginx)
		if err := r.Client.Create(ctx, &deployment); err != nil {
			logger.Error(err, "Failed to create Deployment.")
			return err
		}

		logger.Info("Created Deployment " + deployment.Name)
		return nil
	}

	if err != nil {
		logger.Error(err, "Failed to get Deployment "+nginx.Name)
		return err
	}

	// Ensure that the deployment has the same replicas amount as in the spec
	replicas := nginx.Spec.Replicas
	if *deployment.Spec.Replicas != replicas {
		logger.Info("Changing Deployment " + nginx.Name + " replicas number")
		deployment.Spec.Replicas = &replicas
		err = r.Client.Update(ctx, &deployment)
		if err != nil {
			logger.Error(err, "Failed to update Deployment "+nginx.Name)
			return err
		}

		logger.Info("Deployment " + nginx.Name + " updated successfully")
	}

	return nil
}

func createDeployment(nginx nginxv1beta1.Nginx) *apps.Deployment {
	deployment := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            nginx.Name,
			Namespace:       nginx.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&nginx, nginxv1beta1.GroupVersion.WithKind("Nginx"))},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &nginx.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"custom-nginx.org/deployment-name": nginx.Name,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"custom-nginx.org/deployment-name": nginx.Name,
					},
				},
				Spec: core.PodSpec{
					ServiceAccountName: nginx.Spec.ServiceAccount,
					Containers: []core.Container{
						{
							Name:  "nginx",
							Image: nginx.Spec.Image,
							Ports: []core.ContainerPort{
								{
									Name:          "nginx",
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}
