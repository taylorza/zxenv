package engine

import (
	"embed"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
)

//go:embed templates/*
var tmplfs embed.FS

type project struct {
	DevPath  string
	Emulator string
	Name     string
}

func CreateProject(env *Environment, name string) error {
	devpath, err := filepath.Abs(env.BasePath)
	if err != nil {
		return err
	}

	devpath = filepath.ToSlash(devpath)
	p := project{devpath, env.Emulator, name}

	os.MkdirAll(filepath.Join(devpath, p.Name, "inc"), 0777)
	os.MkdirAll(filepath.Join(devpath, p.Name, "src"), 0777)

	tmplFuncs := template.FuncMap{
		"isWindows": func() bool {
			return runtime.GOOS == "Windows"
		},
	}

	t, err := template.New("templates").Funcs(tmplFuncs).ParseFS(tmplfs, "templates/project/*.tmpl")
	if err != nil {
		return err
	}

	renderFile(p, t, ".vscode", "extensions.json")
	if err != nil {
		return err
	}

	renderFile(p, t, ".vscode", "launch.json")
	if err != nil {
		return err
	}

	renderFile(p, t, ".vscode", "tasks.json")
	if err != nil {
		return err
	}

	renderFile(p, t, "src", "main.asm")
	if err != nil {
		return err
	}

	return nil
}

func renderFile(p project, t *template.Template, subdir, filename string) error {
	dir := path.Join(p.DevPath, p.Name, subdir)
	fullfilename := path.Join(dir, filename)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(fullfilename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, filename+".tmpl", p)
	if err != nil {
		return err
	}
	return nil
}
