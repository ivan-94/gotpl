package gotpl

import (
	"errors"
	"os"
	. "testing"
	"time"
)

const sampleRoot = "./sample"

var (
	accpetedTemplateName = []string{
		"header.html",
		"footer.html",
		"a/index.html",
		"a/about.html",
		"a/nested_a/index.html",
		"b/index.html",
	}
	unacceptedTemplate = []string{
		"body.tpl",
		"a/nested_a/list.tpl",
	}
)

type mockStat struct {
	name    string
	isDir   bool
	modTime time.Time
}

func (m mockStat) Name() string {
	return m.name
}
func (m mockStat) Size() int64 {
	return 0
}
func (m mockStat) Mode() os.FileMode {
	return 0
}
func (m mockStat) ModTime() time.Time {
	return m.modTime
}
func (m mockStat) IsDir() bool {
	return m.isDir
}
func (m mockStat) Sys() interface{} {
	return nil
}
func mockFileInfo(name string, isDir bool, modTime time.Time) os.FileInfo {
	return mockStat{
		name:    name,
		isDir:   isDir,
		modTime: modTime,
	}
}

type walkFuncTable []struct {
	path   string
	info   os.FileInfo
	err    error
	reterr bool
}

func TestSetExt(t *T) {
	tpl := New(sampleRoot)
	if tpl.Ext != ".html" {
		t.Errorf("default ext should be '.html'")
	}
	tpl.SetExt(".tpl")
	if tpl.Ext != ".tpl" {
		t.Errorf("default ext should be '.tpl'")
	}
}

func TestWalkFunc(t *T) {
	tpl := New("/")
	now := time.Now()
	table := walkFuncTable{
		{path: "/a.html", info: mockFileInfo("a.html", false, now), err: nil, reterr: false},
		// conflict
		{path: "/a.html", info: mockFileInfo("a.html", false, now), err: nil, reterr: true},
		{path: "/b.html", info: mockFileInfo("b.html", false, time.Now()), err: nil, reterr: false},
		// conflict
		{path: "/a.html", info: mockFileInfo("a.html", false, now), err: nil, reterr: true},
		{path: "/c", info: mockFileInfo("c", true, now), err: nil, reterr: false},
		{path: "/c/d.html", info: mockFileInfo("c/d.html", false, now), err: nil, reterr: false},
		{path: "/c/e.html", info: mockFileInfo("c/e.html", false, now), err: nil, reterr: false},
		// conflit
		{path: "/c/d.html", info: mockFileInfo("c/d.html", false, now), err: nil, reterr: true},
		{path: "/c/f.html", info: mockFileInfo("c/f.html", false, now), err: errors.New(""), reterr: true},
		{path: "/f.tpl", info: mockFileInfo("f.tpl", false, now), err: nil, reterr: false},
		{path: "/f.tpl.html", info: mockFileInfo("f.tpl.html", false, now), err: nil, reterr: false},
	}

	accpeted := []string{
		"a.html",
		"b.html",
		"c/d.html",
		"c/e.html",
		"f.tpl.html",
	}

	for _, i := range table {
		err := tpl.walkFunc(i.path, i.info, i.err)
		if i.reterr && err == nil || !i.reterr && err != nil {
			t.Errorf("expect walkFunc(%s) return error(%t) but get %s", i.path, i.reterr, err)
		}
	}

	for _, name := range accpeted {
		if _, existed := tpl.files[name]; !existed {
			t.Errorf("expected name '%s' not found in %#v", name, tpl.files)
		}
	}
}

func TestWalk(t *T) {
	checkWalk := func(ext string, accepted, unaccepted []string) {
		tpl := New(sampleRoot)
		tpl.SetExt(ext)
		err := tpl.Walk()
		if err != nil {
			t.Error(err)
		}

		for _, name := range accepted {
			if _, existed := tpl.files[name]; !existed {
				t.Fatalf("expected name '%s' not found in %#v", name, tpl.files)
			}
		}

		for _, name := range unaccepted {
			if _, existed := tpl.files[name]; existed {
				t.Fatalf("unexpected name '%s' found in %#v", name, tpl.files)
			}
		}
	}

	checkWalk(".html", accpetedTemplateName, unacceptedTemplate)
	checkWalk(".tpl", unacceptedTemplate, accpetedTemplateName)
}

func TestParseFiles(t *T) {
	checkparse := func(ext string, accepted []string) {
		tpl := New(sampleRoot)
		tpl.SetExt(ext)
		if err := tpl.Walk(); err != nil {
			t.Error(err)
		}

		if err := tpl.ParseFiles(); err != nil {
			t.Error(err)
		}

		for _, name := range accepted {
			var found bool
			for _, t := range tpl.t.Templates() {
				if name == t.Name() {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("expected template name '%s' not found in assosicated templates %#v", name, tpl.t.Templates())
			}
		}
	}

	checkparse(".html", accpetedTemplateName)
	checkparse(".tpl", unacceptedTemplate)
}

func TestFreshWalkFunc(t *T) {
	old := time.Now().Add(-1 * time.Second)
	now := time.Now()

	tpl := New("/")
	table := walkFuncTable{
		{path: "/a.html", info: mockFileInfo("a.html", false, old), err: nil, reterr: false},
		{path: "/c/d.html", info: mockFileInfo("c/d.html", false, old), err: nil, reterr: false},
		{path: "/c/e.html", info: mockFileInfo("c/e.html", false, old), err: nil, reterr: false},
	}

	for _, i := range table {
		err := tpl.walkFunc(i.path, i.info, i.err)
		if err != nil {
			t.Fatal(err)
		}
	}

	accpeted := []string{
		"a.html",
		"c/d.html",
		"c/e.html",
	}

	for _, name := range accpeted {
		if _, existed := tpl.files[name]; !existed {
			t.Errorf("expected name '%s' not found in %#v", name, tpl.files)
		}
	}

	newItems := walkFuncTable{
		{path: "/b.html", info: mockFileInfo("b.html", false, now), err: nil, reterr: false},
		{path: "/c/d.html", info: mockFileInfo("c/d.html", false, now), err: nil, reterr: false},
		{path: "/c/e.html", info: mockFileInfo("c/e.html", false, old), err: nil, reterr: false},
		{path: "/d/a.html", info: mockFileInfo("d/a.html", false, now), err: nil, reterr: false},
	}

	newFiles := []*tplFile{}
	accpetedNewFiles := []string{
		"b.html",
		"c/d.html",
		"d/a.html",
	}
	walkFunc := tpl.freshWalkFunc(&newFiles)

	for _, i := range newItems {
		err := walkFunc(i.path, i.info, i.err)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, file := range newFiles {
		if file.relPath != accpetedNewFiles[i] {
			t.Fatalf("newItems(%#v) not match (%#v)", newFiles, accpetedNewFiles)
		}
	}
}
