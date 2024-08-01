// models.go
package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressRouteSpec struct {
	// Define aquí los campos de tu IngressRoute
}

type IngressRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IngressRouteSpec `json:"spec,omitempty"`
}

type ApplicationSpec struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ApplicationSpec `json:"spec,omitempty"`
}
