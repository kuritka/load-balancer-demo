package env

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testEnvValueName = "EDT_TEST_ENV_VAL"

func TestGetStringFlagFromEnv(t *testing.T) {
	cases := []struct {
		name     string
		set      bool
		value    string
		expected string
		err      error
	}{
		{name: "Test flag #1", set: true, value: "test", expected: "test", err: nil},
		{name: "Test flag #2", set: true, value: "0123456789", expected: "0123456789", err: nil},
		{name: "Test flag #3", set: false, value: "0123456789", expected: "", err: fmt.Errorf("failed to get flag %s from env", testEnvValueName)},
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			if cases[i].set {
				os.Setenv(testEnvValueName, cases[i].value)
			}
			got, err := GetStringFlagFromEnv(testEnvValueName)
			assert.Equal(t, cases[i].err, err)
			assert.Equal(t, cases[i].expected, got)
			os.Unsetenv(testEnvValueName)
		})
	}
}

func TestMustGetStringFlagFromEnv(t *testing.T) {
	cases := []struct {
		name          string
		value         string
		set           bool
		expected      string
		expectedPanic bool
		err           error
	}{
		{name: "Test flag #1", set: true, value: "test", expected: "test", expectedPanic: false, err: nil},
		{name: "Test flag #2", set: true, value: "0123456789", expected: "0123456789", expectedPanic: false, err: nil},
		{name: "Test flag #3", set: false, value: "0123456789", expected: "", expectedPanic: true, err: errors.New("failed to load string flag from env")},
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					assert.True(t, cases[i].expectedPanic)
					assert.Equal(t, cases[i].err.Error(), r)
				}
			}()
			if cases[i].set {
				os.Setenv(testEnvValueName, cases[i].value)
			}
			got := MustGetStringFlagFromEnv(testEnvValueName)
			assert.Equal(t, cases[i].expected, got)
			os.Unsetenv(testEnvValueName)
		})
	}
}

func TestGetStringFlagFromEnvWithDefault(t *testing.T) {
	cases := []struct {
		name         string
		set          bool
		value        string
		defaultValue string
		expected     string
	}{
		{name: "Test flag #1", set: true, value: "test", defaultValue: "default", expected: "test"},
		{name: "Test flag #2", set: true, value: "0123456789", defaultValue: "default", expected: "0123456789"},
		{name: "Test flag #3", set: false, value: "0123456789", defaultValue: "default", expected: "default"},
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			if cases[i].set {
				os.Setenv(testEnvValueName, cases[i].value)
			}
			got := GetStringFlagFromEnvWithDefault(testEnvValueName, cases[i].defaultValue)
			assert.Equal(t, cases[i].expected, got)
			os.Unsetenv(testEnvValueName)
		})
	}
}
