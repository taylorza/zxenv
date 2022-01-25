package engine

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

const bufferSize = 32 * 1024

func downloadFile(url, destination string, progress func(totalBytes int, bytesRead <-chan int)) error {
	// Create the destination file
	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer f.Close()

	// Send HEAD request to get the file length
	headResp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer headResp.Body.Close()

	totalBytes, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		totalBytes = -1
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if progress != nil {
		update := make(chan int)

		// Kick-off status update function
		go progress(totalBytes, update)

		buffer := make([]byte, bufferSize)
		bytesRead := 0
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				f.Write(buffer[0:n])
				bytesRead += n
				update <- bytesRead
			}
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}
		close(update)
	} else {
		_, err := io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
