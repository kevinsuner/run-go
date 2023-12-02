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

func getGoTarball(version, dst string) error {
	res, err := http.Get(fmt.Sprintf("%s/%s.tar.gz", GO_URL, version))
	if err != nil {
		return fmt.Errorf("%w: %v", errRequestFailed, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", errUnexpectedStatus, err)
	}

	out, err := os.Create(fmt.Sprintf("%s/%s.tar.gz", dst, version))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}

	return nil
}

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
