package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/k1LoW/tmpmod/fs"
)

// root - returns the root directory of the git repository
func root() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if fi, err := os.Stat(filepath.Join(wd, ".git")); err == nil && fi.IsDir() {
			return wd, nil
		}
		if wd == filepath.Dir(wd) {
			return "", errors.New("Not a git repository")
		}
		wd = filepath.Dir(wd)
	}
}

// Switch - create and switch git branch
func Switch(b string) error {
	root, err := root()
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(root)
	if err != nil {
		fmt.Println(1)
		return err
	}
	bn := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", b))

	headRef, err := r.Head()
	if err != nil {
		fmt.Println(2)
		return err
	}
	ref := plumbing.NewHashReference(bn, headRef.Hash())
	if err := r.Storer.SetReference(ref); err != nil {
		fmt.Println(3)
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		fmt.Println(4)
		return err
	}

	if err := w.Checkout(&git.CheckoutOptions{
		Branch: bn,
		Keep:   true,
	}); err != nil {
		fmt.Println(5)
		return err
	}

	return nil
}

// CommitAll - commit all changes
func CommitAll() (string, error) {
	msg := "Renamed module by tmpmod"
	root, err := root()
	if err != nil {
		return "", err
	}
	r, err := git.PlainOpen(root)
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	if _, err := w.Add("."); err != nil {
		return "", err
	}
	h, err := w.Commit(msg, &git.CommitOptions{})
	if err != nil {
		return "", err
	}
	return h.String(), nil
}

// Clone - clone git repository
func Clone(url, root string) (string, error) {
	tmpdir := filepath.Join(root, ".tmp")
	if _, err := os.Stat(tmpdir); err == nil {
		return "", fmt.Errorf("%s already exists", tmpdir)
	}
	defer os.RemoveAll(tmpdir)
	if !strings.Contains(url, "://") {
		url = fmt.Sprintf("%s://%s", "https", url)
	}
	var branch string
	splitted := strings.Split(url, "@")
	if len(splitted) == 2 {
		url = splitted[0]
		branch = splitted[1]
	}
	opt := &git.CloneOptions{
		URL: url,
	}
	if branch != "" {
		opt.ReferenceName = plumbing.NewBranchReferenceName(branch)
		opt.SingleBranch = true
	}
	_, err := git.PlainClone(tmpdir, false, opt)
	m, _, err := fs.ModfileAndGoRoot(tmpdir)
	if err != nil {
		return "", err
	}
	p := filepath.Join(root, m.Module.Mod.Path)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return "", err
	}
	if err := os.Rename(tmpdir, p); err != nil {
		return "", err
	}
	return p, nil
}
