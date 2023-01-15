package util

import (
	"encoding/json"
	"fmt"
	"main/requests"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

/* VARS */

var (
	GIT_HOST       string
	JSON_DATA_FILE string
)

/* VARS */

/* Read Data from Json */
type Languages struct {
	Lang []string `json:"names,omitempty"`
	Main struct {
		Path string `json:"path,omitempty"`
		Url  string `json:"url,omitempty"`
	} `json:"mainfile,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
	Makefile   string `json:"makefile,omitempty"`
	Nixsh      string `json:"nix-shell,omitempty"`

	Commons []struct {
		Path string `json:"path,omitempty"`
		Url  string `json:"url,omitempty"`
	} `json:"commons,omitempty"`
	Utils []struct {
		Path string `json:"path,omitempty"`
		Url  string `json:"url,omitempty"`
	} `json:",omitempty"`
}
type JsonDataFile struct {
	GitHost   string      `json:"githost,omitempty"`
	Languages []Languages `json:"langs,omitempty"`
}

func GetJsonData(filename string) (JsonDataFile, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return JsonDataFile{}, err
	}

	var result JsonDataFile
	err = json.Unmarshal(content, &result)
	if err != nil {
		return JsonDataFile{}, err
	}

	return result, nil
}

/* Read Data from Json */

/*  */

func GetFileFromGit(url string) ([]byte, error) {
	data, err := requests.Get(url)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetLogo() string {
	/*return " [ Project Generator ]" */
	return `
 ___                                _     ___                    _              
(  _ \                _            ( )_  (  _ \                 ( )_            
| |_) )_ __   _      (_)  __    ___|  _) | ( (_)_ __   __    _ _|  _)  _   _ __ 
|  __/(  __)/ _ \    | |/ __ \/ ___) |   | |  _(  __)/ __ \/ _  ) |  / _ \(  __)
| |   | |  ( (_) )   | |  ___/ (___| |_  | (_( ) |  (  ___/ (_| | |_( (_) ) |   
(_)   (_)   \___/ _  | |\____)\____)\__) (____/(_)   \____)\__ _)\__)\___/(_)   
                 ( )_| |                                                        
                  \___/                                                         
`
}

/*  */

/* Create Project */
func NewCreateProject(name, language string) error {
	data, err := GetJsonData(JSON_DATA_FILE)
	if err != nil {
		return err
	}
	return CreateProject(data.Languages, name, language)
}

func CreateProject(language_data []Languages, name, language string) error {
	found := false
	for _, lang := range language_data {
		if slices.Contains(lang.Lang, strings.ToLower(language)) {
			os.Mkdir(name, 0755)

			/* SHELL.NIX */
			if lang.Nixsh != "" {
				content, err := GetFileFromGit(strings.Replace(lang.Nixsh, "<git>", GIT_HOST, -1))
				if err != nil {
					return err
				}

				os.WriteFile(name+"/shell.nix", content, 0755)
			}
			/* SHELL.NIX */

			/* MAKEFILE */
			if lang.Makefile != "" {
				content, err := GetFileFromGit(strings.Replace(lang.Makefile, "<git>", GIT_HOST, -1))
				if err != nil {
					return err
				}

				os.WriteFile(name+"/Makefile", content, 0755)
			}
			/* MAKEFILE */

			/* DOCKERFILE */
			if lang.Dockerfile != "" {
				content, err := GetFileFromGit(strings.Replace(lang.Dockerfile, "<git>", GIT_HOST, -1))
				if err != nil {
					return err
				}

				os.WriteFile(name+"/Dockerfile", content, 0755)
			}
			/* DOCKERFILE */

			/* Mainfile */
			if lang.Main.Url != "" {
				content, err := GetFileFromGit(strings.Replace(lang.Main.Url, "<git>", GIT_HOST, -1))
				if err != nil {
					return err
				}

				mainpath := strings.Split(lang.Main.Path, "/")
				path := strings.Join(mainpath[:len(mainpath)-1], "/")
				os.Mkdir(name+"/"+path, 0755)
				os.WriteFile(name+"/"+lang.Main.Path, content, 0755)
			}
			/* Mainfile */

			/* Commons */
			if lang.Commons != nil {
				for _, common := range lang.Commons {
					content, err := GetFileFromGit(strings.Replace(common.Url, "<git>", GIT_HOST, -1))
					if err != nil {
						return err
					}

					commonpath := strings.Split(common.Path, "/")
					path := strings.Join(commonpath[:len(commonpath)-1], "/")
					os.Mkdir(name+"/"+path, 0755)
					os.WriteFile(name+"/"+common.Path, content, 0755)
				}
			}
			/* Commons */

			/* Util */
			if lang.Utils != nil {
				for _, util := range lang.Utils {
					content, err := GetFileFromGit(strings.Replace(util.Url, "<git>", GIT_HOST, -1))
					if err != nil {
						return err
					}

					os.WriteFile(name+"/"+util.Path, content, 0755)
				}
			}
			/* Util */
			found = true
		}
	}

	if !found {
		return fmt.Errorf("unknown language")
	}

	return nil
}

/* Create Project */
