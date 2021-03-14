package conf

import (
	"fmt"
	"os"
)

var AccessKey = ""
var ServerAddress = ":8090"
var DatabasePath = "./urlo.db"

func init() {
	AccessKey = requireEnv("URLO_ACCESS_KEY")
	ServerAddress = requireEnv("URLO_SERVER_ADDRESS")
	DatabasePath = requireEnv("URLO_DATABASE_PATH")
}

func requireEnv(name string) string {
	result := os.Getenv(name)
	if len(result) == 0 {
		panic(fmt.Sprintf("%s env var should be set", name))
	}
	return result
}
