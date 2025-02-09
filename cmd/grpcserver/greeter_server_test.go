package main_test

import (
	"testing"

	"github.com/mzzz-zzm/go-tdd-practice/adapters"
	"github.com/mzzz-zzm/go-tdd-practice/adapters/grpcserver"
	"github.com/mzzz-zzm/go-tdd-practice/specifications"
)

// go test -v -timeout 300s -run ^TestGreeterServerWithTemplateConfig$ github.com/mzzz-zzm/go-tdd-practice/cmd/grpcserver
func TestGreeterServerWithTemplateConfig(t *testing.T) {
	dockerConfig := adapters.DockerConfig{
		DockerFileName: "Dockerfile",
		ServiceName:    "testsvr_grpc",
		ContainerName:  "testsvr_grpc",
		Port:           50051,
		Protocol:       "",
		BinToBuild:     "grpcserver",
	}

	endPt := adapters.StartDockerServer(t, dockerConfig)
	driver := grpcserver.Driver{
		Addr: endPt,
	}
	defer driver.Close()

	specifications.GreetSpecifications(t, &driver)
	specifications.CurseSpecifications(t, &driver)
}
