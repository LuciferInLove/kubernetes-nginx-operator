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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NginxSpec defines the desired state of Nginx
type NginxSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default:=2
	Replicas int32 `json:"replicas,omitempty"`

	Host              string `json:"host,omitempty"`
	CertManagerIssuer string `json:"certManagerIssuer,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="nginx:latest"
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="default"
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// NginxStatus defines the observed state of Nginx
type NginxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Nginx is the Schema for the nginxes API
type Nginx struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxSpec   `json:"spec,omitempty"`
	Status NginxStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NginxList contains a list of Nginx
type NginxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nginx `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nginx{}, &NginxList{})
}
