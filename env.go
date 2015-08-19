package sellsword

import (
	"fmt"
	"github.com/fatih/color"
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
	e.Name = path.Base(e.Path)
	e.App = a
	return nil
}

func (e *Env) Save() error {
	if e.EnvType != "environment" {
		Logger.Warnf("Environment type %s does not currently support the save operation", e.EnvType)
		return nil
	}
	if d, err := yaml.Marshal(&e.Variables); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(e.Path, d, 0775); err != nil {
			return err
		} else {
			green := GetTermPrinterF(color.FgGreen)
			fmt.Print(green("New environment created at %s\n", e.Path))
			return nil
		}
	}
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
