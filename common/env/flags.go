// Package utils implements helpers.
package env

import (
	"fmt"
	"github.com/kuritka/onho.io/common/log"
	"os"
	"strconv"
)

// GetStringFlagFromEnv attemps to lookup a flag with name 'flagName' from the environment
func GetStringFlagFromEnv(flagName string) (string, error) {
	value, lookup := os.LookupEnv(flagName)
	if !lookup {
		return "", fmt.Errorf("failed to get flag %s from env", flagName)
	}

	return value, nil
}

// MustGetStringFlagFromEnv returns the value of the flag 'flagName' or panics if no value is set
func MustGetStringFlagFromEnv(flagName string) string {
	value, err := GetStringFlagFromEnv(flagName)
	if err != nil {
		log.Logger().Panic().Err(err).Msg("failed to load string flag from env")
	}
	return value
}

func MustGetIntFlagFromEnv(flagName string) int {
	strval, err := GetStringFlagFromEnv(flagName)
	if err == nil {
		value, err := strconv.Atoi(strval)
		if err == nil {
			return value
		}
	}
	log.Logger().Panic().Err(err).Msg("failed to load string flag from env")
	return -1
}

// EnvValueHolder holds values.
type EnvValueHolder struct {
	// Variable name to resolve (without prefix)
	N string
	// Pointer where to store the resolved value
	V *string
	// Default value if variable was not found or empty
	def *string
}

// NewEnvValueHolder creates new EnvValueHolder with name and pointer to the value
func NewEnvValueHolder(name string, value *string) *EnvValueHolder {
	return &EnvValueHolder{
		N: name,
		V: value,
	}
}

// Default sets the default value for a variable
func (h *EnvValueHolder) Default(value string) *EnvValueHolder {
	h.def = &value
	return h
}

func (h *EnvValueHolder) assign(value string) {
	// If variable is empty, set the default value (if we have it)
	if value == "" && h.hasDefault() {
		value = *h.def
	}

	*h.V = value
}

func (h *EnvValueHolder) hasDefault() bool {
	return h.def != nil
}

// MustGetStringFlagsFromEnv resolves provided env variables
func MustGetStringFlagsFromEnv(flagPrefix string, vars ...*EnvValueHolder) {
	for _, k := range vars {
		value, err := GetStringFlagFromEnv(flagPrefix + k.N)
		if err != nil && !k.hasDefault() {
			log.Logger().Panic().Err(err).Msg("failed to load string flag from env")
		}

		k.assign(value)
	}
}

// GetStringFlagFromEnvWithDefault returns the value of the flag 'flagName' or 'defaultValue' if no value is set
func GetStringFlagFromEnvWithDefault(flagName string, defaultValue string) string {
	value, err := GetStringFlagFromEnv(flagName)
	if err != nil {
		return defaultValue
	}
	return value
}
