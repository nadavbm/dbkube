package kubetz

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	deployments "github.com/nadavbm/dbkube/apis/deployments/v1alpha1"
)

const numOfReplicas = 1

// BuildDeployment creates a kubernetes deployment specification
func BuildDeployment(ns string, deploy *deployments.Deployment) *appsv1.Deployment {
	controller := true
	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: deploy.APIVersion,
			Kind:       deploy.Kind,
			Name:       deploy.Name,
			UID:        deploy.UID,
			Controller: &controller,
		},
	}

	name := deploy.Spec.ImageName + "-deploy"
	// TODO: num of replicas should change according to plan (1 for single instance, 2 for master slave config and 3 for cluster)
	// this plan require change in the apis and manifests
	replicas := int32(numOfReplicas)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
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
							Env: getEnvVarsForSecret(deploy.Spec.ImageName),
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

	switch {
	case strings.Contains(deploy.Spec.ImageName, "postgres"):
		deployment.ObjectMeta = buildObjectMetadata("postgres-deploy", ns, ownerRef)
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: buildLabels("postgres-deploy", "postgres"),
		}
		deployment.Spec.Template.ObjectMeta = metav1.ObjectMeta{
			Labels: buildLabels("postgres-deploy", "postgres"),
		}
	case strings.Contains(deploy.Spec.ImageName, "mysql"):
		deployment.ObjectMeta = buildObjectMetadata("mysql-deploy", ns, ownerRef)
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: buildLabels("mysql-deploy", "mysql"),
		}
		deployment.Spec.Template.ObjectMeta = metav1.ObjectMeta{
			Labels: buildLabels("mysql-deploy", "mysql"),
		}
	case strings.Contains(deploy.Spec.ImageName, "mongo"):
		deployment.ObjectMeta = buildObjectMetadata("mongo-deploy", ns, ownerRef)
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: buildLabels("mongo-deploy", "mongo"),
		}
		deployment.Spec.Template.ObjectMeta = metav1.ObjectMeta{
			Labels: buildLabels("mongo-deploy", "mongo"),
		}
	}

	return deployment
}

//
// ------------------------------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------
//

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
