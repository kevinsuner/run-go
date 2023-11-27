/*
    SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
    SPDX-License-Identifier: MIT
*/
package events

import (
	"fmt"
	"os"
	"os/exec"
)

func RunGoProject(name string, data []byte) (string, error) {
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
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", err
	}

	cmd := exec.Command(os.Getenv("GOPATH"), "mod", "tidy")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command(os.Getenv("GOPATH"), "run", "main.go")
	cmd.Dir = projectDir
	out, _ := cmd.CombinedOutput()

	return string(out), nil
}
