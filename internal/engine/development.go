package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Source struct {
	Name string
	Arch string
	OS   string
	Url  string
}

type Sources []Source

func (s Sources) resolve(name string) (string, error) {
	architectures := []string{runtime.GOARCH, ""}
	operatingSystems := []string{runtime.GOOS, ""}

	for _, a := range architectures {
		for _, o := range operatingSystems {
			for _, e := range s {
				if strings.EqualFold(e.Name, name) && strings.EqualFold(e.Arch, a) && strings.EqualFold(e.OS, o) {
					return e.Url, nil
				}
			}
		}
	}
	return "", fmt.Errorf("%v is not valid for %v - %v", name, runtime.GOOS, runtime.GOARCH)
}

var sources = Sources{
	Source{Name: "sources", Arch: "", OS: "", Url: "https://github.com/taylorza/zxenv/blob/main/sources.json?raw=true"},
	Source{Name: "cspect", Arch: "", OS: "", Url: "https://www.dropbox.com/s/6hcl37zyqqars4q/CSpect2_19_4_3.zip?dl=1"},
	Source{Name: "zesarux", Arch: "amd64", OS: "windows", Url: "https://github.com/chernandezba/zesarux/releases/download/ZEsarUX-10.3/ZEsarUX_windows-10.3.zip"},
	Source{Name: "core3", Arch: "", OS: "", Url: "https://github.com/taylorza/zxenv/raw/main/images/tbblue_core_3_02_00_os_2_08.zip"},
	Source{Name: "2gb", Arch: "", OS: "", Url: "http://www.zxspectrumnext.online/cspect/cspect-next-2gb.zip"},
	Source{Name: "4gb", Arch: "", OS: "", Url: "http://www.zxspectrumnext.online/cspect/cspect-next-4gb.zip"},
	Source{Name: "8gb", Arch: "", OS: "", Url: "http://www.zxspectrumnext.online/cspect/cspect-next-8gb.zip"},
	Source{Name: "16gb", Arch: "", OS: "", Url: "http://www.zxspectrumnext.online/cspect/cspect-next-16gb.zip"},
	Source{Name: "asm", Arch: "", OS: "", Url: "https://github.com/z00m128/sjasmplus/releases/download/v1.20.1/sjasmplus-1.20.1.win.zip"},
	Source{Name: "asm", Arch: "amd64", OS: "linux", Url: "https://github.com/taylorza/zxenv/blob/main/images/sjasmplus-1.20.1.linux.zip?raw=true"},
	Source{Name: "asm", Arch: "", OS: "darwin", Url: "https://github.com/taylorza/zxenv/blob/main/images/sjasmplus-1.20.1.mac.zip?raw=true"},
	Source{Name: "hdfmonkey", Arch: "", OS: "", Url: "http://uto.speccy.org/downloads/hdfmonkey_windows.zip"},
}

func SetupDevelopment(env *Environment) error {
	fmt.Println("Setting up development environment")
	fmt.Printf("Emulator: %v\n", env.Emulator)
	fmt.Printf("SD Card Size: %v\n", env.SDSize)
	fmt.Println()

	err := downloadSources(env)
	if err != nil {
		fmt.Println(err)
	}

	err = makeDirectories(env)
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

func downloadSources(env *Environment) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to resolve executable path (%w)", err)
	}

	dir := path.Dir(exePath)
	sourcesFile := filepath.Join(dir, "sources.json")

	if _, err := os.Stat(sourcesFile); os.IsNotExist(err) {
		sourcesFile = filepath.Join(env.BasePath, "sources.json")
	}

	_, err = os.Stat(sourcesFile)
	if os.IsNotExist(err) {
		fmt.Println("Downloading sources")
		err = download("sources", sourcesFile)
	}

	if err != nil {
		return fmt.Errorf("failed to acquire latest sources.json (%w)", err)
	}

	sourcesBytes, err := os.ReadFile("sources.json")
	if err != nil {
		fmt.Printf("failed to read sources.json (%s)\n", err)
	}

	err = json.Unmarshal(sourcesBytes, &sources)
	if err != nil {
		return fmt.Errorf("failed to parse sources.json (%s)\n", err)
	}
	return nil
}

func makeDirectories(env *Environment) error {
	err := os.MkdirAll(env.TempPath(), 0777)
	if err != nil {
		return fmt.Errorf("failed for create temporary directory (%w)", err)
	}

	err = os.MkdirAll(env.EmulatorPath(), 0777)
	if err != nil {
		return fmt.Errorf("failed for create emulator directory (%w)", err)
	}

	err = os.MkdirAll(env.SDPath(), 0777)
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
		if !d.IsDir() && d.Name() == "sjasmplus" || d.Name() == "sjasmplus.exe" {
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
	if err != nil {
		return 0, err
	}

	err = os.Chmod(dst, sourceFileStat.Mode())
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

	url, err := sources.resolve(urlkey)
	if err != nil {
		return err
	}
	downloadFile(url, filename, reportCopyProgress)

	return nil
}
