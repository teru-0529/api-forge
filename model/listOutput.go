/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/teru-0529/api-forge/store"
)

const PROD = "🟢 production"
const MOCK = "🟡 mock"

// FUNCTION: MDファイルの書き込み
func (apis *Apis) ListMd(path string) error {
	// PROCESS: Fileの取得
	file, cleanup, err := store.NewFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	// PROCESS: 書き込み
	file.WriteString("# API list\n")

	for _, service := range apis.services {
		file.WriteString(fmt.Sprintf("\n## %s(%s)\n\n", service.serviceName, service.description))

		file.WriteString("  | Path | Method | OperationId | Summary | ParamNum | RequestBody | Responses | Status |\n")
		file.WriteString("  |---|---|---|---|--:|---|---|---|\n")

		for _, api := range service.apis {
			status := PROD
			if service.isMock(api.operationId) {
				status = MOCK
			}
			file.WriteString(fmt.Sprintf("  | %s | %s | %s | %s | %d | %s | %s | %s |\n",
				api.path, api.method, api.operationId, api.summary, api.request.paramCount, api.request.bodyName(), api.resNames(), status))

		}
	}
	return nil
}

// FUNCTION: Tsvファイルの書き込み
func (apis *Apis) ListTsv(path string) error {
	// PROCESS: Writerの取得
	writer, cleanup, err := store.NewCsvWriter(path)
	if err != nil {
		return err
	}
	defer cleanup()

	// PROCESS: 書き込み
	defer writer.Flush() //内部バッファのフラッシュは必須
	writer.Write([]string{
		"Service",
		"Name",
		"Path",
		"Method",
		"OperationId",
		"Summary",
		"ParamNum",
		"RequestBody",
		"Responses",
		"Status",
	})
	for _, service := range apis.services {
		for _, api := range service.apis {
			status := PROD
			if service.isMock(api.operationId) {
				status = MOCK
			}
			writer.Write([]string{
				service.serviceName,
				service.description,
				api.path,
				api.method,
				api.operationId,
				api.summary,
				strconv.Itoa(api.request.paramCount),
				api.request.bodyName(),
				api.resNames(),
				status,
			})

		}
	}
	return nil
}

// ----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+

// FUNCTION: response名称(連結)
func (oa *Openapi) isMock(operationId string) bool {
	if slices.Contains(oa.ImplementedApis, operationId) {
		return false
	} else if slices.Contains(oa.ReadyApis, operationId) {
		return true
	} else {
		return oa.server.MockBase
	}
}

// FUNCTION: response名称(連結)
func (api *Api) resNames() string {
	resNames := []string{}
	for _, res := range api.responses {
		resNames = append(resNames, res.bodyName())
	}
	return strings.Join(resNames, ", ")
}

// FUNCTION: requeestBody名称
func (req *Request) bodyName() string {
	if req.hasRequestBody {
		return req.bodyDescription
	} else {
		return "N/A"
	}
}

// FUNCTION: responseBody名称
func (res *Response) bodyName() string {
	if res.status == "default" {
		return res.status
	} else {
		return fmt.Sprintf("%s(%s)", res.status, res.bodyDescription)
	}
}
