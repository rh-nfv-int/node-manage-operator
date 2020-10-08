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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeLabelsSpec defines the desired state of NodeLabels
type NodeLabelsSpec struct {
	LabelGroup []LabelGroup   `json:"labelGroup,omitempty"`
	Mode       LabelGroupMode `json:"mode,omitempty"`
}

type LabelGroupMode string

const (
	LabelGroupModeMutuallyExclusive LabelGroupMode = "exclusive"
	LabelGroupModeApplyOnAll        LabelGroupMode = "all"
)

type LabelGroup struct {
	Count  int               `json:"count,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// NodeLabelsStatus defines the observed state of NodeLabels
type NodeLabelsStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NodeLabels is the Schema for the nodelabels API
type NodeLabels struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeLabelsSpec   `json:"spec,omitempty"`
	Status NodeLabelsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeLabelsList contains a list of NodeLabels
type NodeLabelsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeLabels `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeLabels{}, &NodeLabelsList{})
}
