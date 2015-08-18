package sellsword

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path"
)

var Logger *log.Logger

func GetTermPrinter(colorName color.Attribute) func(...interface{}) string {
	newColor := color.New(colorName)
	newColor.EnableColor()
	return newColor.SprintFunc()
}

func GetTermPrinterF(colorName color.Attribute) func(string, ...interface{}) string {
	newColor := color.New(colorName)
	newColor.EnableColor()
	return newColor.SprintfFunc()
}

type AppSet struct {
	Apps []*App
	Home string
}

func (as *AppSet) FindApps(appNames ...string) {
	if _, err := os.Stat(as.Home); os.IsNotExist(err) {
		red := GetTermPrinterF(color.FgRed)
		Logger.Errorln(red("The Home directory that you have specified, %s, does not exist.", as.Home))
	} else {
		if appNames[0] == "all" {
			di, _ := ioutil.ReadDir(as.Home)
			for i := range di {
				if di[i].Name() != "config" {
					a := new(App)
					a.Path = path.Join(as.Home, di[i].Name())
					a.Parse()
					as.Apps = append(as.Apps, a)
				}
			}
		} else {
			for i := range appNames {
				a := new(App)
				a.Path = path.Join(as.Home, appNames[i])
				a.Parse()
				as.Apps = append(as.Apps, a)
			}
		}
	}
}

func (as *AppSet) ListApps(appNames []string) {
	if len(appNames) == 0 {
		as.FindApps("all")
	} else {
		as.FindApps(appNames...)
	}
	for i := range as.Apps {
		cyan := GetTermPrinter(color.FgCyan)
		red := GetTermPrinter(color.FgRed)
		green := GetTermPrinter(color.FgGreen)
		fmt.Printf("%s:\n", cyan(as.Apps[i].Name))
		current, err := as.Apps[i].Current()
		if err != nil {
			fmt.Printf("%s\n", red("No environment currently in use"))
		} else {
			fmt.Printf("\t%s\t%s\n", green(current.Name), green("CURRENT"))
		}
		envs := as.Apps[i].ListEnvs()
		for i := range envs {
			if envs[i].Name != current.Name {
				fmt.Printf("\t%s\n", envs[i].Name)
			}
		}
	}
}
