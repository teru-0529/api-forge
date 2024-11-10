/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/teru-0529/api-forge/store"
	"gopkg.in/yaml.v3"
)

// TITLE: ApiList構造体
type ApiList struct {
	WorkSpaceId string    `yaml:"workSpaceId"`
	InitIsMock  bool      `yaml:"initIsMock"`
	Services    []Service `yaml:"services"`
}

type Service struct {
	ServiceName string `yaml:"serviceName"`
	OpenapiPath string `yaml:"openapiPath"`
	openapi     Openapi
	ProdServer  Server   `yaml:"prodServer"`
	MockServer  Server   `yaml:"mockServer"`
	Apis        []ApiKey `yaml:"apis"`
}

type Server struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	ServiceId string `yaml:"serviceId"`
}
type ApiKey struct {
	Title       string `yaml:"title"`
	OperationId string `yaml:"operationId"`
	KongId      string `yaml:"kongId"`
	ResourceId  string `yaml:"resourceId"`
	Implemented bool   `yaml:"implemented"`
}

// FUNCTION: Apisの作成
func New(path string) (*ApiList, error) {
	// PROCESS: settingの読込み
	apiList, err := newApiList(path)
	if err != nil {
		return nil, err
	}

	// PROCESS: serviceの設定
	for i, service := range apiList.Services {
		// PROCESS: serverの設定
		apiList.Services[i].ProdServer.init()
		apiList.Services[i].MockServer.init()

		// PROCESS: openapi読込み
		openapi, err := NewOpenapi(apiList.Services[i])
		if err != nil {
			return nil, err
		}
		apiList.Services[i].openapi = *openapi

		// PROCESS: APIリスト(不足分)設定
		for _, item := range openapi.apis {
			if !service.registered(item.operationId) {
				apiKey := ApiKey{
					Title:       item.summary,
					OperationId: item.operationId,
					KongId:      uuid.NewString(),
					ResourceId:  generateResourceId(service.ServiceName, len(apiList.Services[i].Apis)),
					Implemented: !apiList.InitIsMock,
				}
				apiList.Services[i].Apis = append(apiList.Services[i].Apis, apiKey)
			}
		}
	}

	// PROCESS: 設定ファイル保存
	apiList.Write(path)

	return apiList, nil
}

// FUNCTION: yamlファイルの書き込み
func (apiList *ApiList) Write(path string) error {
	// PROCESS: Encoderの取得
	encoder, cleanup, err := store.NewYamlEncorder(path)
	if err != nil {
		return err
	}
	defer cleanup()
	err = encoder.Encode(&apiList)
	if err != nil {
		return err
	}

	return nil
}

// FUNCTION: API登録済みかどうか
func (service *Service) registered(operationId string) bool {
	for _, api := range service.Apis {
		if api.OperationId == operationId {
			return true
		}
	}
	return false
}

// FUNCTION: APIKeyの取得
func (service *Service) getApikey(operationId string) (*ApiKey, error) {
	for _, api := range service.Apis {
		if api.OperationId == operationId {
			return &api, nil
		}
	}
	return nil, errors.New("Not found")
}

// FUNCTION: ServiceIdの設定
func (server *Server) init() {
	if server.ServiceId == "" {
		server.ServiceId = uuid.NewString()
	}
}

// FUNCTION: リソースID
func generateResourceId(serviceName string, seq int) string {
	var name string
	if len(serviceName) > 6 {
		name = serviceName[:6]
	} else {
		name = serviceName + strings.Repeat("_", 6-len(serviceName))
	}

	return fmt.Sprintf("API-%s-%06d", name, seq+1)
}

// FUNCTION: ApiList構造のパース
func newApiList(path string) (*ApiList, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// PROCESS: unmarchal
	var param ApiList
	err = yaml.Unmarshal([]byte(file), &param)
	if err != nil {
		return nil, err
	}
	return &param, nil
}
