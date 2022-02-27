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
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NginxReconciler) deleteSecret(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("secret-delete")
	secret := core.Secret{}

	// Check if the target secret exists
	logger.Info("Checking if the target secret exists before deleting...")
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name + "-tls"}, &secret)
	if err != nil {
		return err
	}

	// Delete secret
	logger.Info("Deleting secret " + nginx.Name + "-tls")
	if err := r.Client.DeleteAllOf(ctx, &secret, client.InNamespace(nginx.Namespace), client.MatchingFieldsSelector{
		Selector: fields.AndSelectors(
			fields.OneTermNotEqualSelector("metadata.name", nginx.Name),
		),
	}); err != nil {
		return err
	}

	return nil
}
