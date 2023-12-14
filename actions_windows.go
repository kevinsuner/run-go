/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/

//go:build windows
package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"syscall"
)

func runFromEditor(data []byte) (string, error) {
	file := filepath.Join(os.Getenv("RUNGO_APP_DIR"), fmt.Sprintf("%d.go", time.Now().Unix()))
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(os.Getenv("RUNGO_GO_BIN"), "run", file)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := cmd.CombinedOutput()

	err = os.Remove(file)
	if err != nil {
		return "", err
	}
	
	return string(output), nil
}

func runFromSnippet(snippet string, data []byte) (string, error) {
	dir := filepath.Join(os.Getenv("RUNGO_APP_DIR"), SNIPPETS_DIR, snippet)
	err := os.WriteFile(filepath.Join(dir, "main.go"), data, 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(os.Getenv("RUNGO_GO_BIN"), "mod", "tidy")
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	cmd = exec.Command(os.Getenv("RUNGO_GO_BIN"), "run", "main.go")
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := cmd.CombinedOutput()

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
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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

func listSnippets() ([]string, error) {
	dir := filepath.Join(os.Getenv("RUNGO_APP_DIR"), SNIPPETS_DIR) 
	snippets := make([]string, 0)
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := os.Stat(filepath.Join(dir, path))
		if err != nil {
			return err
		}

		if info.IsDir() {
			snippets = append(snippets, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return snippets[1:], nil
}

func openSnippet(snippet string) (string, error) {
	dir := filepath.Join(os.Getenv("RUNGO_APP_DIR"), SNIPPETS_DIR, snippet)
	_, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "", err
	}

	data, err := os.ReadFile(filepath.Join(dir, "main.go"))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
