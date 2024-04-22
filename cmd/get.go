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
	"github.com/k1LoW/tmpmod/git"
	"github.com/spf13/cobra"
)

var (
	tmpmodRoot string
	all        bool
)

var getCmd = &cobra.Command{
	Use:   "get [REPO]",
	Short: "get renamed module",
	Long:  `get renamed module.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		url := args[0]
		cmd.Println("Getting", url+"...")
		p, hash, err := git.Clone(url, filepath.Join(wd, tmpmodRoot))
		if err != nil {
			return err
		}
		m, root, err := fs.ModfileAndGoRoot(wd)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		as := fmt.Sprintf("%s/%s", m.Module.Mod.Path, rel)
		cmd.Println("Renaming module to", as+"...")
		{
			m, _, err := fs.ModfileAndGoRoot(p)
			if err != nil {
				return err
			}
			from := m.Module.Mod.Path
			rnRoot := p
			if all {
				rnRoot = root
			}
			if err := fs.RenameModule(rnRoot, from, as, true); err != nil {
				return err
			}
		}
		cmd.Println("Cleaning up files...")
		if err := fs.CleanupModuleFiles(p); err != nil {
			return err
		}
		log := fmt.Sprintf("Use %s (%s) as %s temporarily", url, hash, as)
		if err := os.WriteFile(filepath.Join(p, ".tmpmod.log"), []byte(log), os.ModePerm); err != nil {
			return err
		}

		cmd.Println()
		cmd.Printf("Usage: use `%s`\n", as)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&tmpmodRoot, "root", "", "tmpmod", "tmpmod root directory")
	getCmd.Flags().BoolVarP(&all, "rename-all", "", false, "rename also the module path in the importing source codes")
}
