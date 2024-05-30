package walk_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWtnb/go-walk"
)

func getPerm(path string) fs.FileMode {
	s := string(os.PathSeparator)
	elems := strings.Split(path, s)
	for i := 0; i < len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], s)
		if fs, err := os.Stat(p); err == nil {
			return fs.Mode() & os.ModePerm
		}

	}
	return 0700
}

func makeTestDir(path string) error {
	p := getPerm(path)
	err := os.MkdirAll(path, p)
	return err
}

func makeTestFile(path string) error {
	_, err := os.Create(path)
	return err
}

func makeTestTree(root string, dirs []string, files []string) error {
	if err := makeTestDir(root); err != nil {
		return err
	}
	for _, d := range dirs {
		p := filepath.Join(root, d)
		if err := makeTestDir(p); err != nil {
			return err
		}
	}
	for _, f := range files {
		p := filepath.Join(root, f)
		if err := makeTestFile(p); err != nil {
			return err
		}
	}
	return nil
}

func testTree(root string) error {
	ds := []string{"_obsolete", "aa/bb", "aa/cc", "bb/ee"}
	fs := []string{"aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	return makeTestTree(root, ds, fs)
}

func TestGetChildItem(t *testing.T) {
	if err := testTree(`C:\Personal\gotemp\hoge`); err != nil {
		t.Error(err)
		return
	}
	var d walk.Dir
	d.Init(`C:\Personal\gotemp\hoge`, true, -1, "_obsolete")
	found, err := d.GetChildItem()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemInException(t *testing.T) {
	if err := testTree(`C:\Personal\gotemp\hoge`); err != nil {
		t.Error(err)
		return
	}
	var d walk.Dir
	d.Init(`C:\Personal\gotemp\hoge\_obsolete`, true, -1, "_obsolete")
	found, err := d.GetChildItem()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemWithEverything(t *testing.T) {
	if err := testTree(`C:\Personal\gotemp\hoge`); err != nil {
		t.Error(err)
		return
	}
	var d walk.Dir
	d.Init(`C:\Personal\gotemp\hoge`, true, -1, "_obsolete")
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemWithEverythingInException(t *testing.T) {
	if err := testTree(`C:\Personal\gotemp\hoge`); err != nil {
		t.Error(err)
		return
	}
	var d walk.Dir
	d.Init(`C:\Personal\gotemp\hoge\_obsolete`, true, -1, "_obsolete")
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}
