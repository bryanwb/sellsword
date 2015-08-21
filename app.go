package sellsword

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

type App struct {
	Name            string
	EnvType         string `yaml:"type"`
	Path            string
	Root            string
	Target          string
	Definition      string
	Variables       []string
	VariableNames   []string
	ExportVariables map[string]string
}

// NewApp is the constructor for New Apps
func NewApp(name string, sswHome string) (*App, error) {
	a := new(App)
	a.Name = name
	a.Definition = path.Join(sswHome, "config", name+".ssw")
	a.Path = path.Join(sswHome, name)
	Logger.Debugf("Parsing application found at %s", a.Path)
	if data, err := ioutil.ReadFile(a.Definition); err != nil {
		Logger.Errorln(err.Error())
		return a, err
	} else {
		if err := yaml.Unmarshal(data, a); err != nil {
			return a, err
		}
		if a.EnvType == "directory" {
			Logger.Debugf("Target for %s is currently %s", a.Name, a.Target)
			if newTarget, err := expandPath(a.Target); err != nil {
				Logger.Debugf(err.Error())
				return a, err
			} else {
				Logger.Debugf("New target for %s is %s", a.Name, newTarget)
				a.Target = newTarget
				return a, err
			}
		} else {
			if err := a.ParseExportVars(); err != nil {
				return a, err
			}
		}
		return a, nil
	}
	return a, nil
}

func (a *App) ParseExportVars() error {
	a.VariableNames = make([]string, 0)
	a.ExportVariables = make(map[string]string, len(a.Variables))
	for i := range a.Variables {
		keyValue := strings.Split(a.Variables[i], "=")
		a.VariableNames = appendIfMissing(a.VariableNames, keyValue[0])
		a.ExportVariables[keyValue[1]] = keyValue[0]
	}
	return nil
}

func (a *App) Current() (*Env, error) {
	var e *Env
	if realPath, err := resolveSymlink(path.Join(a.Path, "current")); err != nil {
		return e, err
	} else {
		envName := path.Base(realPath)
		if a.EnvType == "environment" {
			a.ParseExportVars()
			return NewEnvironmentEnv(envName, a.Path, a.ExportVariables, a.VariableNames)
		} else {
			return NewDirectoryEnv(envName, a.Path)
		}
	}

}

func (a *App) ListEnvs() []*Env {
	a.ParseExportVars()
	envs := make([]*Env, 0)
	di, _ := ioutil.ReadDir(a.Path)
	for i := range di {
		name := di[i].Name()
		if name != "current" {
			var e *Env
			if a.EnvType == "environment" {
				e, _ = NewEnvironmentEnv(name, a.Path, a.ExportVariables, a.VariableNames)
			} else {
				e, _ = NewDirectoryEnv(name, a.Path)
			}
			envs = append(envs, e)
		}
	}
	return envs
}

func (a *App) Load() {
	if a.EnvType == "environment" {
		Logger.Debugf("Exporting environment variables for application %s\n", a.Name)
		a.ParseExportVars()
		if env, err := a.Current(); err == nil {
			env.PopulateExportVars()
			env.PrintExports()
		} else {
			Logger.Error(err.Error())
		}
	} else {
		Logger.Debugf("Application %s has no environment variables to export, nothing to do\n", a.Name)
	}
}

func (a *App) MakeCurrent(envName string) error {
	red := GetTermPrinterF(color.FgRed)
	envPath := path.Join(a.Path, envName)
	currentEnv, currentErr := a.Current()
	Logger.Debugf("Current env is %s", currentEnv.Name)
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		Logger.Error(err.Error())
		return err
	} else if currentErr == nil && envName == currentEnv.Name {
		Logger.Warn(red("%s is already set as the default environment for application %s. Nothing to do.",
			envName, a.Name))
		return nil
	} else {
		var newEnv *Env
		if a.EnvType == "environment" {
			a.ParseExportVars()
			newEnv, err = NewEnvironmentEnv(envName, a.Path, a.ExportVariables, a.VariableNames)
			newEnv.PopulateExportVars()
			a.UnsetExportVars()
			newEnv.PrintExports()
		} else {
			newEnv, err = NewDirectoryEnv("current", a.Path)
		}
		if err := a.Unlink(); err != nil {
			Logger.Debugf("Encountered error when unlinking current for %s", a.Name)
			return err
		}
		return a.Link(newEnv.Name)
	}
}

func (a *App) EnumerateExportVars() []string {
	vars := make([]string, len(a.ExportVariables))
	i := 0
	for k, _ := range a.ExportVariables {
		vars[i] = k
		i++
	}
	return vars
}

func (a *App) MakeUnsetExportVars() string {
	vars := a.EnumerateExportVars()
	statements := make([]string, 0)
	for i := range vars {
		statements = append(statements, fmt.Sprintf("unset %s", vars[i]))
	}
	sort.Strings(statements)
	return strings.Join(statements, "\n")
}

func (a *App) UnsetExportVars() {
	fmt.Printf(a.MakeUnsetExportVars())
}

func (a *App) Unlink() error {
	current := path.Join(a.Path, "current")
	if _, err := os.Lstat(current); os.IsNotExist(err) {
		Logger.Debugf("Current symlink %s does not exist, nothing to do", current)
		return nil
	} else {
		if a.EnvType == "directory" {
			Logger.Debugf("Removing Target symlink for %s at %s", a.Name, a.Target)
			if err := os.Remove(a.Target); err != nil {
				return err
			}
		}
		return os.Remove(path.Join(a.Path, "current"))
	}
}

func (a *App) Link(envName string) error {
	source := path.Join(a.Path, envName)
	target := path.Join(a.Path, "current")
	if err := os.Symlink(source, target); err != nil {
		Logger.Debugf(err.Error())
		return err
	}
	if a.EnvType == "directory" {
		if err := os.Symlink(source, a.Target); err != nil {
			Logger.Debugf(err.Error())
			return err
		}
	}
	return nil
}

func (a *App) NewEnv(envName string) (*Env, error) {
	if a.EnvType == "environment" {
		return NewEnvironmentEnv(envName, a.Path, a.ExportVariables, a.VariableNames)
	} else {
		return NewDirectoryEnv(envName, a.Path)
	}

}
