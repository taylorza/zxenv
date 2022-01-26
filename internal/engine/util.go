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
	return unzipAndStrip(zipfile, destination, "")
}

func unzipAndStrip(zipfile, destination, stripPath string) error {
	rdr, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer rdr.Close()

	for _, f := range rdr.File {
		source, err := f.Open()
		if err != nil {
			return err
		}
		fname := strings.Replace(f.Name, stripPath, "", 1)
		target := path.Join(destination, fname)
		if f.FileInfo().IsDir() {
			os.MkdirAll(target, f.Mode())
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
		}
	}
	return nil
}
