package engine

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const bufferSize = 32 * 1024

type ProgressFunc func(bytesRead <-chan int)

func copyAll(dst io.Writer, src io.Reader, progress ProgressFunc) error {
	buffer := make([]byte, bufferSize)
	bytesRead := 0
	if progress != nil {
		update := make(chan int)
		defer close(update)

		go progress(update)
		for {
			n, err := src.Read(buffer)
			if n > 0 {
				dst.Write(buffer[0:n])
				bytesRead += n
				select {
				case update <- bytesRead:
				default:
				}
			}
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}
	} else {
		io.Copy(dst, src)
	}
	return nil
}

func downloadFile(url, destination string, progress ProgressFunc) error {
	// Create the destination file
	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = copyAll(f, resp.Body, progress)
	if err != nil {
		return err
	}

	return nil
}

func unzip(zipfile, destination string) error {
	return unzipAndStrip(zipfile, destination, false)
}

func unzipAndStrip(zipfile, destination string, stripFirstDir bool) error {
	err := os.MkdirAll(destination, 0777)
	if err != nil {
		return err
	}

	rdr, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rdr.Close()

	stripDir := ""
	for _, f := range rdr.File {
		source, err := f.Open()
		if err != nil {
			return err
		}
		if f.FileInfo().IsDir() && stripFirstDir && stripDir == "" {
			stripDir = f.Name
		}
		name := strings.Replace(f.Name, stripDir, "", 1)
		target := path.Join(destination, name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(target, f.Mode())
			if err != nil {
				return err
			}
		} else {
			newFile, err := os.Create(target)
			if err != nil {
				return err
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, source)
			if err != nil {
				return err
			}
			err = os.Chmod(target, f.Mode())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
