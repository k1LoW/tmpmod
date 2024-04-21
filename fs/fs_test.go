package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otiai10/copy"
	"github.com/tenntenn/golden"
)

func TestRenameModule(t *testing.T) {
	tests := []struct {
		excludeModfile bool
		golden         string
	}{
		{false, "rename-a"},
		{true, "rename-a-exclude-modfile"},
	}
	for _, tt := range tests {
		t.Run(tt.golden, func(t *testing.T) {
			dir := t.TempDir()
			if err := copy.Copy("testdata/a", dir); err != nil {
				t.Fatal(err)
			}

			if err := RenameModule(dir, "github.com/my/a", tt.excludeModfile); err != nil {
				t.Fatal(err)
			}

			got := golden.Txtar(t, dir)
			update := false
			if os.Getenv("UPDATE_GOLDEN") != "" {
				update = true
			}
			if diff := golden.Check(t, update, "testdata", tt.golden, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestModfileAndGoRoot(t *testing.T) {
	tests := []struct {
		dir  string
		want string
	}{
		{".", "github.com/k1LoW/tmpmod"},
		{"testdata", "github.com/k1LoW/tmpmod"},
		{"testdata/a", "github.com/k1LoW/tmpmod/fs/testdata/a"},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			abs, err := filepath.Abs(tt.dir)
			if err != nil {
				t.Fatal(err)
			}
			m, _, err := ModfileAndGoRoot(abs)
			if err != nil {
				t.Fatal(err)
			}
			got := m.Module.Mod.Path
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
