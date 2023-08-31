package mongodb

import (
	"fmt"
	"strings"

	configmaps "github.com/nadavbm/dbkube/apis/configmaps/v1alpha1"
)

// GenerateMongodFile generates config map data for mongo
// ref http://docs.mongodb.org/manual/reference/configuration-options/
// ref https://www.mongodb.com/docs/manual/reference/configuration-options/
func GenerateMongodFile(cm *configmaps.ConfigMap) string {
	config := fmt.Sprintf("systemLog:\n")
	config += fmt.Sprintf("  destination: file\n")
	config += fmt.Sprintf("  logAppend: true\n")
	config += fmt.Sprintf("  path: /var/log/mongodb/mongod.log\n")
	config += fmt.Sprintf("storage:\n")
	config += fmt.Sprintf("  dbPath: %s\n", cm.Spec.DataDirectory)
	config += fmt.Sprintf("  journal:\n")
	config += fmt.Sprintf("    enabled: true\n")
	config += fmt.Sprintf(`  engine: "wiredTiger"\n`)
	config += fmt.Sprintf("  wiredTiger:\n")
	config += fmt.Sprintf("    engineConfig: true\n")
	config += fmt.Sprintf("      cacheSizeGB: 1\n")
	config += fmt.Sprintf("processManagement:\n")
	config += fmt.Sprintf("  fork: true\n")
	config += fmt.Sprintf("  pidFilePath: /var/run/mongodb/mongod.pid\n")
	config += fmt.Sprintf("net:\n")
	config += fmt.Sprintf("  port: %d\n", cm.Spec.Port)

	cfg := strings.Replace(config, `\n`, "\n", -1)
	return cfg
}
