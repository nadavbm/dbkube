package kubetz

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildObjectMetadata(name, namespace string, ownerRef []metav1.OwnerReference) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       namespace,
		OwnerReferences: ownerRef,
	}
}
