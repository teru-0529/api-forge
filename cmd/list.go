/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/teru-0529/api-forge/model"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Create API list.",
	Long:  "Create API list.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// PROCESS: APIファイルの読み込み
		apis, err := model.New(settingFile)
		if err != nil {
			return err
		}

		// PROCESS:
		apis.ListMd(filepath.Join(distDir, "api-list.md"))
		apis.ListTsv(filepath.Join(distDir, "api-list.tsv"))

		fmt.Println("***command[list] completed.")
		return nil
	},
}

func init() {
}
