package main

import (
	"io"
	"net/http"
	"os"
	"projs/common"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProjectTemplate struct {
	/* languages the project is written in (optional) */
	Languages []string `yaml:"languages"`

	/* whether to initialize a git repo or not */
	Git bool `yaml:"git"`

	/* files to download (map is: `foldername: files` `filename: url`) */
	Files map[string]any `yaml:"files"`

	/* commands to execute after creating downloading files and (maybe) initializing git repo */
	Commands []string `yaml:"cmds"`
}

func main() {
	if len(os.Args) < 3 {
		tmp_exe, _ := os.Executable()
		exe := strings.Split(strings.ReplaceAll(tmp_exe, "\\", "/"), "/")
		common.ERR("Usage: %s <project name> <config file>", exe[len(exe)-1])
		return
	}

	var yml []byte
	var err error

	if strings.HasPrefix(os.Args[2], "http://") || strings.HasPrefix(os.Args[1], "https://") {
		res, err := http.Get(os.Args[2])
		if err != nil {
			panic(err)
		}

		yml, err = io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		res.Body.Close()
	} else {
		yml, err = os.ReadFile(os.Args[2])
		if err != nil {
			panic(err)
		}
	}

	var out ProjectTemplate
	err = yaml.Unmarshal(yml, &out)
	if err != nil {
		panic(err)
	}

	err = CreateProject(os.Args[1], out)
	if err != nil {
		panic(err)
	}
}
