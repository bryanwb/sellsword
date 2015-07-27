package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

var (
	logger *log.Logger
)

type sswConfig struct {
	Name       string
	ConfigType string `yaml:"type"`
	Variables  []string
}

func (c *sswConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c *sswConfig) ParseExportVars() (map[string]string, error) {
	exportVars := make(map[string]string)
	for i := range c.Variables {
		keyValue := strings.Split(c.Variables[i], "=")
		exportVars[keyValue[1]] = keyValue[0]
	}
	return exportVars, nil
}

func mergeEnvMap(dest map[string]string, src map[string]string) {
	for key, value := range src {
		if _, ok := dest[key]; ok {
			logger.Printf("There is already a value present for %s, ignoring new value", key)
		} else {
			dest[key] = value
		}
	}
}

func populateExportVars(exportVars map[string]string, currentVars map[string]string) {
	for key, value := range exportVars {
		if currentValue, ok := currentVars[value]; ok {
			exportVars[key] = currentValue
		} else {
			delete(exportVars, key)
		}
	}
}

func convertToBash(exportVars map[string]string) {
	exports := make([]string, 1)
	for key, value := range exportVars {
		exports = append(exports, "export "+key+"="+value)
	}
	fmt.Println(strings.Join(exports, "\n"))
}

func loadConfigs() {
	allExportedVars := make(map[string]string)
	usr, _ := user.Current()
	homedir := path.Join(usr.HomeDir, "/.ssw/config")
	logger.Println(homedir)
	dirInfo, err := ioutil.ReadDir(homedir)
	configs := make([]string, 0)
	if err == nil {
		for i := range dirInfo {
			if strings.HasSuffix(dirInfo[i].Name(), ".ssw") {
				logger.Println("Found configuration file " + dirInfo[i].Name())
				configs = append(configs, path.Join(homedir, dirInfo[i].Name()))
			}
		}
	} else {
		logger.Fatal("Failed to read directory " + homedir + ". Received error " + err.Error())
	}
	for i := range configs {
		data, err := ioutil.ReadFile(configs[i])
		if err != nil {
			logger.Fatal(err)
		}
		config := new(sswConfig)
		config.Name = strings.TrimSuffix(path.Base(configs[i]), ".ssw")
		if err := config.Parse(data); err != nil {
			logger.Fatal(err)
		}
		if config.ConfigType == "environment" {
			exportVars, err := config.ParseExportVars()
			if err == nil {
				envFile := path.Join(usr.HomeDir, "/.ssw", config.Name, "current-env.ssw")
				logger.Printf("Loading current environment %s", envFile)
				currentVars := make(map[string]string)
				envData, err := ioutil.ReadFile(envFile)
				if err == nil {
					yaml.Unmarshal(envData, currentVars)
					populateExportVars(exportVars, currentVars)
					mergeEnvMap(allExportedVars, exportVars)
				} else {
					logger.Fatal(err)
				}
			}
		}
	}
	convertToBash(allExportedVars)
}

func mkdirP(directories []string) {
	for dir := range directories {
		_, stat_err := os.Stat(directories[dir])
		if stat_err != nil {
			if _, ok := stat_err.(*os.PathError); ok {
				mkdir_err := os.MkdirAll(directories[dir], 0755)
				if mkdir_err != nil {
					logger.Fatal("Received error " + mkdir_err.Error())
					os.Exit(1)
				} else {
					logger.Println("Created directory " + directories[dir])
				}
			} else {
				logger.Fatal("Received error " + stat_err.Error())
				os.Exit(1)
			}
		} else {
			logger.Println("Directory ", directories[dir], "already exists")
		}
	}
}

func main() {
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	var sswCmd = &cobra.Command{
		Use:   "sellsword",
		Short: "Sellsword is a generic command-line tool for switching between application configurations",
		Long:  `Sellsword is a generic tool for switching between application configurations`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please invoke one of Sellsword's subcommands to get started. Type sellsword help for more information")
		},
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of sellsword",
		Long:  `All software has versions. This is Sellsword's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Sellsword version 0.0.1")
		},
	}
	sswCmd.AddCommand(versionCmd)

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize Sellsword",
		Long:  `This command creates the Sellsword directory structure and downloads default configurations from git@github.com:bryanwb/sellsword.git`,
		Run: func(cmd *cobra.Command, args []string) {
			usr, _ := user.Current()
			homedir := path.Join(usr.HomeDir, ".ssw/config")
			awsdir := path.Join(usr.HomeDir, ".ssw/aws")
			chefdir := path.Join(usr.HomeDir, ".ssw/chef")
			mkdirP([]string{homedir, awsdir, chefdir})
		},
	}
	sswCmd.AddCommand(initCmd)

	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Loads current Sellsword configurations",
		Long:  `This command loads all default environment configurations for use by the shell`,
		Run: func(cmd *cobra.Command, args []string) {
			loadConfigs()
		},
	}
	sswCmd.AddCommand(loadCmd)

	sswCmd.Execute()
}
