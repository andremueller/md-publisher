package file

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Exists returns true if the given file or directory exists
func Exists(fileName string) bool {
	_, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalf("Cannot detect if file %s exists - there is a singularity.", fileName)
	}
	return true
}
