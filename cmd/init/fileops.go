package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/oops"
)

// projectRoot anchors all file I/O to the working directory.
type projectRoot struct {
	dir string
}

func newProjectRoot() (projectRoot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return projectRoot{}, oops.Wrapf(err, "get working directory")
	}

	return projectRoot{dir: cwd}, nil
}

func (pr projectRoot) abs(rel string) string {
	return filepath.Join(pr.dir, filepath.Clean(rel))
}

func (pr projectRoot) readFile(rel string) ([]byte, error) {
	data, err := os.ReadFile(pr.abs(rel))
	if err != nil {
		return nil, oops.Wrapf(err, "read %s", rel)
	}

	return data, nil
}

func (pr projectRoot) writeFile(rel string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(pr.abs(rel), data, perm); err != nil {
		return oops.Wrapf(err, "write %s", rel)
	}

	return nil
}

func (pr projectRoot) replaceInFile(fpath, old, replacement string) error {
	data, readErr := pr.readFile(fpath)
	if readErr != nil {
		return readErr
	}

	content := string(data)
	if !strings.Contains(content, old) {
		return nil
	}

	info, statErr := os.Stat(pr.abs(fpath))
	if statErr != nil {
		return oops.Wrapf(statErr, "stat %s", fpath)
	}

	return pr.writeFile(fpath, []byte(strings.ReplaceAll(content, old, replacement)), info.Mode())
}
