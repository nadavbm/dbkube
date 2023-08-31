package kubetz

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configmaps "github.com/nadavbm/dbkube/apis/configmaps/v1alpha1"
	"github.com/nadavbm/dbkube/pkg/postgres"
)

// BuildConfigMap will build a kubernetes config map for postgres
func BuildConfigMap(ns string, cm *configmaps.ConfigMap) *v1.ConfigMap {
	configMap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
	}
	switch {
	case cm.Spec.DatabaseServiceType == "postgres":
		configMap.ObjectMeta = metav1.ObjectMeta{
			Name:      "postgres-config",
			Namespace: ns,
		}
		configMap.Data = map[string]string{
			"pg_hba.conf":     postgres.GenerateHbaConf(cm),
			"postgresql.conf": postgres.GeneratePostgresqlConf(cm),
		}
	case cm.Spec.DatabaseServiceType == "mysql":
		configMap.ObjectMeta = metav1.ObjectMeta{
			Name:      "mysql-config",
			Namespace: ns,
		}
		configMap.Data = map[string]string{
			"my.cnf": postgres.GenerateHbaConf(cm),
		}
	case cm.Spec.DatabaseServiceType == "mongo":
		configMap.ObjectMeta = metav1.ObjectMeta{
			Name:      "mongodb-config",
			Namespace: ns,
		}
		configMap.Data = map[string]string{
			"mongod.conf": postgres.GenerateHbaConf(cm),
		}
	}

	return configMap
}
