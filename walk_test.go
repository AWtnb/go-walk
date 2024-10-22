package walk_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/AWtnb/go-walk"
)

const EXCEPTION = "_obsolete"

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
	ds := []string{EXCEPTION + "/uu", EXCEPTION + "/pp", "aa/bb", "aa/cc", "bb/ee"}
	fs := []string{EXCEPTION + "/pp/mm.txt", "aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	return makeTestTree(root, ds, fs)
}

func testDirPath(name string) string {
	up := os.Getenv("USERPROFILE")
	return filepath.Join(up, "Personal", "gotemp", name)
}

func TestGetChildItem(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under: %s", tdp)
	var d walk.Dir
	d.Init(tdp, true, -1, EXCEPTION)
	found, err := d.GetChildItem()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemOnlyDir(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse dirs under: %s", tdp)
	var d walk.Dir
	d.Init(tdp, false, -1, EXCEPTION)
	found, err := d.GetChildItem()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemWithoutException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under without exception: %s", tdp)
	var d walk.Dir
	d.Init(tdp, true, -1, "")
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
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	r := filepath.Join(tdp, EXCEPTION)
	t.Logf("testing to traverse under exception directory itself: %s", r)
	var d walk.Dir
	d.Init(r, true, -1, EXCEPTION)
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
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under with Everything: %s", tdp)
	var d walk.Dir
	d.Init(tdp, true, -1, EXCEPTION)
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemWithEverythingOnlyDir(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse dirs under with Everything: %s", tdp)
	var d walk.Dir
	d.Init(tdp, false, -1, EXCEPTION)
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestGetChildItemWithEverythingWithoutException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under with Everything without exception: %s", tdp)
	var d walk.Dir
	d.Init(tdp, true, -1, "")
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
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	r := filepath.Join(tdp, EXCEPTION)
	t.Logf("testing to traverse under exception directory itself with Everything: %s", r)
	var d walk.Dir
	d.Init(r, true, -1, EXCEPTION)
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}

func TestSamePathsFound(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under: %s", tdp)
	var d walk.Dir
	d.Init(tdp, true, -1, EXCEPTION)
	found1, err := d.GetChildItem()
	if err != nil {
		t.Error(err)
		return
	}
	found2, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("checking if GetChildItem() and GetChildItemWithEverything() returns same slice")
	t.Log("result of GetChildItem()")
	for _, s := range found1 {
		t.Log(s)
	}
	t.Log("result of GetChildItemWithEverything()")
	for _, s := range found2 {
		t.Log(s)
	}
	for _, s := range found1 {
		if !slices.Contains(found2, s) {
			t.Errorf("%s not found in GetChildItemWithEverything() result", s)
			return
		}
	}
	for _, s := range found2 {
		if !slices.Contains(found1, s) {
			t.Errorf("%s not found in GetChildItem() result", s)
			return
		}
	}
	t.Log("2 slices are the same")
	t.Log("checking 2 slices are the same order")
	for i, s := range found1 {
		if s != found2[i] {
			t.Logf("index %d item is different", i)
			t.Logf("index %d of GetChildItem() is %s", i, s)
			t.Logf("index %d of GetChildItemWithEverything() is %s", i, found2[i])
		}
	}
}
