/*
    SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
    SPDX-License-Identifier: MIT
*/
package events

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const SNIPPETS_DIR = ".run-go/snippets"

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

func CreateGoProject(name string, data []byte) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("%s/%s", home, SNIPPETS_DIR)
	_, err = os.ReadDir(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	// TODO: If dir already exists throw error message
	projectDir := fmt.Sprintf("%s/%s", dir, name)
	if err := os.Mkdir(projectDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/main.go", projectDir)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
