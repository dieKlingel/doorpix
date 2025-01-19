package exec

import (
	"os/exec"
	"strings"
)

func Run(cmd []string) (string, error) {
	script := strings.Join(cmd, "\n")
	out, err := exec.Command("sh", "-c", script).CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
