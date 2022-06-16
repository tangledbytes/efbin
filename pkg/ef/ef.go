package ef

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/utkarsh-pro/efbin/pkg/constants"
	"github.com/utkarsh-pro/efbin/pkg/util"
)

// Run runs the wrapped binary with the given arguments.
func Run(args []string) error {
	transArgs, err := TransformArgsWithEnv(args)
	if err != nil {
		return fmt.Errorf("failed to execute binary: %w", err)
	}

	cmd := exec.Command(util.GetBinaryName(), transArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// TransformArgsWithSet transforms a set of arguments with a set of flags.
func TransformArgsWithEnv(args []string) ([]string, error) {
	transformed, err := TransformArgsWithSet(args, ConvertEnvToFlags(), os.Getenv("UDOCKER__TARGETARG"))
	if err != nil {
		return nil, fmt.Errorf("failed to transform args with env: %w", err)
	}

	return transformed, nil
}

// ConvertEnvToFlags converts a set of environment variables to a set of flags.
//
// The function skips the environment variables that start with the prefix defined in constants.EnvPrefix + "__"
func ConvertEnvToFlags() []string {
	envs := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, constants.EnvPrefix) && !strings.HasPrefix(env, constants.EnvPrefix+"_") {
			envs = append(envs, strings.TrimPrefix(env, constants.EnvPrefix))
		}
	}

	return ConvertStringSetToFlags(envs)
}

// ConvertStringSetToFlags converts a set of strings to a set of flags.
func ConvertStringSetToFlags(set []string) []string {
	flags := []string{}
	for _, env := range set {
		parsed := strings.SplitN(env, "=", 2)
		if len(parsed) != 2 {
			continue
		}

		key := parsed[0]
		value := parsed[1]

		if len(key) == 1 {
			flags = append(flags, "-"+key)
		} else {
			flags = append(flags, "--"+strings.ReplaceAll(key, "_", "-"))
		}

		if len(value) > 0 {
			flags = append(flags, value)
		}
	}

	return flags
}

// TransformArgsWithSet transforms a set of arguments with a set of flags.
// The target argument is the name of the build to transform => target can be of form <arg> or <arg>:<skips>
func TransformArgsWithSet(args, set []string, target string) ([]string, error) {
	newArgs := []string{}
	argName := ""
	argSkips := 0
	targetArg := target

	if targetArg == "" {
		newArgs = append(newArgs, set...)
		newArgs = append(newArgs, args...)
		return newArgs, nil
	}

	parsed := strings.Split(targetArg, ":")
	if len(parsed) == 1 {
		argName = parsed[0]
	}
	if len(parsed) == 2 {
		argName = parsed[0]
		parsedPos, err := strconv.ParseInt(parsed[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid arg skips: %s", parsed[1])
		}

		argSkips = int(parsedPos)
	}
	if len(parsed) > 2 {
		return nil, fmt.Errorf("invalid target arg: %s", parsed[1])
	}

	skips := 0
	for i, arg := range args {
		if arg == argName {
			if skips == argSkips {
				newArgs = append(newArgs, args[:i+1]...)
				newArgs = append(newArgs, set...)
				newArgs = append(newArgs, args[i+1:]...)

				return newArgs, nil
			}

			skips += 1
		}
	}

	return args, nil
}
