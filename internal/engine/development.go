package engine

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

var urls = map[string]string{
	"cspect":     "http://www.javalemmings.com/public/zxnext/CSpect2_15_01.zip",
	"sdcard":     "http://www.zxspectrumnext.online/cspect/tbbluemmc-512mb.zip",
	"openal":     "http://www.zxspectrumnext.online/cspect/openal.zip",
	"asm":        "https://github.com/z00m128/sjasmplus/releases/download/v1.18.3/sjasmplus-1.18.3.win.zip",
	"hdfmonkey":  "http://uto.speccy.org/downloads/hdfmonkey_windows.zip",
	"dezog":      "https://github.com/maziac/DeZogPlugin/releases/download/v2.1.0/DeZogPlugin.dll",
	"dezog-conf": "https://raw.githubusercontent.com/maziac/DeZogPlugin/main/DeZogPlugin.dll.config",
}

func SetupDevelopment(env *Environment) error {
	err := makeDirectories(env)
	if err != nil {
		return err
	}

	err = installEmulator(env)
	if err != nil {
		return err
	}

	err = installDeZog(env)
	if err != nil {
		return err
	}

	err = setupSdCard(env)
	if err != nil {
		return err
	}

	err = installAssembler(env)
	if err != nil {
		return err
	}

	err = installSDTool(env)
	if err != nil {
		return err
	}
	return nil
}

func makeDirectories(env *Environment) error {
	err := os.MkdirAll(env.TempPath(), os.ModeDir)
	if err != nil {
		return fmt.Errorf("failed for create temporary directory (%w)", err)
	}

	err = os.MkdirAll(env.EmulatorPath(), os.ModeDir)
	if err != nil {
		return fmt.Errorf("failed for create emulator directory (%w)", err)
	}

	err = os.MkdirAll(env.SDPath(), os.ModeDir)
	if err != nil {
		return fmt.Errorf("failed for create sd card directory (%w)", err)
	}

	return nil
}

func installEmulator(env *Environment) error {
	zippath := path.Join(env.TempPath(), "cspect.zip")
	log.Println("Downloading emulator")
	err := download("cspect", zippath)
	if err != nil {
		return fmt.Errorf("failed to install emulator (%w)", err)
	}

	log.Println("Unzipping emulator")
	err = unzip(zippath, env.EmulatorPath())
	if err != nil {
		return fmt.Errorf("failed to unzip the emulator (%w)", err)
	}
	return nil
}

func installDeZog(env *Environment) error {
	targetPath := path.Join(env.EmulatorPath(), "DeZogPlugin.dll")

	log.Println("Downloading DeZogPlugin")
	err := download("dezog", targetPath)
	if err != nil {
		return fmt.Errorf("failed to download DeZogPlugin (%w)", err)
	}

	targetPath = path.Join(env.EmulatorPath(), "DeZogPlugin.dll.config")
	log.Println("Downloading DeZogPlugin config")
	err = download("dezog-conf", targetPath)
	if err != nil {
		return fmt.Errorf("failed to download DeZogPlugin config (%w)", err)
	}

	return nil
}

func installSDTool(env *Environment) error {
	targetPath := path.Join(env.EmulatorPath(), "hdfmonkey.exe")

	log.Println("Downloading hdfmonkey")
	err := download("hdfmonkey", targetPath)
	if err != nil {
		return fmt.Errorf("failed to download hdfmonkey (%w)", err)
	}

	return nil
}

func setupSdCard(env *Environment) error {
	zippath := path.Join(env.TempPath(), "sdcard.zip")

	log.Println("Downloading SD Card")
	err := download("sdcard", zippath)
	if err != nil {
		return fmt.Errorf("failed to download sd card image (%w)", err)
	}

	log.Println("Unzipping SD Card")
	err = unzip(zippath, env.TempPath())
	if err != nil {
		return fmt.Errorf("failed to unzip the sd card image (%w)", err)
	}

	files, err := os.ReadDir(env.TempPath())
	if err != nil {
		return err
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".img" || path.Ext(file.Name()) == ".mmc" {
			log.Printf("Copying %v\n", file.Name())
			_, err = copyFile(path.Join(env.TempPath(), file.Name()), path.Join(env.SDPath(), file.Name()))
			if err != nil {
				return err
			}
		}

		if path.Ext(file.Name()) == ".rom" {
			log.Printf("Copying %v\n", file.Name())
			_, err = copyFile(path.Join(env.TempPath(), file.Name()), path.Join(env.EmulatorPath(), file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func installAssembler(env *Environment) error {
	zippath := path.Join(env.TempPath(), "sjasmplus.zip")
	unzippath := path.Join(env.TempPath(), "sjasmplus")

	log.Printf("Downloading sjasmplus")
	err := download("asm", zippath)
	if err != nil {
		return fmt.Errorf("failed to download sd card image (%w)", err)
	}

	log.Printf("Unzipping sjasmplus")
	err = unzip(zippath, unzippath)
	if err != nil {
		return fmt.Errorf("failed to unzip the sd card image (%w)", err)
	}

	err = filepath.WalkDir(unzippath, func(fullpath string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == "sjasmplus.exe" {
			log.Printf("Copying %v\n", d.Name())
			_, err2 := copyFile(fullpath, path.Join(env.EmulatorPath(), d.Name()))
			if err2 != nil {
				return err2
			}
			return io.EOF
		}
		return nil
	})

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func download(urlkey, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	downloadFile(urls[urlkey], filename, func(totalBytes int, progress <-chan int) {
		for bytes := range progress {
			fmt.Printf("[%v / %v]\x1b[K\r", bytes, totalBytes)
		}
		fmt.Println()
	})

	return nil
}

func unzip(zipfile, destination string) error {
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
		target := path.Join(destination, f.Name)
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
