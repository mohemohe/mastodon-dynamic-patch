package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sync"

	"gopkg.in/yaml.v2"
)

type (
	Patch struct {
		File    string `yaml:"file"`
		Regex   string `yaml:"regex"`
		Replace string `yaml:"replace"`
	}
)

var patched bool

func main() {
	patched = false
	go listen()

	patches := load()
	for i, p := range patches {
		log.Println(i, ":", p.File)
		replace(p)
	}
	patched = true

	w := new(sync.WaitGroup)
	w.Add(1)
	w.Wait()
}

func listen() {
	addr := os.Getenv("PATCH_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if patched {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusTooEarly)
			_, _ = w.Write([]byte("NG"))
		}
	})
	log.Println("server listen on " + addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}

func load() []Patch {
	url := os.Getenv("PATCH_CONFIG_URL")
	var config []byte
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		config = b
	} else {
		cd := path.Dir(os.Args[0])
		b, err := ioutil.ReadFile(path.Join(cd, "config.yaml"))
		if err != nil {
			log.Fatalln(err)
		}
		config = b
	}
	patches := make([]Patch, 0)
	if err := yaml.Unmarshal(config, &patches); err != nil {
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