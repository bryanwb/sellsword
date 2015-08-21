package sellsword

import (
	"os"
	"path"
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

// test that Unlink removes current link

// ListEnvs returns correct list

// Test MakeCurrent unsets export vars

// Test MakeCurrent changes current symlink
