package main

import (
	"context"
	"fmt"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("holiiiii")
	inCluster := false
	// Crea el cliente de Kubernetes

	var config *rest.Config
	var err error

	if !inCluster {
		config, err = clientcmd.BuildConfigFromFlags("", "/Users/lramos/.kube/config")
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err)
	}

	// Crea el cliente dinámico
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Define el GroupVersionResource para IngressRoute
	ingressRouteGVR := schema.GroupVersionResource{
		Group:    "traefik.containo.us",
		Version:  "v1alpha1",
		Resource: "ingressroutes",
	}

	applicationGVR := schema.GroupVersionResource{
		Group:    "hajimari.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}

	// Obtiene todos los IngressRoute de todos los namespaces
	ingressRoutes, err := dynamicClient.Resource(ingressRouteGVR).Namespace(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Para cada IngressRoute, crea un Application
	for _, ingressRoute := range ingressRoutes.Items {
		routes, found, err := unstructured.NestedSlice(ingressRoute.Object, "spec", "routes")
		if err != nil || !found {
			fmt.Printf("Error al obtener las rutas: %v\n", err)
			continue
		}

		// Para cada ruta, extrae la URL
		for _, route := range routes {
			routeMap, ok := route.(map[string]interface{})
			if !ok {
				fmt.Println("Error al convertir la ruta a map[string]interface{}")
				continue
			}

			match, ok := routeMap["match"].(string)
			if !ok {
				fmt.Println("Error al obtener el campo match")
				continue
			}

			// Usa una expresión regular para extraer la URL
			re := regexp.MustCompile(`Host\(` + "`" + `(.+?)` + "`" + `\)`)

			matchParts := re.FindStringSubmatch(match)
			if len(matchParts) < 2 {
				fmt.Println("Error al extraer la URL del campo match")
				continue
			}

			url := matchParts[1] // La URL está en la segunda parte del match
			fmt.Println(url)

			fmt.Println(ingressRoute.GetName()) // Imprime el nombre del IngressRoute

			// Crea un Application
			application := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "hajimari.io/v1alpha1",
					"kind":       "Application",
					"metadata": map[string]interface{}{
						"name": ingressRoute.GetName(), // Usa el mismo nombre que el IngressRoute
					},
					"spec": map[string]interface{}{
						"name": ingressRoute.GetName(), // Usa el mismo nombre que el IngressRoute
						"url":  "http://" + url,        // Usa el nombre del IngressRoute para construir la URL
					},
				},
			}
			fmt.Println(application)

			fmt.Println(applicationGVR)

			//Crea el recurso en el cluster
			// _, err = dynamicClient.Resource(applicationGVR).Namespace(ingressRoute.GetNamespace()).Create(context.TODO(), application, metav1.CreateOptions{})
			if err != nil {
				fmt.Printf("Error al crear Application: %v\n", err)
			}
		}
	}

}
