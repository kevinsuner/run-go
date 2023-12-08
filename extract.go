package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFile(file, dst string) error {
	switch {
	case strings.HasSuffix(file, ".tar.gz"):
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
	case strings.HasSuffix(file, ".zip"):
		r, err := zip.OpenReader(file)
		if err != nil {
			return err
		}
		defer r.Close()

		for _, f := range r.File {
			target := filepath.Join(dst, f.Name)
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
				continue
			}

			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			dstFile, err := os.OpenFile(
				target, 
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 
				f.Mode(),
			)
			if err != nil {
				return err
			}

			srcFile, err := f.Open()
			if err != nil {
				return err
			}

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return err
			}

			dstFile.Close()
			srcFile.Close()
		} 
	default:
		return errUnknownFileExt
	}

	return nil
}

