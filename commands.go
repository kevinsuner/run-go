/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/

//go:build !windows
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Either run code from an existing snippet, or create a temporary .go file
// that gets executed and deleted
func runCode(snippet string, data []byte) (string, error) {
	if len(snippet) > 0 {
		dir := filepath.Join(os.Getenv("RUNGO_APP_DIR"), SNIPPETS_DIR, snippet)
		err := os.WriteFile(filepath.Join(dir, "main.go"), data, 0644)
		if err != nil {
			return "", err
		}

		cmd := exec.Command(os.Getenv("RUNGO_GO_BIN"), "mod", "tidy")
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			return "", err
		}

		cmd = exec.Command(os.Getenv("RUNGO_GO_BIN"), "run", "main.go")
		cmd.Dir = dir
		output, _ := cmd.CombinedOutput()

		return string(output), nil
	}

	file := filepath.Join(os.Getenv("RUNGO_APP_DIR"), fmt.Sprintf("%d.go", time.Now().Unix()))
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		return "", err
	}

	output, _ := exec.Command(os.Getenv("RUNGO_GO_BIN"), "run", file).CombinedOutput()

	err = os.Remove(file)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func newSnippet(snippet string, data []byte) error {
	dir := filepath.Join(os.Getenv("RUNGO_APP_DIR"), SNIPPETS_DIR, snippet)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command(os.Getenv("RUNGO_GO_BIN"), "mod", "init", snippet)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dir, "main.go"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

