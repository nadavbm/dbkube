package postgres

import (
	"fmt"
	"strings"

	configmaps "github.com/nadavbm/dbkube/apis/configmaps/v1alpha1"
)

// GenerateHbaConf will generate config map data for postgres pg_hba.conf file
// to define connection setting to the db according to https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
func GenerateHbaConf(cm *configmaps.ConfigMap) string {
	config := ""
	for _, v := range cm.Spec.WhiteList {
		config = config + fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", v.AllowedHostType, v.AllowedDatabase, v.AllowedUser, v.AllowedAddress, v.AllowedIpMask, v.AllowedMethod)
	}

	return config
}

// GeneratePostgresqlConf will generate config map data for postgres postgresql.conf file
// postgresql.conf server settings example - https://github.com/postgres/postgres/blob/master/src/backend/utils/misc/postgresql.conf.sample
func GeneratePostgresqlConf(cm *configmaps.ConfigMap) string {
	config := fmt.Sprintf("data_directory = %s\n", cm.Spec.DataDirectory)
	config += fmt.Sprintf("hba_file = %s\n", cm.Spec.ConfigFileLocation)
	config += fmt.Sprintf("max_connections = %d\n", cm.Spec.MaxConnections)
	config += fmt.Sprintf("port = %d\n", cm.Spec.Port)
	cfg := strings.Replace(config, `\n`, "\n", -1)
	return cfg
}
