package utils

import (
	"fmt"
	"os"

	// external
	log "github.com/sirupsen/logrus"
)

func ShowGolangEnvs() {
	names := []string{"GOGC", "GOTRACEBACK", "GOMAXPROCS", "GODEBUG"}
	values := make([]string, len(names))

	grid := 0
	for index, name := range names {
		values[index] = os.Getenv(name)
		if grid < len(name) {
			grid = len(name)
		}
	}

	f := fmt.Sprintf("\t%%-%ds:  %%s", grid+1)
	for index, name := range names {
		log.Infof(f, name, values[index])
	}
}
