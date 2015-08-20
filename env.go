package sellsword

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Env struct {
	Name            string
	Path            string
	Current         bool
	EnvType         string
	ExportVariables map[string]string
	Variables       map[string]string
}

func NewEnv(name string, basePath string, exportVars map[string]string, vars []string,
	envType string) (*Env, error) {
	env := new(Env)
	env.Name = name
	env.EnvType = envType
	env.Path = path.Join(basePath, name)
	if envType == "environment" {
		env.ExportVariables = exportVars
		env.Variables = arrayToEmptyMap(vars)
	}
	return env, nil
}

// NewEnvironmentEnv is a factory method that properly initializes the Env struct for the env type of Environment
func NewEnvironmentEnv(name string, basePath string, exportVars map[string]string, vars []string) (*Env, error) {
	return NewEnv(name, basePath, exportVars, vars, "environment")
}

// NewDirectoryEnv is a factory method that properly initializes the Env struct for the env type of Directory
func NewDirectoryEnv(name string, basePath string) (*Env, error) {
	return NewEnv(name, basePath, map[string]string{}, []string{}, "directory")
}

// func (e *Env) Parse(a *App) error {
// 	e.Name = path.Base(e.Path)
// 	return nil
// }

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
	yamlVars := make(map[string]string)
	Logger.Debugf("e.Path is %s", e.Path)
	d, err := ioutil.ReadFile(e.Path)
	if err == nil {
		yaml.Unmarshal(d, yamlVars)
		for key, value := range e.ExportVariables {
			if v, ok := yamlVars[value]; ok {
				e.ExportVariables[key] = v
			} else {
				delete(e.ExportVariables, key)
			}
		}
	}
	Logger.Debugf("env export vars are %v", e.ExportVariables)
	return err
}

func (e *Env) PrintExports() {
	exports := make([]string, len(e.ExportVariables))
	for key, value := range e.ExportVariables {
		exports = append(exports, "export "+key+"="+value)
	}
	fmt.Println(strings.Join(exports, "\n"))
}

// *Constructs* a new environment, not to be confused w/ the Constructor NewEnv
// In case of environment type, queries user for values
// Not implemented for other types yet
func (e *Env) Construct() error {
	if e.EnvType == "environment" {
		for k, _ := range e.Variables {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("%s: ", k)
			if text, err := reader.ReadString('\n'); err != nil {
				return err
			} else {
				e.Variables[k] = strings.TrimSpace(text)
			}
		}
		if err := e.Save(); err != nil {
			Logger.Errorf("error: %v", err)
			return err
		}
	} else {
		red := GetTermPrinterF(color.FgRed)
		fmt.Fprint(os.Stderr, red("new command not implemented for environment type %s", e.EnvType))
	}
	return nil
}
