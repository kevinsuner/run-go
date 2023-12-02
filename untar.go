/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func untarFile(file, dst string) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	gzipR, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzipR.Close()

	tarR := tar.NewReader(gzipR)
	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tarR); err != nil {
				return err
			}

			f.Close()
		}
	}

	return nil
}
