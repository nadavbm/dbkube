package kubetz

import (
	"math/rand"
	"strings"

	secrets "github.com/nadavbm/dbkube/apis/secrets/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$&*"

func BuildSecret(ns string, s *secrets.Secret) *v1.Secret {
	controller := true
	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: s.APIVersion,
			Kind:       s.Kind,
			Name:       s.Name,
			UID:        s.UID,
			Controller: &controller,
		},
	}

	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		}}
	switch {
	case s.Spec.DatabaseServiceType == "postgres":
		secret.ObjectMeta = buildObjectMetadata("postgres-secret", ns, ownerRef)
		secret.StringData = createPostgresSecret(s)
	case s.Spec.DatabaseServiceType == "mysql":
		secret.ObjectMeta = buildObjectMetadata("mysql-secret", ns, ownerRef)
		secret.StringData = createMySQLSecret(s)
	case s.Spec.DatabaseServiceType == "mongo":
		secret.ObjectMeta = buildObjectMetadata("mongo-secret", ns, ownerRef)
		secret.StringData = createMongoDBSecret(s)
	}

	return secret
}

//
// --------------------------------------------------------------------------------- deployment helpers -----------------------------------------------------------------------------
//

func getEnvVarsForSecret(imageName string) []v1.EnvVar {
	switch {
	case strings.Contains(imageName, "postgres"):
		return []v1.EnvVar{
			getEnvVarSecretSource("POSTGRES_DATABASE", "postgres-secret", "postgres_database"),
			getEnvVarSecretSource("POSTGRES_USER", "postgres-secret", "postgres_user"),
			getEnvVarSecretSource("POSTGRES_PASSWORD", "postgres-secret", "postgres_password"),
		}
	case strings.Contains(imageName, "mysql"):
		return []v1.EnvVar{
			getEnvVarSecretSource("MYSQL_DATABASE", "mysql-secret", "mysql_database"),
			getEnvVarSecretSource("MYSQL_USER", "mysql-secret", "mysql_user"),
			getEnvVarSecretSource("MYSQL_PASSWORD", "mysql-secret", "mysql_password"),
			getEnvVarSecretSource("MYSQL_ROOT_PASSWORD", "mysql-secret", "mysql_root_password"),
		}
	case strings.Contains(imageName, "mongo"):
		return []v1.EnvVar{
			getEnvVarSecretSource("MONGO_USER", "mongo-secret", "mongo_user"),
			getEnvVarSecretSource("MONGO_PASSWORD", "mongo-secret", "mongo_password"),
		}
	}
	return []v1.EnvVar{}
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

//
// ------------------------------------------------------------------------------ secret generator helpers -----------------------------------------------------------------------------
//

func buildSecretObjectMetadata(name, namespace string, s *secrets.Secret) metav1.ObjectMeta {
	controller := true
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: s.APIVersion,
				Kind:       s.Kind,
				Name:       name,
				UID:        s.UID,
				Controller: &controller,
			},
		},
	}
}

func createPostgresSecret(secret *secrets.Secret) map[string]string {
	m := make(map[string]string)
	m["postgres_password"] = randStringBytes(12)
	m["postgres_database"] = secret.Spec.Database
	m["postgres_user"] = secret.Spec.User
	return m
}

func createMySQLSecret(secret *secrets.Secret) map[string]string {
	m := make(map[string]string)
	m["mysql_password"] = randStringBytes(12)
	m["mysql_root_password"] = randStringBytes(12)
	m["mysql_database"] = secret.Spec.Database
	m["mysql_user"] = secret.Spec.User
	return m
}

func createMongoDBSecret(secret *secrets.Secret) map[string]string {
	m := make(map[string]string)
	m["mongo_password"] = randStringBytes(12)
	m["mongo_user"] = secret.Spec.User
	return m
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
