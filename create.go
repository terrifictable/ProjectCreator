package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"projs/common"
	"strings"
)

func CreateProject(name string, data ProjectTemplate) error {
	err := os.Mkdir(name, 0755)
	if err != nil {
		return err
	}

	if data.Git {
		cmd := exec.Command("git", "init")
		cmd.Dir = "./" + name + "/"
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		common.CMD("%s", out)
	}

	err = createFiles(name, data.Files)
	if err != nil {
		return err
	}

	for _, command := range data.Commands {
		_cmd := strings.Split(command, " ")
		cmd := exec.Command(_cmd[0], _cmd[1:]...)
		cmd.Dir = "./" + name + "/"
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		common.CMD("$ %s \n\t-> %s", command, out)
	}

	return nil
}

func createFiles(name string, files map[string]any) error {
	for file, v := range files {
		if url, ok := v.(string); ok {
			f, err := os.Create(name + "/" + file)
			if err != nil {
				return err
			}

			res, err := http.Get(url)
			if err != nil {
				return err
			}

			data, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			res.Body.Close()

			f.WriteString(string(data))
			f.Close()
		} else if folder, ok := v.(map[string]any); ok {
			err := os.Mkdir(name+"/"+file, 0755)
			if err != nil {
				return err
			}
			createFiles(name+"/"+file, folder)
		}
	}

	return nil
}
