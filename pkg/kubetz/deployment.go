package kubetz

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	deployments "github.com/nadavbm/dbkube/apis/deployments/v1alpha1"
)

const numOfReplicas = 1

// BuildDeployment creates a kubernetes deployment specification
func BuildDeployment(ns string, deploy *deployments.Deployment) *appsv1.Deployment {
	name := deploy.Spec.ImageName
	// TODO: num of replicas should change according to plan (1 for single instance, 2 for master slave config and 3 for cluster)
	// this plan require change in the apis and manifests
	replicas := int32(numOfReplicas)

	var app string

	switch {
	case strings.Contains(deploy.Spec.ImageName, "postgres"):
		app = "postgres"
	case strings.Contains(deploy.Spec.ImageName, "mysql"):
		app = "mysql"
	case strings.Contains(deploy.Spec.ImageName, "mongo"):
		app = "mongodb"
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: buildLabels(app, name),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: buildLabels(app, name),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           deploy.Spec.ImageName + ":" + deploy.Spec.ImageVersion,
							ImagePullPolicy: "IfNotPresent",
							Ports: []v1.ContainerPort{
								{
									Protocol:      v1.ProtocolTCP,
									ContainerPort: deploy.Spec.ContainerPort,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse(deploy.Spec.MemoryLimit),
									v1.ResourceCPU:    resource.MustParse(deploy.Spec.CpuLimit),
								},
								Requests: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse(deploy.Spec.MemoryRequest),
									v1.ResourceCPU:    resource.MustParse(deploy.Spec.CpuRequest),
								},
							},
							Env: []v1.EnvVar{
								getEnvVarSecretSource("POSTGRES_DATABASE", "db-secret", "postgres_database"),
								getEnvVarSecretSource("POSTGRES_USER", "db-secret", "postgres_user"),
								getEnvVarSecretSource("POSTGRES_PASSWORD", "db-secret", "postgres_password"),
							},
						},
					},
					RestartPolicy:                 v1.RestartPolicyAlways,
					DNSPolicy:                     v1.DNSClusterFirst,
					SchedulerName:                 v1.DefaultSchedulerName,
					TerminationGracePeriodSeconds: func(i int64) *int64 { return &i }(30),
				},
			},
		},
	}
}

//
// ------------------------------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------
//

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

func getEnvVarSecretSource(envName, name, secret string) v1.EnvVar {
	return v1.EnvVar{
		Name: envName,
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
				Key: secret,
			},
		},
	}
}

func getEnvVarConfigMapSource(configName, fileName string) v1.EnvVar {
	return v1.EnvVar{
		Name: configName,
		ValueFrom: &v1.EnvVarSource{
			ConfigMapKeyRef: &v1.ConfigMapKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configName,
				},
				Key: fileName,
			},
		},
	}
}
