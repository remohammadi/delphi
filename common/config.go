package common

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

var defaults = map[string]string{
	"LOG_LEVEL":     "INFO",
	"SERVER_NAME":   "delphi-remohammadi.rhcloud.com",
	"TEMPLATES_DIR": "./templates",
	"STATIC_DIR":    "./static",
}

func ConfigString(name string) string {
	val := os.Getenv(fmt.Sprintf("%s", name))
	if val != "" {
		return val
	}
	return defaults[name]
}

func ConfigByteArray(name string) []byte {
	base64Value := ConfigString(name)
	value, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Warn("Error while decoding config value")
		return nil
	}
	return value
}
