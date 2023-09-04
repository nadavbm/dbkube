package kubetz

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configmaps "github.com/nadavbm/dbkube/apis/configmaps/v1alpha1"
	"github.com/nadavbm/dbkube/pkg/mongodb"
	"github.com/nadavbm/dbkube/pkg/mysql"
	"github.com/nadavbm/dbkube/pkg/postgres"
)

// BuildConfigMap will build a kubernetes config map for postgres
func BuildConfigMap(ns string, cm *configmaps.ConfigMap) *v1.ConfigMap {
	controller := true
	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: cm.APIVersion,
			Kind:       cm.Kind,
			Name:       cm.Name,
			UID:        cm.UID,
			Controller: &controller,
		},
	}
	configMap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
	}

	switch {
	case cm.Spec.DatabaseServiceType == "postgres":
		configMap.ObjectMeta = buildObjectMetadata("postgres-config", ns, ownerRef)
		configMap.Data = map[string]string{
			"pg_hba.conf":     postgres.GenerateHbaConf(cm),
			"postgresql.conf": postgres.GeneratePostgresqlConf(cm),
		}
	case cm.Spec.DatabaseServiceType == "mysql":
		configMap.ObjectMeta = buildObjectMetadata("mysql-config", ns, ownerRef)

		configMap.Data = map[string]string{
			"my.cnf": mysql.GenerateMyCnfFile(cm),
		}
	case cm.Spec.DatabaseServiceType == "mongo":
		configMap.ObjectMeta = buildObjectMetadata("mongo-config", ns, ownerRef)

		configMap.Data = map[string]string{
			"mongod.conf": mongodb.GenerateMongodFile(cm),
		}
	}

	return configMap
}
