/*
Copyright © 2024 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/k1LoW/tmpmod/fs"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

var revertCmd = &cobra.Command{
	Use:   "revert [DIR]",
	Short: "revert modified module",
	Long:  `revert modified module.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
			return fmt.Errorf("directory not found: %s", dir)
		}
		mf := filepath.Join(dir, fs.RenamedModfile)
		if _, err := os.Stat(mf); err != nil {
			return fmt.Errorf("%s is not modified by tmpmod", dir)
		}
		b, err := os.ReadFile(mf)
		if err != nil {
			return err
		}
		f, err := modfile.Parse("", b, nil)
		if err != nil {
			return err
		}
		to := f.Module.Mod.Path
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		m, root, err := fs.ModfileAndGoRoot(wd)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			return err
		}
		from := fmt.Sprintf("%s/%s", m.Module.Mod.Path, rel)
		cmd.Println("Reverting module from", from, "to", to+"...")
		if err := fs.RenameModule(root, from, to, true); err != nil {
			return err
		}

		if err := os.RemoveAll(dir); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Reverted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
