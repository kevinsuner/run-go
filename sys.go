package main

import (
	"fmt"
	"os"
	"runtime"
)

func getOSAndArch() (string, string, error) {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "amd64" {
			return runtime.GOOS, runtime.GOARCH, nil
		} else if runtime.GOARCH == "arm64" {
			return runtime.GOOS, runtime.GOARCH, nil
		}

		return "", "", errUnsupportedArch
	case "linux":
		if runtime.GOARCH == "amd64" {
			return runtime.GOOS, runtime.GOARCH, nil
		} else if runtime.GOARCH == "arm64" {
			return runtime.GOOS, runtime.GOARCH, nil
		}

		return "", "", errUnsupportedArch
	case "windows":
		if runtime.GOARCH == "amd64" {
			return runtime.GOOS, runtime.GOARCH, nil
		}

		return "", "", errUnsupportedArch
	default:
		return "", "", errUnsupportedOS
	}
}

// Set the required env variables for the application to work
func setEnvVariables(home, version, osys, arch string) error {
	if err := os.Setenv(
		"RUNGO_APP_DIR", 
		fmt.Sprintf("%s/%s", home, APP_DIR),
	); err != nil {
		return err
	}

	if err := os.Setenv("RUNGO_GO_VER", version); err != nil {
		return err
	}

	if osys == "windows" {
		if err := os.Setenv(
			"RUNGO_GO_BIN",
			fmt.Sprintf("%s/%s/%s/%s.%s-%s/bin/go.exe",
				home,
				APP_DIR,
				GOS_DIR,
				version,
				osys,
				arch,
			),
		); err != nil {
			return err
		}

		return nil
	}

	if err := os.Setenv(
		"RUNGO_GO_BIN",
		fmt.Sprintf("%s/%s/%s/%s.%s-%s/bin/go",
			home,
			APP_DIR,
			GOS_DIR,
			version,
			osys,
			arch,
		),
	); err != nil {
		return err
	}

	return nil
}

// Update the version of the Go binary used in the application
func updateGoBinEnvVariable(dir, version, osys string) error {
	if osys == "windows" {
		if err := os.Setenv(
			"RUNGO_GO_BIN",
			fmt.Sprintf("%s/%s/bin/go.exe", dir, version),
		); err != nil {
			return err
		}

		return nil
	}

	if err := os.Setenv(
		"RUNGO_GO_BIN",
		fmt.Sprintf("%s/%s/bin/go", dir, version),
	); err != nil {
		return err
	}

	return nil
}
