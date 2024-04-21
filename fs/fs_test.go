package fs

import (
	"os"
	"testing"

	"github.com/otiai10/copy"
	"github.com/tenntenn/golden"
)

func TestRenameModule(t *testing.T) {
	dir := t.TempDir()

	if err := copy.Copy("testdata/a", dir); err != nil {
		t.Fatal(err)
	}

	if err := RenameModule(dir, "github.com/my/a", true); err != nil {
		t.Fatal(err)
	}

	got := golden.Txtar(t, dir)
	update := false
	if os.Getenv("UPDATE_GOLDEN") != "" {
		update = true
	}
	if diff := golden.Check(t, update, "testdata", "rename-a", got); diff != "" {
		t.Error(diff)
	}
}
