// SPDX-License-Identifier: MIT
package events

import (
	"fmt"
	"os"
)

func LoadGoProject(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	projectDir := fmt.Sprintf("%s/%s/%s", home, SNIPPETS_DIR, name)
	_, err = os.ReadDir(projectDir)
	if os.IsNotExist(err) {
		return "", err
	}

	filename := fmt.Sprintf("%s/main.go", projectDir)
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
