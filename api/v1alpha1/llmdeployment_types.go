package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LLMDeploymentSpec defines the desired state of LLMDeployment
type LLMDeploymentSpec struct {
	// ModelRef is the name of the LLMModel to deploy
	ModelRef string `json:"modelRef"`

	// Replicas is the number of replicas to deploy
	Replicas int `json:"replicas"`

	// Port is the port the model service will listen on
	// +kubebuilder:default=8080
	Port int `json:"port,omitempty"`
}

// DeploymentCondition represents a condition of the deployment
type DeploymentCondition struct {
	// Type of deployment condition
	Type string `json:"type"`

	// Status of the condition, one of True, False, Unknown
	Status string `json:"status"`

	// LastTransitionTime is the last time the condition transitioned
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief reason for the condition's last transition
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human-readable explanation for the condition
	// +optional
	Message string `json:"message,omitempty"`
}

// LLMDeploymentStatus defines the observed state of LLMDeployment
type LLMDeploymentStatus struct {
	// AvailableReplicas is the number of available replicas
	AvailableReplicas int32 `json:"availableReplicas"`

	// Conditions represent the current service state
	Conditions []DeploymentCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LLMDeployment is the Schema for the llmdeployments API
type LLMDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LLMDeploymentSpec   `json:"spec,omitempty"`
	Status LLMDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LLMDeploymentList contains a list of LLMDeployment
type LLMDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LLMDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LLMDeployment{}, &LLMDeploymentList{})
}
