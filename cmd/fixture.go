/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teru-0529/api-forge/model"
)

// fixtureCmd represents the fixture command
var fixtureCmd = &cobra.Command{
	Use:   "fixture",
	Short: "Create test fixture template.",
	Long:  "Create test fixture template.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// PROCESS: APIファイルの読み込み
		apiList, err := model.New(settingFile)
		if err != nil {
			return err
		}

		// PROCESS: Fixture出力
		apiList.Fixture(distDir)

		fmt.Println("***command[fixture] completed.")
		return nil
	},
}

func init() {
}
