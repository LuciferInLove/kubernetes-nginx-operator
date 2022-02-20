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
	cert "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NginxReconciler) certificateReconcile(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("certificate")

	// Check if the Certificate already exists and create one if it doesn't
	certificate := cert.Certificate{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name}, &certificate)
	if apierrors.IsNotFound(err) {
		logger.Info("Could not find existing Certificate for " + nginx.Name + ", creating...")

		certificate = *createCertificate(nginx)
		if err := r.Client.Create(ctx, &certificate); err != nil {
			logger.Error(err, "Failed to create Certificate.")
			return err
		}

		logger.Info("Created Certificate " + certificate.Name)
		return nil
	}

	if err != nil {
		logger.Error(err, "Failed to get Certificate "+nginx.Name)
		return err
	}

	return nil
}

func createCertificate(nginx nginxv1beta1.Nginx) *cert.Certificate {
	certificate := cert.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:            nginx.Name,
			Namespace:       nginx.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&nginx, nginxv1beta1.GroupVersion.WithKind("Nginx"))},
		},
		Spec: cert.CertificateSpec{
			DNSNames: []string{nginx.Spec.Host},
			IssuerRef: cmmeta.ObjectReference{
				Name: nginx.Spec.CertManagerIssuer,
			},
			SecretName: nginx.Name + "-tls",
		},
	}
	return &certificate
}

func (r *NginxReconciler) deleteSecret(ctx context.Context, nginx nginxv1beta1.Nginx) error {
	logger := log.Log.WithName("controller").WithName("certificatdelete")
	secret := core.Secret{}

	// Check if the target secret exists
	logger.Info("Checking if the target secret exists before deleting...")
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: nginx.Namespace, Name: nginx.Name + "-tls"}, &secret)
	if apierrors.IsNotFound(err) {
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

	if err != nil {
		return err
	}

	return nil
}
