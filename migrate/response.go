package migrate

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Response struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata"`
		Spec              interface{} `json:"spec,omitempty"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}
