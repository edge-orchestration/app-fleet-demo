/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PublisherSpec struct {
	// Port is the port number to expose on the Publisher Pod
	// +kubebuilder:default=4001
	Port *int32 `json:"port,omitempty"`
	// Replicas is the number of deployement replicas to scale
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
}

type ConsumerSpec struct {
	// Port is the port number to expose on the Consumer Pod
	// +kubebuilder:default=4002
	Port *int32 `json:"port,omitempty"`
	// Replicas is the number of deployement replicas to scale
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
}

type RabbitmqSpec struct {
	// Replicas is the number of deployement replicas to scale
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
}

// RmqDemoSpec defines the desired state of RmqDemo
type RmqDemoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Publisher App
	Publisher *PublisherSpec `json:"publisher,omitempty"`

	// Consumer App
	Consumer *ConsumerSpec `json:"consumer,omitempty"`

	// Rabbitmq Server
	Rabbitmq *RabbitmqSpec `json:"rabbitmq,omitempty"`

	// force Re-deploy
	ForceRedeploy bool `json:"forceRedeploy,omitempty"`
}

// RmqDemoStatus defines the observed state of RmqDemo
type RmqDemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RmqDemo is the Schema for the rmqdemoes API
type RmqDemo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RmqDemoSpec   `json:"spec,omitempty"`
	Status RmqDemoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RmqDemoList contains a list of RmqDemo
type RmqDemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RmqDemo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RmqDemo{}, &RmqDemoList{})
}
