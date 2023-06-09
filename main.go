package main

import (
	"encoding/json"
	"fmt"
	"os"
	"projs/common"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	common.DBG("Hi")

	if len(os.Args) < 2 {
		tmp_exe, _ := os.Executable()
		exe := strings.Split(strings.ReplaceAll(tmp_exe, "\\", "/"), "/")
		common.ERR("Usage: %s <config file>", exe[len(exe)-1])
		return
	}

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var out map[string]interface{}
	err = yaml.Unmarshal(file, &out)
	if err != nil {
		panic(err)
	}

	out_pretty, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out_pretty)

}
