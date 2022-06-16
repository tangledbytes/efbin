package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/utkarsh-pro/efbin/pkg/constants"
)

// Exit takes a message and exits the program with an error code.
func Exit(status int, msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(status)
}

// GetEnvOrDefault returns the value of the environment variable named by the key
// and if not found, returns the default value.
func GetEnvOrDefault(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return val
}

// IsWrappedBinPresent takes the name of a binary and checks in the $PATH
// if the binary exists in the $PATH or not
func IsWrappedBinPresent(bin string) bool {
	if bin == "" {
		return false
	}

	_, err := exec.LookPath(bin)
	return err == nil
}

// GetBinaryName returns the name of the binary which efbin needs to wrap
func GetBinaryName() string {
	return os.Getenv(fmt.Sprintf("%s_BIN", constants.EnvPrefix))
}

// PreventFuckUp ensures that the binary doesn't attempt to execute itself
func PreventFuckUp() {
	bin, err := os.Executable()
	if err != nil {
		Exit(1, "failed to get current executable: ", err)
	}

	if bin == GetBinaryName() {
		Exit(
			1,
			"failed to execute binary: binary is the same as the executable - consider giving fullpath to the binary in the environment variable: ",
			fmt.Sprintf("%s_BIN", constants.EnvPrefix),
		)
	}
}
