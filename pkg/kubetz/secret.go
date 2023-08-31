package kubetz

import (
	"math/rand"

	secrets "github.com/nadavbm/dbkube/apis/secrets/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$&*"

func BuildSecret(ns, app string, secret *secrets.Secret) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "db-secret",
			Namespace: ns,
		},
		StringData: createSecret(secret),
	}
}

//
// ------------------------------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------
//

func createSecret(secret *secrets.Secret) map[string]string {
	m := make(map[string]string)
	m["postgres_password"] = randStringBytes(12)
	m["postgres_database"] = secret.Spec.Database
	m["postgres_user"] = secret.Spec.User
	return m
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
