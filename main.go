package main

import (
	"fmt"
	"main/ui"
	"main/util"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		execpath := strings.Split(os.Args[0], "/")
		thisFile := execpath[len(execpath)-1]

		fmt.Printf("Usage: %s [json_language_file]\n", thisFile)

		os.Exit(1)
	}

	util.JSON_DATA_FILE = os.Args[1]
	data, err := util.GetJsonData(util.JSON_DATA_FILE)
	if err != nil {
		fmt.Printf("Error opening json file: %s\n", err)
		os.Exit(1)
	}
	util.GIT_HOST = data.GitHost

	if err := ui.InitUI(); err != nil {
		panic(err)
	}
}
