package engine

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

var urls = map[string]string{
	"cspect":      "http://www.javalemmings.com/public/zxnext/CSpect2_16_5.zip",
	"zesarux":     "https://github.com/chernandezba/zesarux/releases/download/ZEsarUX-10.1/ZEsarUX_windows-10.1.zip",
	"core3-128mb": "https://github.com/taylorza/zxenv/blob/main/images/tbblue_core_3_01_10_os_2_07g.zip?raw=true",
	"2gb":         "http://www.zxspectrumnext.online/cspect/cspect-next-2gb.zip",
	"4gb":         "http://www.zxspectrumnext.online/cspect/cspect-next-4gb.zip",
	"8gb":         "http://www.zxspectrumnext.online/cspect/cspect-next-8gb.zip",
	"16gb":        "http://www.zxspectrumnext.online/cspect/cspect-next-16gb.zip",
	"asm":         "https://github.com/z00m128/sjasmplus/releases/download/v1.20.0/sjasmplus-1.20.0.win.zip",
	"hdfmonkey":   "http://uto.speccy.org/downloads/hdfmonkey_windows.zip",
	"dezog":       "https://github.com/maziac/DeZogPlugin/releases/download/v2.1.0/DeZogPlugin.dll",
	"dezog-conf":  "https://raw.githubusercontent.com/maziac/DeZogPlugin/main/DeZogPlugin.dll.config",
}

func SetupDevelopment(env *Environment) error {
	fmt.Println("Setting up development environment")
	fmt.Printf("Emulator: %v\n", env.Emulator)
	fmt.Printf("SD Card Size: %v\n", env.SDSize)
	fmt.Println()

	err := makeDirectories(env)
	if err != nil {
		return err
	}

	err = installEmulator(env)
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

	err = env.Save()
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
	zippath := path.Join(env.TempPath(), "emulator.zip")
	fmt.Println("Downloading emulator")
	err := download(env.Emulator, zippath)
	if err != nil {
		return fmt.Errorf("failed to install emulator (%w)", err)
	}

	fmt.Println("Unzipping emulator")
	if env.Emulator == "zesarux" {
		err = unzipAndStrip(zippath, env.EmulatorPath(), true)
	} else {
		err = unzip(zippath, env.EmulatorPath())
	}
	if err != nil {
		return fmt.Errorf("failed to unzip the emulator (%w)", err)
	}
	return nil
}

func installSDTool(env *Environment) error {
	downloadDst := path.Join(env.TempPath(), "hdfmonkey.zip")

	fmt.Println("Downloading hdfmonkey")
	err := download("hdfmonkey", downloadDst)
	if err != nil {
		return fmt.Errorf("failed to download hdfmonkey (%w)", err)
	}

	fmt.Println("Unzipping hdfmonkey")
	err = unzip(downloadDst, env.EmulatorPath())
	if err != nil {
		return fmt.Errorf("failed to unzip hdfmonkey (%w)", err)
	}

	return nil
}

func setupSdCard(env *Environment) error {
	downloadDst := path.Join(env.TempPath(), env.SDCardName())

	fmt.Println("Downloading SD Card")
	err := download(env.SDSize, downloadDst)
	if err != nil {
		return fmt.Errorf("failed to download sd card image (%w)", err)
	}

	fmt.Println("Unzipping SD Card")
	err = unzip(downloadDst, env.TempPath())
	if err != nil {
		return fmt.Errorf("failed to unzip the sd card image (%w)", err)
	}

	files, err := os.ReadDir(env.TempPath())
	if err != nil {
		return err
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".img" || path.Ext(file.Name()) == ".mmc" {
			fmt.Printf("Copying %v\n", file.Name())
			_, err = copyFile(path.Join(env.TempPath(), file.Name()), path.Join(env.SDPath(), "tbblue-dev.sd"))
			if err != nil {
				return err
			}
		}

		if path.Ext(file.Name()) == ".rom" {
			fmt.Printf("Copying %v\n", file.Name())
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

	fmt.Println("Downloading sjasmplus")
	err := download("asm", zippath)
	if err != nil {
		return fmt.Errorf("failed to download sd card image (%w)", err)
	}

	fmt.Println("Unzipping sjasmplus")
	err = unzip(zippath, unzippath)
	if err != nil {
		return fmt.Errorf("failed to unzip the sd card image (%w)", err)
	}

	err = filepath.WalkDir(unzippath, func(fullpath string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == "sjasmplus.exe" {
			fmt.Printf("Copying %v\n", d.Name())
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

func reportCopyProgress(progress <-chan int) {
	for bytes := range progress {
		fmt.Printf("[%v]\x1b[K\r", bytes)
	}
	fmt.Println()
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

	downloadFile(urls[urlkey], filename, reportCopyProgress)

	return nil
}
