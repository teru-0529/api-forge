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
	formatVersion string
	title         string
	description   string
	version       string
	apis          []Api
}

type Api struct {
	path        string
	method      string
	operationId string
	summary     string
	description string
	request     Request
	responses   []Response
}

type Request struct {
	paramCount int
	hasBody    bool
	name       string
}

type Response struct {
	status string
	name   string
}

// FUNCTION: Apiパース
func NewOpenapi(service Service) (*Openapi, error) {
	log.Printf("parse '%s' openapi file.", service.ServiceName)
	openapi := Openapi{}

	// PROCESS: openapi.yamlの読込み
	row, err := readOpenapi(service.OpenapiPath)
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

			hasBody := false
			name := ""
			// INFO: リクエストボディ(has,description)
			bodyDescription, err := p.M("requestBody").M("description").String()
			if err == nil {
				hasBody = true
				name = bodyDescription
			}
			api.request = Request{paramCount: paramCount, hasBody: hasBody, name: name}

			// PROCESS: response
			// INFO: レスポンス(status,description)
			ress := []Response{}
			res, _ := p.M("responses").Map()
			for status, resItem := range res {
				description, _ := dproxy.New(resItem).M("description").String()

				ress = append(ress, Response{status: status, name: description})
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
