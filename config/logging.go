package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

/**
 * Init Logger
 */
func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stderr)
}
