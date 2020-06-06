package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

type (
	Patch struct {
		File    string `yaml:"file"`
		Regex   string `yaml:"regex"`
		Replace string `yaml:"replace"`
	}
)

func main() {
	patches := load()
	for i, p := range patches {
		log.Println(i, ":", p.File)
		replace(p)
	}
}

func load() []Patch {
	cd := path.Dir(os.Args[0])
	b, err := ioutil.ReadFile(path.Join(cd, "config.yaml"))
	if err != nil {
		log.Fatalln(err)
	}
	patches := make([]Patch, 0)
	if err := yaml.Unmarshal(b, patches); err != nil {
		log.Fatalln(err)
	}
	return patches
}

func replace(patch Patch) {
	b, err := ioutil.ReadFile(patch.File)
	if err != nil {
		log.Println("file read error:", patch.File)
		return
	}
	file := string(b)

	switch {
	case patch.Regex != "":
		regex := regexp.MustCompile(patch.Regex)
		file = regex.ReplaceAllString(file, patch.Replace)
	}

	if err := ioutil.WriteFile(patch.File, []byte(file), 644); err != nil {
		log.Println("file write error:", patch.File)
	}
}
