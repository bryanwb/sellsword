package sellsword

import (
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"testing"
)

func setUpTest() string {
	dir, _ := os.Getwd()
	tmpdir := path.Join(dir, "test/tmp")
	os.MkdirAll(tmpdir, 0755)
	Logger = logrus.New()
	return tmpdir
}

func TestResolveSymlinkForRealLink(t *testing.T) {
	tmpdir := setUpTest()
	source := path.Join(tmpdir, "source")
	target := path.Join(tmpdir, "target")
	ioutil.WriteFile(source, []byte{}, 0755)
	os.Symlink(source, target)
	actualSource, _ := resolveSymlink(target)
	if actualSource != source {
		t.Errorf("Expected %s but received %s", source, actualSource)
	}
}

func TestResolveSymlinkForFakeLinkFails(t *testing.T) {
	_, err := resolveSymlink("/does/not/exist")
	if _, ok := err.(*os.PathError); !ok {
		t.Error("Expected os.PathError to be raised when resolving nonexistent symlink")
	}
}

func TestExpandPath(t *testing.T) {
	setUpTest()
	var actual string
	var expected string
	actual, _ = expandPath("~/.chef/foo")
	usr, _ := user.Current()
	expected = path.Join(usr.HomeDir, ".chef/foo")
	if actual != expected {
		t.Errorf("Expected %s but received %s", expected, actual)
	}
	wd, _ := os.Getwd()
	expected = path.Join(path.Dir(wd), "foobar")
	actual, _ = expandPath("../foobar")
	if actual != expected {
		t.Errorf("Expected %s but received %s", expected, actual)
	}
}

func TestContains(t *testing.T) {
	a := []string{"foo", "bar", "baz"}
	if !contains(a, "foo") {
		t.Errorf("Expected contains method to report that %v contains %s\n", a, "foo")
	}
	if contains(a, "nope") {
		t.Errorf("Expected contains method to report that %v does not contain %s\n", a, "nope")
	}
}

func TestAppendIfMissing(t *testing.T) {
	var expectedLen int
	var updated []string
	a := []string{"foo", "bar", "baz"}
	expectedLen = len(a)
	updated = appendIfMissing(a, "foo")
	if expectedLen != len(updated) {
		t.Errorf("Expected method to not add duplicate value %s", "foo")
	}
	expectedLen = len(a)
	updated = appendIfMissing(a, "blah")
	if expectedLen == len(updated) {
		t.Errorf("Expected method to add new value %s", "blah")
	}

}

func TestArrayToMap(t *testing.T) {
	myarr := []string{"foo", "bar", "baz"}
	expected := map[string]string{"foo": "", "bar": "", "baz": ""}
	actual := arrayToEmptyMap(myarr)
	for k, v := range expected {
		if v != actual[k] {
			t.Errorf("Expected %s to be present. It appears that %v and %v do not match", v, actual, expected)
		}
	}
}
