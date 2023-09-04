package kubetz

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	services "github.com/nadavbm/dbkube/apis/services/v1alpha1"
)

// BuildService in kubernetes with pgDeploy port
func BuildService(ns string, svc *services.Service) *v1.Service {
	controller := true
	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: svc.APIVersion,
			Kind:       svc.Kind,
			Name:       svc.Name,
			UID:        svc.UID,
			Controller: &controller,
		},
	}

	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
	}
	switch {
	case svc.Spec.DatabaseServiceType == "postgres":
		service.ObjectMeta = buildObjectMetadata("postgres-service", ns, ownerRef)
		service.Spec = v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Ports:    buildServicePorts(svc),
			Selector: buildLabels("pg-service", "postgres"),
		}
	case svc.Spec.DatabaseServiceType == "mysql":
		service.ObjectMeta = buildObjectMetadata("mysql-service", ns, ownerRef)
		service.Spec = v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Ports:    buildServicePorts(svc),
			Selector: buildLabels("mysql-service", "mysql"),
		}
	case svc.Spec.DatabaseServiceType == "mongo":
		service.ObjectMeta = buildObjectMetadata("mongo-service", ns, ownerRef)
		service.Spec = v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Ports:    buildServicePorts(svc),
			Selector: buildLabels("mongodb-service", "mongodb"),
		}
	}

	return service
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
