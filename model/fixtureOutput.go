/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/teru-0529/api-forge/store"
)

// FUNCTION: Fixtureの書き込み
func (apiList *ApiList) Fixture(dir string) error {

	for _, service := range apiList.Services {
		dirPath := filepath.Join(dir, service.ServiceName)
		for _, api := range service.openapi.apis {
			fileName := fmt.Sprintf("%s_ok.yaml", api.operationId)
			outPath := filepath.Join(dirPath, fileName)

			// PROCESS: Fileの取得
			file, cleanup, err := store.NewFile(outPath)
			if err != nil {
				return err
			}
			defer cleanup()

			// PROCESS: 書き込み
			file.WriteString("dataType: api_test_fixture\n")
			file.WriteString("version: 1.0.0\n")

			file.WriteString("# ----+----+----+\n")
			file.WriteString("# 基本情報\n")
			file.WriteString(fmt.Sprintf("description: %s(成功ケース)\n", api.summary))

			file.WriteString("# ----+----+----+\n")
			file.WriteString("# 情報クリア\n")
			file.WriteString("reset:\n")
			file.WriteString("  sequences:\n")
			file.WriteString("    ## @TODO: 数値を1にリセットするシーケンスをスキーマ毎の配列で指示します。\n")
			file.WriteString("    ## - schema: orders\n")
			file.WriteString("    ##   items: [order_no_seed, product_id_seed]\n")
			file.WriteString("  tables:\n")
			file.WriteString("    ## @TODO: トランケートするテーブルをスキーマ毎の配列で指示します。\n")
			file.WriteString("    ## - schema: orders\n")
			file.WriteString("    ##   items: [products, receivings]\n")

			file.WriteString("# ----+----+----+\n")
			file.WriteString("# 事前処理\n")
			file.WriteString("setupTable:\n")
			file.WriteString("  ## @TODO: テスト実施前に登録するデータをスキーマ,テーブル,データ(JSON)毎の配列で指示します。\n")
			file.WriteString("  ## - schema: orders\n")
			file.WriteString("  ##   table: products\n")
			file.WriteString("  ##   body: '[\n")
			file.WriteString("  ##     {\"product_name\": \"日本刀\",\"cost_price\": 20000},\n")
			file.WriteString("  ##     {\"product_name\": \"火縄銃\",\"cost_price\": 40000},\n")
			file.WriteString("  ##     {\"product_name\": \"弓\",\"cost_price\": 15000},\n")
			file.WriteString("  ##     ]'\n")

			file.WriteString("# ----+----+----+\n")
			file.WriteString("# API実行\n")
			file.WriteString("execute:\n")
			file.WriteString("  ## @TODO: hostKeyを個別に設定します(環境変数で`host:port`の形式)。\n")
			file.WriteString("  hostKey: '@@@@@'\n")
			file.WriteString(fmt.Sprintf("  method: %s\n", strings.ToUpper(api.method)))
			file.WriteString("  ## @TODO: Pathパラメータがある場合は適宜変換します。Queryメータがある場合も設定します。\n")
			file.WriteString(fmt.Sprintf("  path: %s\n", re.ReplaceAllString(api.path, `@@@@@`)))
			file.WriteString("  headers:\n")
			file.WriteString("    - key: x-account-id\n")
			file.WriteString("      ## @TODO: HeaderパラメータとしてアカウントIDを指定します。\n")
			file.WriteString("      value: '@@@@@'\n")
			if api.request.hasBody {
				file.WriteString("  ## @TODO: RequestBosyをJson文字列で指定します。\n")
				file.WriteString("  body: '@@@@@'\n")
			} else {
				file.WriteString("  body: null\n")
			}

			file.WriteString("# ----+----+----+\n")
			file.WriteString("# 検証\n")
			file.WriteString("verification:\n")
			status := api.normalStatus()
			file.WriteString(fmt.Sprintf("  httpStatus: %s\n", status))
			if status == "200" {
				file.WriteString("  execResult:\n")
				file.WriteString("    check: true\n")
				file.WriteString("    ## @TODO: 検証を除外したい項目(例えば日付項目)を配列で指定します。指定方法はJsonPathに従います。($は配列のすべての要素の意。)\n")
				file.WriteString("    ## exclude: [/$/created_at,]\n")
				file.WriteString("    exclude: []\n")
				file.WriteString("  tables: []\n")

			} else {
				file.WriteString("  execResult:\n")
				file.WriteString("    check: false\n")
				file.WriteString("    exclude: []\n")
				file.WriteString("  tables:\n")
				file.WriteString("    ## @TODO: 検証対象にするテーブルをスキーマ,テーブル,検証除外項目の配列で指示します。\n")
				file.WriteString("    ## - schema: orders\n")
				file.WriteString("    ##   table: receiving_details\n")
				file.WriteString("    ##   exclude: [/$/order_date, /$/created_at, /$/updated_at,]\n")
			}
		}
	}
	return nil
}

// FUNCTION: normalStatus
func (api *Api) normalStatus() string {
	for _, res := range api.responses {
		if res.status[:1] == "2" {
			return res.status
		}
	}
	return "default"
}
