package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

func NewConfigLoader(configURI string) Loader {
	return Loader{
		fileURI: configURI,
	}
}

type Loader struct {
	binary  []byte
	fileURI string
}

func (instance *Loader) ReadConfiguration(ptr interface{}) error {
	return doGetConfiguration(instance.fileURI, ptr)
}

func doGetConfiguration(fileURI string, ptr interface{}) error {
	binary, err := os.ReadFile(fileURI)
	if err != nil {
		return err
	}
	content, err := interpolate(string(binary))
	if err != nil {
		return err
	}
	if strings.HasSuffix(fileURI, ".json") {
		err = json.Unmarshal([]byte(content), ptr)
	} else {
		err = yaml.Unmarshal([]byte(content), ptr)
	}
	return err
}

func interpolate(s string) (string, error) {
	var err error
	re := regexp.MustCompile(`\$\{([^}]+)}`)
	result := re.ReplaceAllStringFunc(s, func(match string) string {
		if err != nil {
			return ""
		}
		content := match[2 : len(match)-1]
		parts := strings.Split(content, ":")
		// If there's no default value
		if len(parts) == 1 {
			envValue, exists := os.LookupEnv(parts[0])
			if !exists {
				err = errors.New("environment variable not found: " + parts[0])
				return ""
			}
			return envValue
		}
		// If there's a default value
		envValue, exists := os.LookupEnv(parts[0])
		if !exists {
			return parts[1]
		}
		return envValue
	})

	return result, err
}

func ToZapLogLevel(level string) (zapcore.Level, error) {
	if level == "" {
		return ToZapLogLevel("INFO")
	}
	switch level {
	case "INFO":
		return zap.InfoLevel, nil
	case "DEBUG":
		return zap.DebugLevel, nil
	case "ERROR":
		return zap.ErrorLevel, nil
	case "WARN":
		return zap.WarnLevel, nil
	default:
		return zap.PanicLevel, fmt.Errorf("logging.level=%s is unsupported", level)
	}
}
