package mysql

import (
	"fmt"
	"strings"

	configmaps "github.com/nadavbm/dbkube/apis/configmaps/v1alpha1"
)

// GenerateMyCnfFile generates the data for mysql config map
// example here: https://www.ibm.com/docs/en/ztpf/2022?topic=performance-mysql-configuration-file-example
func GenerateMyCnfFile(cm *configmaps.ConfigMap) string {
	config := fmt.Sprintf("[mysqld]\n")
	config += fmt.Sprintf("port = %d\n", cm.Spec.Port)
	config += fmt.Sprintf("socket = /tmp/mysql.sock\n")
	config += fmt.Sprintf("max_connections = %d\n", cm.Spec.MaxConnections)
	cfg := strings.Replace(config, `\n`, "\n", -1)
	return cfg
}
