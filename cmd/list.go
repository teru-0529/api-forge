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
		apiList, err := model.New(settingFile)
		if err != nil {
			return err
		}

		// PROCESS: リスト出力
		apiList.ListMd(filepath.Join(distDir, "api-list.md"))
		apiList.ListTsv(filepath.Join(distDir, "api-list.tsv"))

		fmt.Println("***command[list] completed.")
		return nil
	},
}

func init() {
}
