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
	ds := []string{EXCEPTION + "/uu", EXCEPTION + "/pp", "aa/bb", "aa/.zz/uu", "aa/cc", "bb/ee"}
	fs := []string{EXCEPTION + "/pp/mm.txt", "aa/.zz/i.txt", "aa/.xx", "bb/_obsolete", "aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	return makeTestTree(root, ds, fs)
}

func testDirPath(name string) string {
	up := os.Getenv("USERPROFILE")
	return filepath.Join(up, "Personal", "gotemp", name)
}

func TestTraverse(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under: %s", tdp)
	var w walk.Walker
	w.Init(tdp, true, -1, EXCEPTION)
	found, err := w.Traverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestTraverseWithDepth(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	depths := []int{0, 1, 2}
	for _, dep := range depths {
		t.Logf("---------------------\ntesting to traverse under: %s (depth: %d)", tdp, dep)
		var w walk.Walker
		w.Init(tdp, true, dep, EXCEPTION)
		found, err := w.Traverse()
		if err != nil {
			t.Error(err)
			return
		}
		for _, s := range found {
			r, _ := filepath.Rel(tdp, s)
			t.Logf("'%s' was found", r)
		}
	}
}

func TestTraverseOnlyDir(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse dirs under: %s", tdp)
	var w walk.Walker
	w.Init(tdp, false, -1, EXCEPTION)
	found, err := w.Traverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestTraverseWithoutException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under without exception: %s", tdp)
	var w walk.Walker
	w.Init(tdp, true, -1, "")
	found, err := w.Traverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestTraverseInException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	r := filepath.Join(tdp, EXCEPTION)
	t.Logf("testing to traverse under exception directory itself: %s", r)
	var w walk.Walker
	w.Init(r, true, -1, EXCEPTION)
	found, err := w.Traverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestEverythingTraverse(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under with Everything: %s", tdp)
	var w walk.Walker
	w.Init(tdp, true, -1, EXCEPTION)
	found, err := w.EverythingTraverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestEverythingTraverseWithDepth(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	depths := []int{0, 1, 2}
	for _, dep := range depths {
		t.Logf("---------------------\ntesting to traverse under: %s (depth: %d)", tdp, dep)
		var w walk.Walker
		w.Init(tdp, true, dep, EXCEPTION)
		found, err := w.EverythingTraverse()
		if err != nil {
			t.Error(err)
			return
		}
		for _, s := range found {
			r, _ := filepath.Rel(tdp, s)
			t.Logf("'%s' was found", r)
		}
	}
}

func TestEverythingTraverseOnlyDir(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse dirs under with Everything: %s", tdp)
	var w walk.Walker
	w.Init(tdp, false, -1, EXCEPTION)
	found, err := w.EverythingTraverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestEverythingTraverseWithoutException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under with Everything without exception: %s", tdp)
	var w walk.Walker
	w.Init(tdp, true, -1, "")
	found, err := w.EverythingTraverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestEverythingTraverseInException(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	r := filepath.Join(tdp, EXCEPTION)
	t.Logf("testing to traverse under exception directory itself with Everything: %s", r)
	var w walk.Walker
	w.Init(r, true, -1, EXCEPTION)
	found, err := w.EverythingTraverse()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		r, _ := filepath.Rel(tdp, s)
		t.Logf("'%s' was found", r)
	}
}

func TestSamePathsFound(t *testing.T) {
	tdp := testDirPath("hoge")
	if err := testTree(tdp); err != nil {
		t.Error(err)
		return
	}
	t.Logf("testing to traverse under: %s", tdp)
	var w walk.Walker
	w.Init(tdp, true, -1, EXCEPTION)
	found1, err := w.Traverse()
	if err != nil {
		t.Error(err)
		return
	}
	found2, err := w.EverythingTraverse()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("checking if Traverse() and EverythingTraverse() returns same slice")
	t.Log("result of Traverse()")
	for _, s := range found1 {
		t.Log(s)
	}
	t.Log("result of EverythingTraverse()")
	for _, s := range found2 {
		t.Log(s)
	}
	for _, s := range found1 {
		if !slices.Contains(found2, s) {
			t.Errorf("%s not found in EverythingTraverse() result", s)
			return
		}
	}
	for _, s := range found2 {
		if !slices.Contains(found1, s) {
			t.Errorf("%s not found in Traverse() result", s)
			return
		}
	}
	t.Log("2 slices are the same")
	t.Log("checking 2 slices are the same order")
	isDifferentOrder := false
	for i, s := range found1 {
		if s != found2[i] {
			isDifferentOrder = true
			t.Logf("index %d item is different", i)
			t.Logf("index %d of Traverse() is %s", i, s)
			t.Logf("index %d of EverythingTraverse() is %s", i, found2[i])
		}
	}
	if !isDifferentOrder {
		t.Log("2 slices are the same order")
	}
}
