package health

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"reflect"
	"testing"
)

func Test_healthServer_Check(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *grpc_health_v1.HealthCheckRequest
	}

	args1 := args{
		ctx: nil,
		in: &grpc_health_v1.HealthCheckRequest{
			Service: "",
		},
	}

	want1 := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}

	tests := []struct {
		name    string
		args    args
		want    *grpc_health_v1.HealthCheckResponse
		wantErr bool
	}{
		{"success", args1, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthServer{}
			got, err := s.Check(tt.args.ctx, tt.args.in)
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

func Test_healthServer_Watch(t *testing.T) {
	type args struct {
		in  *grpc_health_v1.HealthCheckRequest
		srv grpc_health_v1.Health_WatchServer
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//todo: NOT USED
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthServer{}
			if err := s.Watch(tt.args.in, tt.args.srv); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
