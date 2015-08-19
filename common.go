package sellsword

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
)

var Logger *log.Logger

var Version = "0.0.1"

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
