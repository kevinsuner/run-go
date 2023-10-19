// SPDX-License-Identifier: MIT
package events

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func CreateTempAndRun(data []byte) (string, error) {
	filename := fmt.Sprintf("%d.go", time.Now().Unix())
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", err
	}

	out, _ := exec.Command("go", "run", filename).CombinedOutput()

	if err := os.Remove(filename); err != nil {
		return "", err
	}

	return string(out), nil
}
