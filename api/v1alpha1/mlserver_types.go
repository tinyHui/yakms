/*


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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MLServerSpec defines the desired state of MLServer
type MLServerSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MaxLength=64
	ServerName string `json:"serverName"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`

	HelloWorld string `json:"hello,omitempty"`
}

// MLServerStatus defines the observed state of MLServer
type MLServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// MLServer is the Schema for the mlservers API
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.spec.replicas`
type MLServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MLServerSpec   `json:"spec,omitempty"`
	Status MLServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MLServerList contains a list of MLServer
type MLServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MLServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MLServer{}, &MLServerList{})
}
