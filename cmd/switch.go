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

	"github.com/k1LoW/tmpmod/fs"
	"github.com/k1LoW/tmpmod/git"
	"github.com/spf13/cobra"
)

var as string

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "create and switch git branch for renamed module",
	Long:  `create and switch git branch for renamed module.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if as == "" {
			return fmt.Errorf("module name is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		b := fmt.Sprintf("renamed-%s-by-tmpmod", as)
		cmd.Println("Switching to", b+"...")
		if err := git.Switch(b); err != nil {
			return err
		}
		cmd.Println("Renaming module to", as+"...")
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		m, root, err := fs.ModfileAndGoRoot(wd)
		if err != nil {
			return err
		}
		from := m.Module.Mod.Path
		if err := fs.RenameModule(root, from, as, false); err != nil {
			return err
		}
		hash, err := git.CommitAll()
		if err != nil {
			return err
		}
		cmd.Println("Committed")
		cmd.Println()
		cmd.Printf("Usage: push %s and use `go get %s@%s`\n", b, as, hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringVarP(&as, "as", "", "", "module name to be renamed")
}
