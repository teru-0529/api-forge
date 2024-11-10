/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/teru-0529/api-forge/store"
)

const PROD = "🟢 production"
const MOCK = "🟡 mock"

// FUNCTION: MDファイルの書き込み
func (apiList *ApiList) ListMd(path string) error {
	// PROCESS: Fileの取得
	file, cleanup, err := store.NewFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	// PROCESS: 書き込み
	file.WriteString("# API list\n")

	for _, service := range apiList.Services {
		file.WriteString(fmt.Sprintf("\n## %s(%s)\n\n", service.ServiceName, service.openapi.description))

		file.WriteString("  | ResourceId | Path | Method | Name | ParamNum | RequestBody | Responses | Status |\n")
		file.WriteString("  |---|---|---|---|--:|---|---|---|\n")

		for _, api := range service.openapi.apis {
			// ApiKeyの取得
			apiKey, err := service.getApikey(api.operationId)
			if err != nil {
				return err
			}
			// Production/Mock
			status := MOCK
			if apiKey.Implemented {
				status = PROD
			}

			file.WriteString(fmt.Sprintf("  | %s | %s | %s | %s | %d | %s | %s | %s |\n",
				apiKey.ResourceId,
				api.path,
				api.method,
				fmt.Sprintf("%s(%s)", api.summary, api.operationId),
				api.request.paramCount,
				api.request.bodyName(),
				api.resNames(),
				status,
			))
		}
	}
	return nil
}

// FUNCTION: Tsvファイルの書き込み
func (apiList *ApiList) ListTsv(path string) error {
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
		"ResourceId",
		"Path",
		"Method",
		"OperationId",
		"Summary",
		"ParamNum",
		"RequestBody",
		"Responses",
		"Status",
	})
	for _, service := range apiList.Services {
		for _, api := range service.openapi.apis {
			// ApiKeyの取得
			apiKey, err := service.getApikey(api.operationId)
			if err != nil {
				return err
			}
			// Production/Mock
			status := MOCK
			if apiKey.Implemented {
				status = PROD
			}

			writer.Write([]string{
				service.ServiceName,
				service.openapi.description,
				apiKey.ResourceId,
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
	if req.hasBody {
		return req.name
	} else {
		return "N/A"
	}
}

// FUNCTION: responseBody名称
func (res *Response) bodyName() string {
	if res.status == "default" {
		return res.status
	} else {
		return fmt.Sprintf("%s(%s)", res.status, res.name)
	}
}
