/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/mod/semver"
)

// Downloads specified Go source tar or zip file in RUNGO_APP_DIR directory
func getGoSource(version, osys, dst string) error {
	ext := "tar.gz"
	if osys == "windows" {
		ext = "zip"
	}

	res, err := http.Get(fmt.Sprintf("%s/%s.%s",
		GO_URL,
		version,
		ext,
	))
	if err != nil {
		return fmt.Errorf("%w: %v", errRequestFailed, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", errUnexpectedStatus, res.Status)
	}

	out, err := os.Create(fmt.Sprintf("%s/%s.%s",
		dst,
		version,
		ext,
	))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}

	return nil
}

func getLatestGoVersion() (string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.Split(string(body), "\n")[0], nil
}

// Scrapes all Go versions that match the regexp and sorts them using
// semver package
func getGoVersions() ([]string, error) {
	// TODO: This should be cached
	res, err := http.Get(GO_URL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errRequestFailed, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", errUnexpectedStatus, err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	rawVersions := make([]string, 0)
	doc.Find(".toggleButton").Each(func(i int, s *goquery.Selection) {
		version := s.Find("span").Text()

		// Match versions from go1.16+ ahead and replace the leading "go"
		// prefix for a "v" prefix to sort it using the semver package
		r := regexp.MustCompile(`^go(\d+)\.(1[6-9]|[2-9]\d+)(?:\.(\d+))?$`)
		if r.MatchString(version) {
			rawVersions = append(rawVersions, strings.Replace(version, "go", "v", 1))
		}
	})

	rawVersions = slices.Compact(rawVersions)
	semver.Sort(rawVersions)
	slices.Reverse(rawVersions)

	versions := make([]string, 0)
	for _, rawVersion := range rawVersions {
		versions = append(versions, strings.Replace(rawVersion, "v", "go", 1))
	}

	return versions, nil
}
