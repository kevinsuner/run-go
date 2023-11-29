package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"time"
)

func runFromEditor(data []byte) (string, error) {
	f := fmt.Sprintf("%d.go", time.Now().Unix())
	if err := os.WriteFile(f, data, 0644); err != nil {
		return "", err
	}

	output, _ := exec.Command(os.Getenv("RUNGO_GOBIN"), "run", f).
					CombinedOutput()

	if err := os.Remove(f); err != nil {
		return "", err
	}

	return string(output), nil
}

func runFromSnippet(snippet string, data []byte) (string, error) {
	dir := fmt.Sprintf("%s/%s/%s/%s",
		os.Getenv("RUNGO_HOME"),
		APP_DIR,
		SNIPPETS_DIR,
		snippet,
	)

	f := fmt.Sprintf("%s/main.go", dir)
	if err := os.WriteFile(f, data, 0644); err != nil {
		return "", err
	}

	cmd := exec.Command(os.Getenv("RUNGO_GOBIN"), "mod", "tidy")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command(os.Getenv("RUNGO_GOBIN"), "run", "main.go")
	cmd.Dir = dir
	output, _ := cmd.CombinedOutput()

	return string(output), nil
}

func newSnippet(snippet string, data []byte) error {
	dir := fmt.Sprintf("%s/%s/%s/%s",
		os.Getenv("RUNGO_HOME"),
		APP_DIR,
		SNIPPETS_DIR,
		snippet,
	)

	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	cmd := exec.Command(os.Getenv("RUNGO_GOBIN"), "mod", "init", snippet)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	f := fmt.Sprintf("%s/main.go", dir)
	if err := os.WriteFile(f, data, 0644); err != nil {
		return err
	}

	return nil
}

func listSnippets() ([]string, error) {
	dir := fmt.Sprintf("%s/%s/%s",
		os.Getenv("RUNGO_HOME"),
		APP_DIR,
		SNIPPETS_DIR,
	)

	snippets := make([]string, 0)
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := os.Stat(fmt.Sprintf("%s/%s", dir, path))
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
	dir := fmt.Sprintf("%s/%s/%s/%s",
		os.Getenv("RUNGO_HOME"),
		APP_DIR,
		SNIPPETS_DIR,
		snippet,
	)

	_, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "", err
	}

	f := fmt.Sprintf("%s/main.go", dir)
	data, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
