package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LLMModelSpec defines the desired state of LLMModel
type LLMModelSpec struct {
	// ModelName is the name of the LLM model
	ModelName string `json:"modelName"`

	// Image is the container image for the model
	Image string `json:"image"`

	// Resources defines the resource requirements for the model
	Resources ModelResources `json:"resources,omitempty"`
}

// ModelResources defines resource requirements
type ModelResources struct {
	// CPU resource requirements
	CPU string `json:"cpu,omitempty"`

	// Memory resource requirements
	Memory string `json:"memory,omitempty"`
}

// LLMModelStatus defines the observed state of LLMModel
type LLMModelStatus struct {
	// Phase indicates the current phase of the model (Pending, Ready, Failed)
	Phase string `json:"phase,omitempty"`

	// Message provides additional details about the current phase
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LLMModel is the Schema for the llmmodels API
type LLMModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LLMModelSpec   `json:"spec,omitempty"`
	Status LLMModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LLMModelList contains a list of LLMModel
type LLMModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LLMModel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LLMModel{}, &LLMModelList{})
}
