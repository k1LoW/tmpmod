package fs

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// RenameModule - rename module name in all go.mod, go.sum and go files.
func RenameModule(wd, to string, excludeModfile bool) error {
	m, root, err := ModfileAndGoRoot(wd)
	if err != nil {
		return err
	}
	from := m.Module.Mod.Path
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".mod" && ext != ".sum" {
			return nil
		}
		if excludeModfile && (ext == ".mod" || ext == ".sum") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		bb := bytes.ReplaceAll(b, []byte(from), []byte(to))
		if err := os.WriteFile(path, bb, 0644); err != nil { //nolint:gosec
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// ModfileAndGoRoot - returns the parsed go.mod file and the root directory of the go module.
func ModfileAndGoRoot(wd string) (*modfile.File, string, error) {
	for {
		if wd == filepath.Dir(wd) {
			return nil, "", errors.New("Not a go module")
		}
		if fi, err := os.Stat(filepath.Join(wd, "go.mod")); err != nil || fi.IsDir() {
			wd = filepath.Dir(wd)
			continue
		}
		b, err := os.ReadFile(filepath.Join(wd, "go.mod"))
		if err != nil {
			return nil, "", err
		}
		f, err := modfile.Parse("", b, nil)
		if err != nil {
			return nil, "", err
		}
		return f, wd, nil
	}
}

func CleanupModuleFiles(wd string) error {
	if err := os.RemoveAll(filepath.Join(wd, ".git")); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
		if err := os.Rename(filepath.Join(wd, "go.mod"), filepath.Join(wd, "go.mod.tmpmod")); err != nil {
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(wd, "go.sum")); err == nil {
		if err := os.Rename(filepath.Join(wd, "go.sum"), filepath.Join(wd, "go.sum.tmpmod")); err != nil {
			return err
		}
	}
	return nil
}
