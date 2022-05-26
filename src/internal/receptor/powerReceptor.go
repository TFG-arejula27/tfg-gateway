package receptor

import (
	"errors"
	"os/exec"
	"strings"
)

//Run executes powerstate
func runPowerstat() (string, error) {
	cmd := exec.Command("powerstat", "-R", "-n", "1", "60")
	data, err := cmd.Output()
	if err != nil {
		return "", err

	}
	output := string(data)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")

		if strings.HasPrefix(line, "Average") {
			parsedLine := strings.Split(line, " ")
			return parsedLine[12], nil
		} else {
			continue
		}

	}

	return "", errors.New("cant gather power")

}
