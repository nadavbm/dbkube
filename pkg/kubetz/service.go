package kubetz

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	services "github.com/nadavbm/dbkube/apis/services/v1alpha1"
)

// BuildService in kubernetes with pgDeploy port
func BuildService(ns, app string, service *services.Service) *v1.Service {
	name := "pg-service"
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: buildMetadata(name, ns, app, service.APIVersion, service.Kind, service.UID),
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Ports:    buildServicePorts(service),
			Selector: buildLabels(app, name),
		},
	}
}

//
// ------------------------------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------
//

func buildServicePorts(serivce *services.Service) []v1.ServicePort {
	var ports []v1.ServicePort

	for _, s := range serivce.Spec.Ports {
		sp := v1.ServicePort{
			Name:     s.Name,
			Protocol: v1.Protocol(s.Protocol),
			Port:     s.Port,
		}
		ports = append(ports, sp)
	}

	return ports
}
