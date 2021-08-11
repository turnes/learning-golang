package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turnes/learning-golang/Gophercises/URLShortener/urlshort"
)

var URLPATH = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

const YAML = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

func main() {
	yamlFile := flag.String("yaml", "", "yaml file")
	flag.Parse()
	yaml, err := loadYaml(*yamlFile)
	if err != nil {
		panic(err)
	}
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(URLPATH, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe("0.0.0.0:8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func loadYaml(yamlFile string) ([]byte, error) {
	if isFlagPassed("yaml") {
		file, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			return nil, nil
		}
		return file, nil
	}
	return []byte(YAML), nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
