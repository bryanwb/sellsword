package sellsword

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

// create new app and check values
func TestNewApp(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	home := path.Join(wd, "test")
	a, _ := NewApp("aws", home)
	expectedDefinition := path.Join(home, "config/aws.ssw")
	if a.Definition != expectedDefinition {
		t.Errorf("Expected definition file to be %s, found %s", expectedDefinition, a.Definition)
	}
	if a.EnvType != "environment" {
		t.Errorf("Expected envType to be %s, found %s", "environment", a.EnvType)
	}
	expectedPath := path.Join(home, "aws")
	if a.Path != expectedPath {
		t.Errorf("Expected path to be %s, found %s", a.Path, expectedPath)
	}
}

func TestNewAppDoesNotExist(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	home := path.Join(wd, "test")
	msg := "Expected os.PathError error for app definition that does not exist but did not receive one"
	if _, err := NewApp("doesnotexist", home); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			t.Error(msg)
		}
	} else {
		t.Error(msg)
	}
}

// test parsing export vars
func TestNewAppParsingExportVars(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	home := path.Join(wd, "test")
	expectedExports := map[string]string{"AWS_ACCESS_KEY_ID": "access_key", "AWS_ACCESS_ID": "access_key",
		"AWS_SECRET_ACCESS_KEY": "secret_key", "AWS_SECRET_KEY": "secret_key",
		"AWS_DEFAULT_REGION": "region", "AWS_REGION": "region"}
	a, _ := NewApp("aws", home)
	for k, v := range expectedExports {
		if _, ok := a.ExportVariables[k]; !ok {
			t.Errorf("Expected exported variables to contain %s but not found", k)
		} else if a.ExportVariables[k] != v {
			t.Errorf("Expected exported variable key to have value %s, found %s", v, a.ExportVariables[k])
		}
	}
}

// test that Current returns correct link
func TestAppCheckCurrent(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	source := path.Join(wd, "test/aws/acme")
	currentLink := path.Join(wd, "test/aws/current")
	relink(source, currentLink)
	a, _ := NewApp("aws", path.Join(wd, "test"))
	e, _ := a.Current()
	if e.Path != source {
		t.Errorf("Expected current env to point to %s, found %s", source, e.Path)
	}
}

// test that Current returns correct link for a Directory environment
func TestAppDirectoryCheckCurrent(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	source := path.Join(wd, "test/chef/acme")
	currentLink := path.Join(wd, "test/chef/current")
	target := path.Join(wd, "fixtures/.chef")
	relink(source, currentLink)
	relink(source, target)
	a, _ := NewApp("chef", path.Join(wd, "test"))
	e, _ := a.Current()
	if e.Path != source {
		t.Errorf("Expected current env to point to %s, found %s", source, e.Path)
	}
	actual, _ := os.Readlink(a.Target)
	if actual != source {
		t.Errorf("Expected current env to point to %s, found %s", source, actual)
	}
}

// test that Unlink removes current link and target directory
func TestAppUnlink(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	source := path.Join(wd, "test/chef/acme")
	currentLink := path.Join(wd, "test/chef/current")
	target := path.Join(wd, "fixtures/.chef")
	relink(source, currentLink)
	relink(source, target)
	a, _ := NewApp("chef", path.Join(wd, "test"))
	a.Unlink()
	_, err1 := os.Lstat(currentLink)
	_, err2 := os.Lstat(target)
	if err1 == nil || err2 == nil {
		t.Errorf("Expected %s and %s to be deleted but either 1 or both were not deleted", currentLink, target)
	}
}

func TestAppUnsetVars(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	a, _ := NewApp("aws", path.Join(wd, "test"))
	unsets := strings.TrimSpace(a.MakeUnsetExportVars())
	expected := strings.TrimSpace(`unset AWS_ACCESS_ID
unset AWS_ACCESS_KEY_ID
unset AWS_DEFAULT_REGION
unset AWS_REGION
unset AWS_SECRET_ACCESS_KEY
unset AWS_SECRET_KEY
`)
	if unsets != expected {
		t.Errorf("Expected %s and found %s", expected, unsets)
	}
}

// ListEnvs returns correct list
func TestListEnvs(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	a, _ := NewApp("aws", path.Join(wd, "test"))
	envs := a.ListEnvs()
	envNames := make([]string, 0)
	expected := []string{"acme", "dyncorp"}
	for i := range envs {
		envNames = append(envNames, envs[i].Name)
	}
	for i := range envNames {
		if envNames[i] != expected[i] {
			t.Errorf("List of envs did not match expected %v, found %v", expected, envNames)
		}
	}
}

// Test MakeCurrent changes current symlink and target symlink
func TestAppMakeCurrent(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	a, _ := NewApp("aws", path.Join(wd, "test"))
	dyncorpPath := path.Join(wd, "test/aws/dyncorp")
	a.MakeCurrent("dyncorp")
	source, _ := os.Readlink(path.Join(wd, "test/aws/current"))
	if source != dyncorpPath {
		t.Errorf("Expected make current to set source to %s, found %s", source, dyncorpPath)
	}
}

// test error cases for MakeCurrent

func TestAppLoadAction(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	output := path.Join(wd, "test/tmp/current")
	os.Remove(output)
	a, _ := NewApp("ssh", path.Join(wd, "test"))
	expected := path.Join(wd, "test/ssh/personal")
	current := path.Join(wd, "test/ssh/current")
	os.Symlink(expected, current)
	e, _ := a.Current()
	a.Load()
	d, _ := ioutil.ReadFile(output)
	actual := strings.TrimSpace(string(d))
	if expected != e.Path {
		t.Errorf("Unload action for App failed or did not run, expected value %s, found %s",
			expected, actual)
	}
}

func TestAppUnloadAction(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	output := path.Join(wd, "test/tmp/current")
	os.Remove(output)
	a, _ := NewApp("ssh", path.Join(wd, "test"))
	expected := path.Join(wd, "test/ssh/personal")
	current := path.Join(wd, "test/ssh/current")
	os.Symlink(expected, current)
	os.Remove(output)
	a.Unload()
	d, _ := ioutil.ReadFile(output)
	actual := strings.TrimSpace(string(d))
	if expected != actual {
		t.Errorf("Unload action for App failed or did not run, expected value %s, found %s",
			expected, actual)
	}
}

// test load action

// Delete the current symlink, which points who knows where, and link it
// to source
func relink(source string, link string) {
	os.Remove(link)
	os.Symlink(source, link)
}
