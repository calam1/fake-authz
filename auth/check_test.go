package auth

import (
	"crypto/rsa"
	"fmt"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_service_auth_v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/status"
	"grainger.com/auth_proxy/v1/jwt/exchange"
	"math/big"
	"reflect"
	"testing"
)


func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func getMockPrivateKey() *rsa.PrivateKey {
	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: fromBase10("14314132931241006650998084889274020608918049032671858325988396851334124245188214251956198731333464217832226406088020736932173064754214329009979944037640912127943488972644697423190955557435910767690712778463524983667852819010259499695177313115447116110358524558307947613422897787329221478860907963827160223559690523660574329011927531289655711860504630573766609239332569210831325633840174683944553667352219670930408593321661375473885147973879086994006440025257225431977751512374815915392249179976902953721486040787792801849818254465486633791826766873076617116727073077821584676715609985777563958286637185868165868520557"),
			E: 3,
		},
		D: fromBase10("9542755287494004433998723259516013739278699355114572217325597900889416163458809501304132487555642811888150937392013824621448709836142886006653296025093941418628992648429798282127303704957273845127141852309016655778568546006839666463451542076964744073572349705538631742281931858219480985907271975884773482372966847639853897890615456605598071088189838676728836833012254065983259638538107719766738032720239892094196108713378822882383694456030043492571063441943847195939549773271694647657549658603365629458610273821292232646334717612674519997533901052790334279661754176490593041941863932308687197618671528035670452762731"),
		Primes: []*big.Int{
			fromBase10("130903255182996722426771613606077755295583329135067340152947172868415809027537376306193179624298874215608270802054347609836776473930072411958753044562214537013874103802006369634761074377213995983876788718033850153719421695468704276694983032644416930879093914927146648402139231293035971427838068945045019075433"),
			fromBase10("109348945610485453577574767652527472924289229538286649661240938988020367005475727988253438647560958573506159449538793540472829815903949343191091817779240101054552748665267574271163617694640513549693841337820602726596756351006149518830932261246698766355347898158548465400674856021497190430791824869615170301029"),
		},
	}
}

