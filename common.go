package sellsword

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

var Logger *log.Logger

var Version = "0.0.2"

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

func resolveSymlink(symlink string) (string, error) {
	fi, err := os.Lstat(symlink)
	if err != nil {
		Logger.Debugf("Path %s does not exist\n", symlink)
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

func contains(l []string, s string) bool {
	for _, str := range l {
		if s == str {
			return true
		}
	}
	return false
}

func appendIfMissing(slice []string, s string) []string {
	if !contains(slice, s) {
		return append(slice, s)
	} else {
		return slice
	}
}
