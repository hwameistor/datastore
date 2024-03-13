package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	DataLoadingPhasePreparing = "Preparing"
	DataLoadingPhaseLoading   = "Loading"
	DataLoadingPhaseCompleted = "Completed"
	DataLoadingPhaseFailed    = "Failed"
)

// DataLoadRequestSpec defines the desired state of DataLoadRequest
type DataLoadRequestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// indicate if the request is for all or not
	IsGlobal bool `json:"isGlobal"`

	// name of the node who will loads the data, and it works only when isglobal is false
	Node string `json:"node"`

	// name of the data source
	DataSource string `json:"dataSource"`
	SubDir     string `json:"subDir,omitempty"`

	// +kubebuilder:default:=false
	Retry bool `json:"retry"`
	// +kubebuilder:default:=1
	RetryInterval int64 `json:"retryInterval"`
}

// DataLoadRequestStatus defines the observed state of DataLoadRequest
type DataLoadRequestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:Enum:=Preparing;Loading;Completed;Failed
	Phase string `json:"phase,omitempty"`

	// +kubebuilder:default:=0
	RetryCount int64 `json:"retryCount"`

	LoadingStartTime metav1.Time `json:"loadingStartTime,omitempty"`

	LoadingCompleteTime metav1.Time `json:"loadingCompleteTime,omitempty"`

	Error string `json:"error,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataLoadRequest is the Schema for the dataloadrequests API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=dataloadrequests,scope=Namespaced,shortName=dlr
// +kubebuilder:printcolumn:JSONPath=".spec.dataSource",name=DataSource,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.subDir",name=SubDir,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.retry",name=Retry,type=boolean
// +kubebuilder:printcolumn:JSONPath=".spec.retryInterval",name=RetryInterval,type=integer,priority=1
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Phase,type=string
// +kubebuilder:printcolumn:JSONPath=".status.retryCount",name=RetryCount,type=integer,priority=1
// +kubebuilder:printcolumn:JSONPath=".status.loadingStartTime",name=LoadingStartTime,type=date
// +kubebuilder:printcolumn:JSONPath=".status.loadingCompleteTime",name=LoadingCompleteTime,type=date
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date
// +kubebuilder:printcolumn:JSONPath=".status.error",name=Error,type=string
type DataLoadRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataLoadRequestSpec   `json:"spec,omitempty"`
	Status DataLoadRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataLoadRequestList contains a list of DataLoadRequest
type DataLoadRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataLoadRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataLoadRequest{}, &DataLoadRequestList{})
}
