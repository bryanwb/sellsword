package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func resolveSymlink(symlink string) (string, error) {
	fi, err := os.Lstat(symlink)
	if err != nil {
		log.Debugf("Path %s does not exist\n", symlink)
		return "", err
	} else {
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			if realPath, err := os.Readlink(symlink); err == nil {
				return realPath, nil
			} else {
				return "", err
			}
		} else {
			return "", errors.New(fmt.Sprintf("Path %s exists but is not a symlink\n", symlink))
		}
	}
}

// Why the fuck isn't this in the golang stdlib?
func expandPath(pathName string) (string, error) {
	if string(pathName[0]) == "~" {
		relative := strings.Split(pathName, "~")[1]
		usr, _ := user.Current()
		return path.Join(usr.HomeDir, relative), nil
	} else {
		return filepath.Abs(pathName)
	}
}

type App struct {
	Name            string
	EnvType         string `yaml:"type"`
	Path            string
	Root            string
	Target          string
	Definition      string
	Variables       []string
	ExportVariables map[string]string
}

func (a *App) Parse() error {
	log.Debugf("Parsing application found at %s", a.Path)
	a.Name = path.Base(a.Path)
	a.Definition = path.Join(path.Dir(a.Path), "config", a.Name+".ssw")
	if data, err := ioutil.ReadFile(a.Definition); err != nil {
		log.Errorln(err.Error())
		return err
	} else {
		if err := yaml.Unmarshal(data, a); err != nil {
			return err
		}
		if a.EnvType == "directory" {
			log.Debugf("Target for %s is currently %s", a.Name, a.Target)
			if newTarget, err := expandPath(a.Target); err != nil {
				log.Debugf(err.Error())
				return err
			} else {
				log.Debugf("New target for %s is %s", a.Name, newTarget)
				a.Target = newTarget
				return nil
			}
		}
		return nil
	}
}

func (a *App) ParseExportVars() error {
	a.ExportVariables = make(map[string]string, len(a.Variables))
	for i := range a.Variables {
		keyValue := strings.Split(a.Variables[i], "=")
		a.ExportVariables[keyValue[1]] = keyValue[0]
	}
	return nil
}

func (a *App) Current() (*Env, error) {
	env := new(Env)
	env.EnvType = a.EnvType
	currentPath := path.Join(a.Path, "current")
	realPath, err := resolveSymlink(currentPath)
	if err == nil {
		env.Path = realPath
		err := env.Parse(a)
		return env, err
	} else {
		return env, err
	}
}

func (a *App) ListEnvs() []*Env {
	envs := make([]*Env, 0)
	di, _ := ioutil.ReadDir(a.Path)
	for i := range di {
		name := di[i].Name()
		if name != "current" {
			e := new(Env)
			e.Path = path.Join(a.Path, name)
			e.Parse(a)
			envs = append(envs, e)
		}
	}
	return envs
}

func (a *App) Load() {
	if a.EnvType == "environment" {
		log.Debugf("Exporting environment variables for application %s\n", a.Name)
		a.ParseExportVars()
		if env, err := a.Current(); err == nil {
			env.PopulateExportVars()
			env.PrintExports()
		} else {
			log.Error(err.Error())
		}
	} else {
		log.Debugf("Application %s has no environment variables to export, nothing to do\n", a.Name)
	}
}

func (a *App) DetermineEnvPath(envName string) string {
	if a.EnvType == "environment" {
		return path.Join(a.Path, envName+"-env.ssw")
	} else {
		return path.Join(a.Path, envName)
	}
}

func (a *App) MakeCurrent(envName string) error {
	red := getTermPrinterF(color.FgRed)
	envPath := a.DetermineEnvPath(envName)
	currentEnv, currentErr := a.Current()
	log.Debugf("Current env is %s", currentEnv.Name)
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Error(err.Error())
		return err
	} else if currentErr == nil && envName == currentEnv.Name {
		log.Warn(red("%s is already set as the default environment for application %s. Nothing to do.",
			envName, a.Name))
		return nil
	} else {
		newEnv := new(Env)
		newEnv.Path = envPath
		newEnv.EnvType = a.EnvType
		newEnv.Parse(a)
		if a.EnvType == "environment" {
			a.ParseExportVars()
			newEnv.PopulateExportVars()
			a.UnsetExportVars()
			newEnv.PrintExports()
		}
		if err := a.Unlink(); err != nil {
			log.Debugf("Encountered error when unlinking current for %s", a.Name)
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

func (a *App) UnsetExportVars() {
	vars := a.EnumerateExportVars()
	for i := range vars {
		fmt.Printf("unset %s\n", vars[i])
	}
}

func (a *App) Unlink() error {
	current := path.Join(a.Path, "current")
	if _, err := os.Lstat(current); os.IsNotExist(err) {
		log.Debugf("Current symlink %s does not exist, nothing to do", current)
		return nil
	} else {
		if a.EnvType == "directory" {
			log.Debugf("Removing Target symlink for %s at %s", a.Name, a.Target)
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
		log.Debugf(err.Error())
		return err
	}
	if a.EnvType == "directory" {
		if err := os.Symlink(source, a.Target); err != nil {
			log.Debugf(err.Error())
			return err
		}
	}
	return nil
}
