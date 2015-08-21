package sellsword

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewDirectoryEnv(t *testing.T) {
	tmp := setUpTest()
	e, _ := NewDirectoryEnv("foo", tmp)
	expected := path.Join(tmp, "foo")
	if e.Path != expected {
		t.Errorf("Path for new env is not set correctly, expected %s and got %s\n", expected, e.Path)
	}
	if e.Name != "foo" {
		t.Errorf("Name for new env is not set correctly, expected %s and got %s\n", "foo", e.Name)
	}
	if e.EnvType != "directory" {
		t.Errorf("Type for new env is not set correctly, expected %s and got %s\n", "directory", e.EnvType)
	}
}

func TestNewEnvironmentEnv(t *testing.T) {
	tmp := setUpTest()
	exportVars := map[string]string{"USERNAME": "", "PASSWORD": "", "REGION": ""}
	vars := []string{"username", "password", "region"}
	e, _ := NewEnvironmentEnv("acme", tmp, exportVars, vars)
	expectedPath := path.Join(tmp, "acme")
	if e.Path != expectedPath {
		t.Errorf("Expected path to be %s but found %s", expectedPath, e.Path)
	}
	if e.EnvType != "environment" {
		t.Errorf("Expected environment to be %s, found %s\n", "environment", e.EnvType)
	}
}

func TestPopulateExportVars(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	dir := path.Join(wd, "test/vanilla")
	exportVars := map[string]string{"USERNAME": "username", "PASSWORD": "password", "REGION": "region"}
	exportKeys := []string{"USERNAME", "PASSWORD", "REGION"}
	vars := []string{"username", "password", "region"}
	values := []string{"mcmuffin", "holdthestuffin", "nowhere"}
	e, _ := NewEnvironmentEnv("acme", dir, exportVars, vars)
	e.PopulateExportVars()
	for i := range exportKeys {
		if e.ExportVariables[exportKeys[i]] != values[i] {
			t.Errorf("Expected %s to map to %s, found %s\n", exportKeys[i], values[i],
				e.ExportVariables[exportKeys[i]])
		}
	}
}

func TestPopulateExportVarsRemovesMissingKey(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	dir := path.Join(wd, "test/vanilla")
	exportVars := map[string]string{"USERNAME": "username", "PASSWORD": "password", "REGION": "region",
		"PROFILE": "profile"}
	vars := []string{"username", "password", "region"}
	e, _ := NewEnvironmentEnv("acme", dir, exportVars, vars)
	e.PopulateExportVars()
	if _, ok := e.ExportVariables["PROFILE"]; ok {
		t.Errorf("Expected %s to have been removed because it is not present in actual environment file",
			"PROFILE")
	}
}

func TestPopulateExportVarsNonexistentFile(t *testing.T) {
	setUpTest()
	exportVars := map[string]string{"USERNAME": "username", "PASSWORD": "password", "REGION": "region"}
	vars := []string{"username", "password", "region"}
	e, _ := NewEnvironmentEnv("foobar", "/does/not/exist", exportVars, vars)
	err := e.PopulateExportVars()
	if _, ok := err.(*os.PathError); !ok {
		t.Error("Expected PathError for nonexistent path but no error received")
	}
}

func TestInvalidYaml(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	dir := path.Join(wd, "test/vanilla")
	exportVars := map[string]string{"USERNAME": "username", "PASSWORD": "password", "REGION": "region"}
	vars := []string{"username", "password", "region"}
	e, _ := NewEnvironmentEnv("dyncorp", dir, exportVars, vars)
	err := e.PopulateExportVars()
	if _, ok := err.(error); !ok {
		t.Error("Expected error but did not receive one for invalid yaml")
	}
}

// test printexports
func TestEnvPrintExports(t *testing.T) {
	setUpTest()
	wd, _ := os.Getwd()
	dir := path.Join(wd, "test/vanilla")
	exportVars := map[string]string{"USERNAME": "username", "PASSWORD": "password", "REGION": "region",
		"PROFILE": "profile"}
	vars := []string{"username", "password", "region"}
	e, _ := NewEnvironmentEnv("acme", dir, exportVars, vars)
	e.PopulateExportVars()
	// Using TrimSpace so that extra new lines don't fail this test
	actual := strings.TrimSpace(e.MakeExportStatements())
	expected := strings.TrimSpace(`export PASSWORD=holdthestuffin
export REGION=nowhere
export USERNAME=mcmuffin
`)
	if actual != expected {
		t.Errorf("Expected export statements did not match actual. Actual statements were \n%s\nExpected was %s",
			actual, expected)
	}
}

// test save

// test construct
