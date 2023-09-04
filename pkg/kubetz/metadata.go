package kubetz

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func buildObjectMetadata(name, namespace string, ownerRef []metav1.OwnerReference) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       namespace,
		OwnerReferences: ownerRef,
	}
}

func buildMetadata(name, namespace, app, apiVersion, kind string, uid types.UID) metav1.ObjectMeta {
	controlled := true
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    buildLabels(name, app),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: apiVersion,
				Kind:       kind,
				Name:       name,
				UID:        uid,
				Controller: &controlled,
			},
		},
	}
}

func buildLabels(name, app string) map[string]string {
	m := make(map[string]string)
	m["app"] = app
	m["app.kubernetes.io/name"] = name
	m["app.kubernetes.io/component"] = name
	return m
}
