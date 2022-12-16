// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LemmingSpec defines the desired state of Lemming.
type LemmingSpec struct {
	// Image is the container image to run.
	Image string `json:"image,omitempty"`
	// Command is the name of the executable to run.
	Command string `json:"command,omitempty"`
	// Args are the args to pass to the command.
	Args []string `json:"args,omitempty"`
	// Env are the environment variables to set for the container.
	// +listType=map
	// +listMapKey=name
	Env []corev1.EnvVar `json:"env,omitempty"`
	// ConfigPath is the mount point for configuration inside the pod.
	ConfigPath string `json:"configPath,omitempty"`
	// ConfigFile is the default configuration file name for the pod.
	ConfigFile string `json:"configFile,omitempty"`
	// InitImage is the docker image to use as an init container for the pod.
	InitImage string `json:"initImage,omitempty"`
	// Ports are ports to create on the service.
	Ports map[string]ServicePort `json:"ports,omitempty"`
	// InterfaceCount is number of interfaces to be attached to the pod.
	// +optional
	InterfaceCount int `json:"interfaceCount"`
	// InitSleep is the time sleep in the init container
	// +optional
	InitSleep int `json:"initSleep"`
	// Resources are the K8s resources to allocate to lemming container.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources"`
	// TLS is the configuration the key/certs to use for management.
	// +optional
	TLS TLSSpec `json:"tls"`
}

// TLSSpec is the configuration the key/certs to use for management.
type TLSSpec struct {
	// SelfSigned generates a new self signed certificate.
	// +optional
	SelfSigned SelfSignedSpec `json:"selfSigned"`
}

// SelfSignedSpec is the configuration to generate a self-signed cert.
type SelfSignedSpec struct {
	/// Common name to set in the cert.
	CommonName string `json:"commonName"`
	// RSA keysize to use for key generation.
	KeySize int `json:"keySize"`
}

// ServicePort describes an external L4 port on the device.
type ServicePort struct {
	// InnerPort is port on the container to expose.
	InnerPort int32 `json:"innerPort"`
	// OuterPort is port on the container to expose.
	OuterPort int32 `json:"outerPort"`
}

// LemmingPhase is the overall status of the Lemming.
type LemmingPhase string

const (
	// Running indicates a successfully running lemming.
	Running LemmingPhase = "Running"
	// Failed indicates an error state.
	Failed LemmingPhase = "Failed"
	// Unknown indicates an unknown state.
	Unknown LemmingPhase = "Unknown"
	// Pending indicates a pending state.
	Pending LemmingPhase = "Pending"
)

// LemmingStatus defines the observed state of Lemming
type LemmingStatus struct {
	// Phase is the overall status of the Lemming.
	Phase LemmingPhase `json:"phase"`
	// Message describes why the lemming is in the current phase.
	Message string `json:"message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Lemming is the Schema for the lemmings API
type Lemming struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LemmingSpec   `json:"spec,omitempty"`
	Status LemmingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LemmingList contains a list of Lemming
type LemmingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lemming `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lemming{}, &LemmingList{})
}
