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
	b.WriteString("# THIS FILE WAS GENERATED AUTOMATICALLY USING loe.yt/sct\n\n")
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

type Host struct {
	Name  string
	Value interface{}
}

type Hosts []Host

func (hs *Hosts) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	for name, value := range v {
		*hs = append(*hs, Host{Name: name, Value: value})
	}
	sort.Sort(*hs)
	return nil
}

func (hs Hosts) Len() int {
	return len(hs)
}

func (hs Hosts) Less(i, j int) bool {
	return hs[i].Name < hs[j].Name
}

func (hs Hosts) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

func loadHosts(name string) (Hosts, error) {
	var v Hosts
	file, err := os.Open(path.Join(configDir, name, "hosts.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := json.NewDecoder(file)
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

const mainTemplate = `{{range .}}{{template "template" .}}{{end}}`

var mainTpl = template.Must(template.New("main").Parse(mainTemplate))

func execute(name string, hosts Hosts, w io.Writer) error {
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
