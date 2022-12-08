/*
Copyright 2022.

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
	"fmt"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// TemplateMetadata is a metadata type only for labels and annotations.
type TemplateMetadata struct {
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type HPATemplate struct {
	Metadata *TemplateMetadata                              `json:"metadata,omitempty"`
	Spec     autoscalingv2beta2.HorizontalPodAutoscalerSpec `json:"spec"`
}

type Cron struct {
	Name        string `json:"name"`
	Schedule    string `json:"schedule"`
	Timezone    string `json:"timezone"`
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
}

// CronHPASpec defines the desired state of CronHPA
type CronHPASpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Template HPATemplate `json:"template"`
	Cron     []Cron      `json:"cron"`
}

// CronHPAStatus defines the observed state of CronHPA
type CronHPAStatus struct {
	// LastCronTimestamp is the time of last cron job.
	LastCronTimestamp *metav1.Time `json:"lastCronTimestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cronhpa
// +kubebuilder:printcolumn:name="MINPODS",type="string",JSONPath=".spec.template.spec.minReplicas",description="MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down"
// +kubebuilder:printcolumn:name="MINPODS",type="string",JSONPath=".spec.template.spec.minReplicas",description="MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down"

// CronHPA is the Schema for the cronhpas API
type CronHPA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronHPASpec   `json:"spec,omitempty"`
	Status CronHPAStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CronHPAList contains a list of CronHPA
type CronHPAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronHPA `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CronHPA{}, &CronHPAList{})
}

func (h *CronHPA) String() string {
	fmt.Printf("%+v", h.Spec.Cron)
	return ("true")
}
