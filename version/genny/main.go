package main

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

type VersionGennyStruct struct {
	BinName          string
	Version          string
	CompileTimestamp string
}

var versionFileLoc = "./version/version.go" // relative to genny.go
var versionTemplate = `// Generated Code. DO NOT EDIT.
package version

const (
	BIN = "{{ $.BinName }}"
	VERSION = "{{ $.Version }}"
	COMPILED_AT = "{{ $.CompileTimestamp }}"
)
`

func main() {
	templater, err := template.New("vers").Parse(versionTemplate)
	if err != nil {
		fmt.Println("version-genny: %s", err)
		os.Exit(1)
	}

	fout, err := os.Create(versionFileLoc)
	if err != nil {
		fmt.Println("version-genny: %s", err)
		os.Exit(1)
	}
	defer fout.Close()

	err = templater.Execute(fout, &VersionGennyStruct{
		BinName: "noname001",
		Version: "0.0.1",
		CompileTimestamp: time.Now().UTC().Format("20060102T150405Z"),
	})
	if err != nil {
		fmt.Println("version-genny: %s", err)
		os.Exit(1)
	}
}
