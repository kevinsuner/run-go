// SPDX-License-Identifier: MIT
package events

import (
	"fmt"
	"io/fs"
	"os"
)

func ListGoProjects() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	projectDir := fmt.Sprintf("%s/%s", home, SNIPPETS_DIR)
	_, err = os.ReadDir(projectDir)
	if os.IsNotExist(err) {
		return nil, err
	}

	projects := make([]string, 0)
	if err := fs.WalkDir(
		os.DirFS(projectDir),
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			fileInfo, err := os.Stat(fmt.Sprintf("%s/%s", projectDir, path))
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				projects = append(projects, path)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return projects[1:], nil
}
