/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/teru-0529/api-forge/store"
)

const TRACE_ID = "SYS_SETUP"

// Pathパラメータにヒットする正規表現
var re = regexp.MustCompile(`\{[^}]*\}`)

// FUNCTION: Kong用SQLの書き込み
func (apiList *ApiList) Sql4Kong(path string) error {
	// PROCESS: Fileの取得
	file, cleanup, err := store.NewFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	// PROCESS: 書き込み
	file.WriteString("-- # service and route data for kong.\n")

	file.WriteString("\n-- ----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+\n\n")
	file.WriteString("-- ## delete tables\n")
	file.WriteString("DELETE FROM route;\n")
	file.WriteString("DELETE FROM service;\n")

	for _, service := range apiList.Services {
		file.WriteString("\n-- ----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+\n\n")
		file.WriteString(fmt.Sprintf("-- ## %s(%s)\n", service.ServiceName, service.openapi.description))

		file.WriteString("\n-- ### Service\n")
		// prod service
		file.WriteString(fmt.Sprintf("INSERT INTO service VALUES (%s);\n", serviceParam(
			service.ProdServer.ServiceId,
			service.openapi.description,
			service.ProdServer.Host,
			service.ProdServer.Port,
			fmt.Sprintf("'%s'", service.ServiceName),
			apiList.WorkSpaceId),
		))
		// mock service
		file.WriteString(fmt.Sprintf("INSERT INTO service VALUES (%s);\n", serviceParam(
			service.MockServer.ServiceId,
			fmt.Sprintf("%s(MOCK)", service.openapi.description),
			service.MockServer.Host,
			service.MockServer.Port,
			fmt.Sprintf("'%s', 'mock'", service.ServiceName),
			apiList.WorkSpaceId),
		))

		file.WriteString("\n-- ### Route\n")
		for _, api := range service.openapi.apis {
			// ApiKeyの取得
			apiKey, err := service.getApikey(api.operationId)
			if err != nil {
				return err
			}
			// Production/Mock
			serverId := service.MockServer.ServiceId
			tag := fmt.Sprintf("'%s', 'mock'", service.ServiceName)
			msg := "-- ★★MOCK★★"
			if apiKey.Implemented {
				serverId = service.ProdServer.ServiceId
				tag = fmt.Sprintf("'%s'", service.ServiceName)
				msg = ""
			}

			file.WriteString(fmt.Sprintf("INSERT INTO route VALUES (%s); %s\n", routeParams(
				service.ServiceName,
				apiKey.KongId,
				api.summary,
				api.operationId,
				serverId,
				strings.ToUpper(api.method),
				api.path,
				tag,
				apiList.WorkSpaceId,
			), msg))
		}
	}
	return nil
}

// FUNCTION: Acl用SQLの書き込み
func (apiList *ApiList) Sql4Acl(path string) error {
	// PROCESS: Fileの取得
	file, cleanup, err := store.NewFile(path)
	if err != nil {
		return err
	}
	defer cleanup()

	// PROCESS: 書き込み
	file.WriteString("-- # api resource data for acl.\n")

	file.WriteString("\n-- ----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+\n\n")
	file.WriteString("-- ## delete tables\n")
	file.WriteString("DELETE FROM acl.api_resources;\n")
	file.WriteString("DELETE FROM acl.resources WHERE resource_type = 'ACL';\n")

	for _, service := range apiList.Services {
		file.WriteString("\n-- ----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+\n\n")
		file.WriteString(fmt.Sprintf("-- ## %s(%s)\n", service.ServiceName, service.openapi.description))

		file.WriteString("\n-- ### Resources / ApiResources\n")
		for _, api := range service.openapi.apis {
			// ApiKeyの取得
			apiKey, err := service.getApikey(api.operationId)
			if err != nil {
				return err
			}

			file.WriteString(fmt.Sprintf("INSERT INTO acl.resources VALUES (%s);\n",
				resourcesParam(apiKey.ResourceId, fmt.Sprintf("%s(%s)", api.summary, api.operationId))))
			file.WriteString(fmt.Sprintf("INSERT INTO acl.api_resources VALUES (%s);\n",
				apiResourcesParam(apiKey.ResourceId, apiKey.KongId)))
		}
	}

	return nil
}

// FUNCTION: serviceParams
func serviceParam(serviceId string, description string, host string, port int, tag string, wsId string) string {
	return fmt.Sprintf("'%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', 5, 'http', '%s', %d, null, 60000, 60000, 60000, ARRAY[%s], null, null, null, null, '%s', true",
		serviceId,
		description,
		host,
		port,
		tag,
		wsId,
	)
}

// FUNCTION: routeParams
func routeParams(serviceName string, kongId string, summary string, operationId string, serviceId string, method string, path string, tag string, wsId string) string {
	return fmt.Sprintf("'%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s', ARRAY['http', 'https'], ARRAY['%s'], null, '(\"%s\")', null, null, null, 0, false, false, ARRAY[%s], 426, null, 'v0', '%s', true, true, null, null",
		kongId,
		fmt.Sprintf("%s(%s)", summary, operationId),
		serviceId,
		method,
		fmt.Sprintf("~/%s%s", serviceName, re.ReplaceAllString(path, `[A-Za-z0-9_-]+`)),
		tag,
		wsId,
	)
}

// FUNCTION: resourcesParam
func resourcesParam(recourceId string, apiName string) string {
	return fmt.Sprintf("'%s', 'API', '%s', %s",
		recourceId,
		apiName,
		traceColumns(),
	)
}

// FUNCTION: api_resourcesParam
func apiResourcesParam(recourceId string, kongId string) string {
	return fmt.Sprintf("'%s', '%s', %s",
		recourceId,
		kongId,
		traceColumns(),
	)
}

// FUNCTION: traceColumns
func traceColumns() string {
	return fmt.Sprintf("CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s'",
		TRACE_ID,
		TRACE_ID,
	)
}
