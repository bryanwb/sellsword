package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path"
)

type AppSet struct {
	Apps []*App
	Home string
}

func (as *AppSet) findApps(appNames ...string) {
	if _, err := os.Stat(as.Home); os.IsNotExist(err) {
		red := getTermPrinterF(color.FgRed)
		log.Errorln(red("The Home directory that you have specified, %s, does not exist.", as.Home))
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

func (as *AppSet) listApps(appNames []string) {
	if len(appNames) == 0 {
		as.findApps("all")
	} else {
		as.findApps(appNames...)
	}
	for i := range as.Apps {
		cyan := getTermPrinter(color.FgCyan)
		red := getTermPrinter(color.FgRed)
		green := getTermPrinter(color.FgGreen)
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
