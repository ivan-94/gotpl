// Package gotpl produces an useful template loader and some helpers to
// enhance html/template development experience.
// see https://github.com/carney520/gotpl for more information.
package gotpl

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type tplFile struct {
	path    string
	relPath string
	modTime time.Time
}

// Tpl is a wrapper of "html/template" , produces an useful loader and some helpers
type Tpl struct {
	// Root is templates root directory
	Root string
	// Ext is template file extension, default is ".html"
	Ext string
	// mux to protect files
	mux   sync.Mutex
	files map[string]*tplFile
	// main template
	t *template.Template
}

func (t *Tpl) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() || filepath.Ext(path) != t.Ext {
		return nil
	}

	relPath, err := filepath.Rel(t.Root, path)

	if err != nil {
		return fmt.Errorf("failed to get relative path for '%s'", path)
	}

	_, existed := t.files[relPath]
	if existed {
		return fmt.Errorf("template '%s' existed", relPath)
	}

	t.files[relPath] = &tplFile{
		path:    path,
		relPath: relPath,
		modTime: info.ModTime(),
	}

	return nil
}

// Walk template root and get all templates info
func (t *Tpl) Walk() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	if len(t.files) != 0 {
		t.files = make(map[string]*tplFile)
	}
	return filepath.Walk(t.Root, t.walkFunc)
}

// ParseFiles 解析文件
func (t *Tpl) ParseFiles() error {
	if len(t.files) == 0 {
		return errors.New("files empty, you may call tpl.Walk() before ParaseFiles")
	}

	for tpPath, fileInfo := range t.files {
		data, err := ioutil.ReadFile(fileInfo.path)
		if err != nil {
			return err
		}
		template := t.t.New(tpPath)
		if _, err := template.Parse(string(data)); err != nil {
			return err
		}
	}
	return nil
}

// Load 加载模板文件并解析
func (t *Tpl) Load() error {
	if err := t.Walk(); err != nil {
		return err
	}
	return t.ParseFiles()
}

func (t *Tpl) freshWalkFunc(files *[]*tplFile) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != t.Ext {
			return nil
		}
		relPath, err := filepath.Rel(t.Root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for '%s'", path)
		}

		file, existed := t.files[relPath]
		modTime := info.ModTime()
		if existed && file.modTime == modTime {
			// not modified
			return nil
		}

		// file update
		*files = append(*files, &tplFile{
			path:    path,
			relPath: relPath,
			modTime: modTime,
		})
		return nil
	}
}

func (t *Tpl) walkFreshFiles() ([]*tplFile, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	files := []*tplFile{}
	err := filepath.Walk(t.Root, t.freshWalkFunc(&files))

	if err != nil {
		return nil, err
	}
	return files, nil
}

// Reload walks template root and add new templates or update modified templates
// it's useful in development
func (t *Tpl) Reload() error {
	files, err := t.walkFreshFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file.path)
		if err != nil {
			return err
		}
		template := t.t.New(file.relPath)
		if _, err := template.Parse(string(data)); err != nil {
			return err
		}
		t.files[file.relPath] = file
	}

	return nil
}

// Template get underlying template.Template
func (t *Tpl) Template() *template.Template {
	return t.t
}

// Funcs adds the elements of the argument map to the template's function map
func (t *Tpl) Funcs(funcMap template.FuncMap) *Tpl {
	t.t.Funcs(funcMap)
	return t
}

// SetExt set template file extension
func (t *Tpl) SetExt(ext string) {
	t.Ext = ext
}

// InstallHelpers intsall helpers from github.com/Masterminds/sprig
// and builin helpers
func (t *Tpl) InstallHelpers() {
}

// New create a new Tpl
func New(root string) *Tpl {
	tpl := &Tpl{
		Root:  root,
		Ext:   ".html",
		files: make(map[string]*tplFile),
		t:     template.New(""),
	}
	tpl.InstallHelpers()
	return tpl
}