func TestAuthorizationServer_Check(t *testing.T) {
	type args struct {
		ctx context.Context
		req *envoy_service_auth_v2.CheckRequest
	}

	request := &envoy_service_auth_v2.CheckRequest{
		Attributes: &envoy_service_auth_v2.AttributeContext{
			Source:      nil,
			Destination: nil,
			Request: &envoy_service_auth_v2.AttributeContext_Request{
				Time: nil,
				Http: &envoy_service_auth_v2.AttributeContext_HttpRequest{
					Id:       "",
					Method:   "GET",
					Headers:  nil,
					Path:     "/",
					Host:     "",
					Scheme:   "",
					Query:    "",
					Fragment: "",
					Size:     0,
					Protocol: "",
					Body:     "",
				},
			},
			ContextExtensions: nil,
			MetadataContext:   nil,
		},
	}
	args1 := args{
		ctx: nil,
		req: request,
	}

	want1 := &envoy_service_auth_v2.CheckResponse{
		Status: &status.Status{
			Code:    0,
			Message: "OK",
			Details: nil,
		},
		HttpResponse: &envoy_service_auth_v2.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v2.OkHttpResponse{
				Headers: []*envoy_api_v2_core.HeaderValueOption{
					{
						Header: &envoy_api_v2_core.HeaderValue{
							Key:   "x-custom-header-from-authz",
							Value: "authenticated via hybris",
						},
					},
				},
			},
		},
	}

	userFromHybris = func() (int, *User, error) {
		return 200, &User{"123"}, nil
	}

	exchange.JwtTokenExchange = func(id string) (*string, error) {
		token := "fakeToken"
		return &token, nil
	}

	tests := []struct {
		name    string
		args    args
		want    *envoy_service_auth_v2.CheckResponse
		wantErr bool
	}{
		{"success", args1, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationServer{}
			got, err := a.Check(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationServer_Check_Failure(t *testing.T) {
	type args struct {
		ctx context.Context
		req *envoy_service_auth_v2.CheckRequest
	}

	request := &envoy_service_auth_v2.CheckRequest{
		Attributes: &envoy_service_auth_v2.AttributeContext{
			Source:      nil,
			Destination: nil,
			Request: &envoy_service_auth_v2.AttributeContext_Request{
				Time: nil,
				Http: &envoy_service_auth_v2.AttributeContext_HttpRequest{
					Id:       "",
					Method:   "GET",
					Headers:  nil,
					Path:     "/",
					Host:     "",
					Scheme:   "",
					Query:    "",
					Fragment: "",
					Size:     0,
					Protocol: "",
					Body:     "",
				},
			},
			ContextExtensions: nil,
			MetadataContext:   nil,
		},
	}
	args1 := args{
		ctx: nil,
		req: request,
	}

	want1 := &envoy_service_auth_v2.CheckResponse{
		Status: &status.Status{
			Code:    16,
			Message: "unauthenticated in hybris",
			Details: nil,
		},
		HttpResponse: &envoy_service_auth_v2.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v2.DeniedHttpResponse{
				Status:  &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Headers: nil,
				Body:    "Hybris failed to authenticate",
			},
		},
	}

	userFromHybris = func() (int, *User, error) {
		return 400, &User{}, fmt.Errorf("error from hybris")
	}

	tests := []struct {
		name    string
		args    args
		want    *envoy_service_auth_v2.CheckResponse
		wantErr bool
	}{
		{"failure", args1, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationServer{}
			got, err := a.Check(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Check() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDeniedResponse(t *testing.T) {
	type args struct {
		statusCode rpc.Code
		statusMsg  string
		bodyMsg    string
	}

	args1 := args{
		statusCode: rpc.UNAUTHENTICATED,
		statusMsg: "",
		bodyMsg: "Hybris failed to authenticate",
	}

	want1 := &envoy_service_auth_v2.CheckResponse{
		Status: &status.Status{
			Code:    16,
			Message: "",
			Details: nil,
		},
		HttpResponse: &envoy_service_auth_v2.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v2.DeniedHttpResponse{
				Status:  &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Headers: nil,
				Body:    "Hybris failed to authenticate",
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *envoy_service_auth_v2.CheckResponse
		wantErr bool
	}{
		{"success", args1, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDeniedResponse(tt.args.statusCode, tt.args.statusMsg, tt.args.bodyMsg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildDeniedResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDeniedResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildOkResponse(t *testing.T) {
	type args struct {
		statusCode   rpc.Code
		statusMsg    string
		headerValues []headerValue
	}

	args1 := args{
		statusCode: rpc.UNAUTHENTICATED,
		statusMsg: "",
		headerValues: []headerValue{
			{
				key: "x-custom-header-from-authz",
				value: "authenticated via hybris",
			},
		},
	}

	want1 := &envoy_service_auth_v2.CheckResponse{
		Status: &status.Status{
			Code:    16,
			Message: "",
			Details: nil,
		},
		HttpResponse: &envoy_service_auth_v2.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v2.OkHttpResponse{
				Headers: []*envoy_api_v2_core.HeaderValueOption{
					{
						Header: &envoy_api_v2_core.HeaderValue{
							Key:   "x-custom-header-from-authz",
							Value: "authenticated via hybris",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *envoy_service_auth_v2.CheckResponse
		wantErr bool
	}{
		{"success", args1, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOkResponse(tt.args.statusCode, tt.args.statusMsg, tt.args.headerValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildOkResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildOkResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
//
//func Test_getSessionFromCookie(t *testing.T) {
//	type args struct {
//		req *envoy_service_auth_v2.CheckRequest
//	}
//
//	request := &envoy_service_auth_v2.CheckRequest{
//		Attributes: &envoy_service_auth_v2.AttributeContext{
//			Source:      nil,
//			Destination: nil,
//			Request: &envoy_service_auth_v2.AttributeContext_Request{
//				Time: nil,
//				Http: &envoy_service_auth_v2.AttributeContext_HttpRequest{
//					Id:       "",
//					Method:   "GET",
//					Headers:  nil,
//					Path:     "/",
//					Host:     "",
//					Scheme:   "",
//					Query:    "",
//					Fragment: "",
//					Size:     0,
//					Protocol: "",
//					Body:     "",
//				},
//			},
//			ContextExtensions: nil,
//			MetadataContext:   nil,
//		},
//	}
//
//	args1 := args{
//		req: request,
//	}
//
//	x := "test"
//	want1 := &x
//
//	tests := []struct {
//		name    string
//		args    args
//		want    *string
//		wantErr bool
//	}{
//		{"success", args1, want1, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := getSessionFromCookie(tt.args.req)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("getSessionFromCookie() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getSessionFromCookie() got = %v, want %v", *got, *tt.want)
//			}
//		})
//	}
//}