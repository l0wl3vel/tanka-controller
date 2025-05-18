/*
Copyright 2024 The Flux authors

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

package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TankaReleaseSpec defines the desired state of TankaRelease
type TankaReleaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ChartRef holds a reference to a source controller resource containing the
	// Tanka Bundle artifact.
	TankaRef CrossNamespaceSourceReference `json:"chartRef"`

	// Environment of the Tanka to deploy
	Environment string `json:"environment,omitempty"`

	// Interval at which to reconcile the Tanka release.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	// +required
	Interval metav1.Duration `json:"interval"`

	// Suspend tells the controller to suspend reconciliation for this TankaRelease,
	// it does not apply to already started reconciliations. Defaults to false.
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

const (
	// SourceIndexKey is the key used for indexing HelmReleases based on
	// their sources.
	SourceIndexKey string = ".metadata.source"
)

// TankaReleaseStatus defines the observed state of TankaRelease
type TankaReleaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// TankaRelease is the Schema for the tankareleases API
type TankaRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TankaReleaseSpec   `json:"spec,omitempty"`
	Status TankaReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TankaReleaseList contains a list of TankaRelease
type TankaReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TankaRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TankaRelease{}, &TankaReleaseList{})
}
