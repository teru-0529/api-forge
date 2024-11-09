/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"log"
	"os"

	"github.com/koron/go-dproxy"
	"gopkg.in/yaml.v3"
)

// TITLE: Openapi構造体
type Openapi struct {
	serviceName     string
	server          Server
	ImplementedApis []string
	ReadyApis       []string
	formatVersion   string
	title           string
	description     string
	version         string
	apis            []Api
}

type Api struct {
	path        string
	method      string
	operationId string
	summary     string
	description string
	request     Request
	responses   []Response
	// FIXME:
}

type Request struct {
	paramCount      int
	hasRequestBody  bool
	bodyDescription string
}

type Response struct {
	status          string
	bodyDescription string
}

// FUNCTION: Apiパース
func parseOpenapi(param Setting) (*Openapi, error) {
	log.Printf("parse '%s' openapi.", param.ServiceName)
	openapi := Openapi{serviceName: param.ServiceName, server: param.Server, ImplementedApis: param.ImplementedApis, ReadyApis: param.ReadyApis}

	// PROCESS: openapi.yamlの読込み
	row, err := readOpenapi(param.OpenapiPath)
	if err != nil {
		return nil, err
	}

	proxy := dproxy.New(row)
	// PROCESS: transfer(base)
	formatVersion, _ := proxy.M("openapi").String()
	openapi.formatVersion = formatVersion

	info := proxy.M("info")
	title, _ := info.M("title").String()
	openapi.title = title

	description, _ := info.M("description").String()
	openapi.description = description

	version, _ := info.M("version").String()
	openapi.version = version

	// PROCESS: transfer(path)
	apis, _ := proxy.M("paths").Map()
	ls := []Api{}
	for path, pathItem := range apis {
		items, _ := dproxy.New(pathItem).Map()
		for method, item := range items {
			api := Api{path: path, method: method}
			p := dproxy.New(item)

			operationId, _ := p.M("operationId").String()
			api.operationId = operationId
			summary, _ := p.M("summary").String()
			api.summary = summary
			description, _ := p.M("description").String()
			api.description = description

			// PROCESS: request
			// INFO: リクエストパラメータ(num)
			paramCount := p.Q("parameters").Len()

			hasRequestBody := false
			requestBodyName := ""
			// INFO: リクエストボディ(has,description)
			requestBody, err := p.M("requestBody").M("description").String()
			if err == nil {
				hasRequestBody = true
				requestBodyName = requestBody
			}
			api.request = Request{paramCount: paramCount, hasRequestBody: hasRequestBody, bodyDescription: requestBodyName}

			// PROCESS: response
			// INFO: レスポンス(status,description)
			ress := []Response{}
			res, _ := p.M("responses").Map()
			for status, resItem := range res {
				description, _ := dproxy.New(resItem).M("description").String()

				ress = append(ress, Response{status: status, bodyDescription: description})
			}
			api.responses = ress

			ls = append(ls, api)
		}
	}
	openapi.apis = ls

	return &openapi, nil
}

// FUNCTION: Apiファイル読込み
func readOpenapi(path string) (interface{}, error) {

	// PROCESS: openapi.yamlの読込み
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// PROCESS: パース
	var rowApi interface{}
	err = yaml.Unmarshal([]byte(file), &rowApi)
	if err != nil {
		return nil, err
	}
	return rowApi, nil
}
