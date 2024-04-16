package walk_test

import (
	"testing"

	"github.com/AWtnb/go-walk"
)

func TestGetChildItem(t *testing.T) {
	d := walk.Dir{All: true, Root: `C:\Personal\gotemp\hoge`}
	d.SetWalkDepth(-1)
	d.SetWalkException("_obsolete")
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
	d := walk.Dir{All: true, Root: `C:\Personal\gotemp\hoge\_obsolete`}
	d.SetWalkDepth(-1)
	d.SetWalkException("_obsolete")
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
	d := walk.Dir{All: true, Root: `C:\Personal\gotemp\hoge`}
	d.SetWalkDepth(-1)
	d.SetWalkException("_obsolete")
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
	d := walk.Dir{All: true, Root: `C:\Personal\gotemp\hoge\_obsolete`}
	d.SetWalkDepth(-1)
	d.SetWalkException("_obsolete")
	found, err := d.GetChildItemWithEverything()
	if err != nil {
		t.Error(err)
		return
	}
	for _, s := range found {
		t.Logf("'%s' was found", s)
	}
}
