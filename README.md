# Project Creator
A Simple Program to Create projects written in GoLang using [bubbletea](https://github.com/charmbracelet/bubbletea/) for its Terminal UI


https://user-images.githubusercontent.com/85793326/212573374-4d925130-30b0-49a3-b629-8f58398cd859.mp4


# Using this
```
go build .
./main [config file]
```


# Configuration JSON File
The Configuration is done using a simple json file, which should looke like this: ([configuration file with example values](example.json))
```json
{
    "githost": "ip or hostname for git host",
    "langs": [
        {
            "names": ["language names", "and aliases"],
            "mainfile": {
                "path": "path to main file (including filename.filetype)",
                "url": "url to file on git repo, use <git> for the program to replace it with githost automatically"
            },
            "makefile": "url to makefile on git repo, rest is same as above",
            "nix-shell": "url to shell.nix file on git repo, rest is the same as above",
            
            "commons": [
                {
                    "path": "same as mainfile",
                    "url": "same as mainfile you can use an array tho"
                }
            ]
        }
    ]
}
```



