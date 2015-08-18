package sellsword

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

type Env struct {
	Name      string
	App       *App
	Path      string
	Current   bool
	EnvType   string
	Variables map[string]string
}

func (e *Env) Parse(a *App) error {
	e.Name = strings.Split(path.Base(e.Path), "-env.ssw")[0]
	e.App = a
	return nil
}

func (e *Env) Generate() error {
	return nil
}

func (e *Env) Remove() error {
	return nil
}

func (e *Env) PopulateExportVars() error {
	e.Variables = make(map[string]string)
	envVars := make(map[string]string)
	envData, err := ioutil.ReadFile(e.Path)
	if err == nil {
		yaml.Unmarshal(envData, envVars)
		for key, value := range e.App.ExportVariables {
			if value, ok := envVars[value]; ok {
				e.Variables[key] = value
			}
		}
	}
	return err
}

func (e *Env) PrintExports() {
	exports := make([]string, len(e.Variables))
	for key, value := range e.Variables {
		exports = append(exports, "export "+key+"="+value)
	}
	fmt.Println(strings.Join(exports, "\n"))
}
