package auth

import (
	"fmt"

	"encoding/json"
	"github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	fmt.Println("Check request: ", req)
	apiKey, err := getSessionFromCookie(req)
	fmt.Println("x-api-key: ", *apiKey)

	if err != nil || *apiKey == "" || *apiKey != "123abc" && *apiKey != "456def" {
		fmt.Println("error in getting apiKey", err)
		return buildDeniedResponse(rpc.UNAUTHENTICATED, "unauthenticated in OPA", "OPA failed to authenticate")
	}

	headerValues := []headerValue{}
	headerValue := headerValue{
		key:   "x-api-key-header",
		value: *apiKey,
	}
	headerValues = append(headerValues, headerValue)

	return buildOkResponse(rpc.OK, "OK", headerValues)
}

type headerValue struct {
	key   string
	value string
}

func buildOkResponse(statusCode rpc.Code, statusMsg string, headerValues []headerValue) (*envoy_service_auth_v3.CheckResponse, error) {
	requestCounts := map[string]string{
		"123abc": "5",
		"456def": "10",
	}

	status := &status.Status{
		Code:    int32(statusCode),
		Message: statusMsg,
	}

	var headerValueOptions []*envoy_config_core_v3.HeaderValueOption

	headerValueOptions = append(headerValueOptions, &envoy_config_core_v3.HeaderValueOption{
		Header: &envoy_config_core_v3.HeaderValue{
			Key:   "x-test-api-header-value",
			Value: headerValues[0].value,
		},
	})

	headerValueOptions = append(headerValueOptions, &envoy_config_core_v3.HeaderValueOption{
		Header: &envoy_config_core_v3.HeaderValue{
			Key:   "x-ext-auth-ratelimit",
			Value: requestCounts[headerValues[0].value],
			// Value: "3",
		},
	})

	headerValueOptions = append(headerValueOptions, &envoy_config_core_v3.HeaderValueOption{
		Header: &envoy_config_core_v3.HeaderValue{
			Key:   "x-ext-auth-ratelimit-unit",
			Value: "MINUTE",
		},
	})

	okHttpResponse := &envoy_service_auth_v3.OkHttpResponse{
		Headers: headerValueOptions,
	}

	okResponse := &envoy_service_auth_v3.CheckResponse_OkResponse{
		OkResponse: okHttpResponse,
	}

	responseHeader := &envoy_service_auth_v3.CheckResponse{
		Status:       status,
		HttpResponse: okResponse,
	}

	return responseHeader, nil
}

func buildDeniedResponse(statusCode rpc.Code, statusMsg string, bodyMsg string) (*envoy_service_auth_v3.CheckResponse, error) {
	status := &status.Status{
		Code:    int32(statusCode),
		Message: statusMsg,
	}

	// deniedStatus := &envoy_type.HttpStatus{
	deniedStatus := &envoy_type_v3.HttpStatus{
		Code: envoy_type_v3.StatusCode_Unauthorized,
	}

	deniedHttpResonse := &envoy_service_auth_v3.DeniedHttpResponse{
		Status:  deniedStatus,
		Headers: nil,
		Body:    bodyMsg,
	}

	deniedResponse := &envoy_service_auth_v3.CheckResponse_DeniedResponse{
		DeniedResponse: deniedHttpResonse,
	}

	responseHeader := &envoy_service_auth_v3.CheckResponse{
		Status:       status,
		HttpResponse: deniedResponse,
	}

	return responseHeader, nil
}

func getSessionFromCookie(req *envoy_service_auth_v3.CheckRequest) (*string, error) {
	b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
	fmt.Println("request headers: ", string(b))
	if err == nil {
		fmt.Println("Inbound Headers: ")
		fmt.Println((string(b)))
	}

	if err != nil {
		fmt.Println("error in getting attributes")
		fmt.Println("Inbound Headers: ")
		fmt.Println((string(b)))
	}

	m := make(map[string]string)
	err = json.Unmarshal(b, &m)

	if err != nil {
		fmt.Println("error in unmarshall header attributes from byte []")
	}

	//todo:  remove
	fmt.Println("unmarshall x-api-key", m)
	x := m["x-api-key"]
	fmt.Println("unmarshall x-api-key value", x)
	return &x, err
}
