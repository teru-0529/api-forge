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

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Create Kong and ACL insert data.",
	Long:  "Create Kong and ACL insert data.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// PROCESS: APIファイルの読み込み
		apiList, err := model.New(settingFile)
		if err != nil {
			return err
		}

		// PROCESS: SQL出力
		apiList.Sql4Kong(filepath.Join(distDir, "kongData.sql"))
		apiList.Sql4Acl(filepath.Join(distDir, "aclData.sql"))

		fmt.Println("***command[sql] completed.")
		return nil
	},
}

func init() {
}
