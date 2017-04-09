package main // import "loe.yt/sct"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"text/template"
)

var (
	homeDir    = os.Getenv("HOME")
	configFile = path.Join(homeDir, ".ssh", "config")
	configDir  = path.Join(homeDir, ".sct")
)

const marker = "# THIS FILE WAS GENERATED AUTOMATICALLY USING loe.yt/sct\n\n"

func main() {
	log.SetFlags(0)
	dir, err := os.Open(configDir)
	if err != nil {
		log.Fatalf("can't open %s: %v\n", configDir, err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatalf("unable to read directory names: %v\n", err)
	}
	sort.Strings(names)
	b := bytes.NewBuffer(nil)
	b.WriteString(marker)
	for _, name := range names {
		hosts, err := loadHosts(name)
		if err != nil {
			log.Fatalf("error loading hosts for %s: %v\n", name, err)
		}
		err = execute(name, hosts, b)
		if err != nil {
			log.Fatalf("error executing template for %s: %v\n", name, err)
		}
	}
	err = ioutil.WriteFile(configFile, b.Bytes(), 0600)
	if err != nil {
		log.Fatalf("error writing to %s: %v\n", configFile, err)
	}
}

type host struct {
	Name  string
	Value interface{}
}

func loadHosts(name string) ([]host, error) {
	file, err := os.Open(path.Join(configDir, name, "hosts.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := json.NewDecoder(file)
	var v map[string]interface{}
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	l := make([]host, 0)
	for name, value := range v {
		l = append(l, host{Name: name, Value: value})
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].Name < l[j].Name
	})
	return l, nil
}

const mainTemplate = `{{range .}}{{template "template" .}}{{end}}`

var mainTpl = template.Must(template.New("main").Parse(mainTemplate))

func execute(name string, hosts []host, w io.Writer) error {
	tpl, err := mainTpl.Clone()
	if err != nil {
		return err
	}
	tpl, err = tpl.ParseFiles(path.Join(configDir, name, "template"))
	if err != nil {
		return err
	}
	err = tpl.ExecuteTemplate(w, "main", hosts)
	if err != nil {
		return err
	}
	return nil
}
