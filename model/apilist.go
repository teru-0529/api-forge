/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TITLE: Apis構造体
type Apis struct {
	services []Openapi
}

// TITLE: InputParam構造体
type InputParam struct {
	WorkSpaceId string    `yaml:"kongWorkSpace"`
	Settings    []Setting `yaml:"settings"`
}

type Setting struct {
	ServiceName     string   `yaml:"service"`
	OpenapiPath     string   `yaml:"path"`
	Server          Server   `yaml:"server"`
	ImplementedApis []string `yaml:"implementedApis"`
	ReadyApis       []string `yaml:"readyApis"`
}

type Server struct {
	ProdHost string `yaml:"prodHost"`
	ProdPort int    `yaml:"prodPort"`
	MockHost string `yaml:"mockHost"`
	MockPort int    `yaml:"mockPort"`
	MockBase bool   `yaml:"mockBase"`
}

// FUNCTION: Apisの作成
func New(path string) (*Apis, error) {
	// PROCESS: settingの読込み
	param, err := inputParams(path)
	if err != nil {
		return nil, err
	}

	// PROCESS: settingの読込み
	apis := Apis{}
	for _, setting := range param.Settings {
		api, err := parseOpenapi(setting)
		if err != nil {
			return nil, err
		}
		apis.services = append(apis.services, *api)
	}
	return &apis, nil
}

// FUNCTION: InputParam構造のパース
func inputParams(path string) (*InputParam, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// PROCESS: unmarchal
	var param InputParam
	err = yaml.Unmarshal([]byte(file), &param)
	if err != nil {
		return nil, err
	}
	return &param, nil
}
