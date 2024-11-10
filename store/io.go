/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package store

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FUNCTION: 書き込みファイルの準備（フォルダが無ければ作成する）
func NewFile(fileName string) (*os.File, func(), error) {
	dir := filepath.Dir(fileName)

	// PROCESS: フォルダが存在しない場合作成する
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, nil, fmt.Errorf("cannot create directory: %s", err.Error())
		}
	}

	// PROCESS: 出力用ファイルのオープン
	file, err := os.Create(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create file: %s", err.Error())
	}
	return file, func() { file.Close() }, nil
}

// FUNCTION: csv用のライター
func NewCsvWriter(fileName string) (*csv.Writer, func(), error) {
	// PROCESS: ↑のファイル準備を呼び出す
	file, cleanup, err := NewFile(fileName)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	// PROCESS: Excelで文字化けしないようにする設定。BOM付きUTF8をfileの先頭に付与
	buf := bufio.NewWriter(file)
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	// PROCESS: tsv形式でデータを書き込み
	writer := csv.NewWriter(buf)
	writer.Comma = '\t' //タブ区切りに変更

	return writer, cleanup, nil
}

// FUNCTION: yaml用のエンコーダー
func NewYamlEncorder(fileName string) (*yaml.Encoder, func(), error) {
	// PROCESS: ↑のファイル準備を呼び出す
	file, cleanup, err := NewFile(fileName)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	// PROCESS: encode
	yamlEncoder := yaml.NewEncoder(file)
	yamlEncoder.SetIndent(2) // this is what you're looking for
	return yamlEncoder, cleanup, nil
}
