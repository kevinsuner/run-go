/*
	SPDX-FileCopyrightText: 2023 Kevin Suñer <keware.dev@proton.me>
	SPDX-License-Identifier: MIT
*/
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func uncompressTarFile(file, dst string) error {
	reader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			_, err = os.Stat(target)
			if err != nil {
				err = os.MkdirAll(target, 0755)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg :
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(f, tarReader)
			if err != nil {
				return err
			}

			f.Close()
		}
	}

	err = os.Remove(file)
	if err != nil {
		return err
	}

	return nil
}

func uncompressZipFile(file, dst string) error {
	reader, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		target := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(target), 0755)
		if err != nil {
			return err
		}

		dstFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}

		dstFile.Close()
		srcFile.Close()
	}

	err = os.Remove(file)
	if err != nil {
		return err
	}

	return nil
}
