package main_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/mzzz-zzm/go-tdd-practice/adapters/httpserver"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"

	"github.com/mzzz-zzm/go-tdd-practice/specifications"
)

// go test -v -timeout 300s -run ^TestGreeterServer$ github.com/mzzz-zzm/go-tdd-practice/cmd/httpserver
func TestGreeterServer(t *testing.T) {
	compose, err := tc.NewDockerCompose("../../.devcontainer/docker-compose.yml")
	assert.NoError(t, err)
	t.Cleanup(func() {
		err := compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
		assert.NoError(t, err)
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	err = compose.Up(ctx, tc.Wait(true))
	assert.NoError(t, err)

	// ctx := context.Background()

	// req := testcontainers.ContainerRequest{
	// 	FromDockerfile: testcontainers.FromDockerfile{
	// 		Context:    "../../.",
	// 		Dockerfile: "./Dockerfile",
	// 		// set to false if you want less spam, but this is helpful if you're having troubles
	// 		PrintBuildLog: true,
	// 	},
	// 	ExposedPorts: []string{"8080:8080"},
	// 	WaitingFor:   wait.ForHTTP("/").WithPort("8080"),
	// }
	// container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	// 	ContainerRequest: req,
	// 	Started:          true,
	// })
	// assert.NoError(t, err)
	// t.Cleanup(func() {
	// 	assert.NoError(t, container.Terminate(ctx))
	// })

	// ip, err := container.Host(ctx)
	// assert.NoError(t, err)
	// mappedPort, err := container.MappedPort(ctx, "8080")
	// assert.NoError(t, err)
	// containerUrl := fmt.Sprintf("http://%s:%s", ip, mappedPort)

	// containerUrl := "http://testsvr:8080"
	testsvrContainer, err := compose.ServiceContainer(ctx, "testsvr")
	assert.NoError(t, err)

	port, err := testsvrContainer.MappedPort(ctx, "8080/tcp")
	assert.NoError(t, err)
	host, err := testsvrContainer.Host(ctx)
	assert.NoError(t, err)
	url := fmt.Sprintf("http://%s:%s", host, port.Port())

	// endPt, err := testsvrContainer.Endpoint(ctx, "http")
	// assert.NoError(t, err)

	driver := go_specs_greet.Driver{BaseURL: url}
	specifications.GreetSpecifications(t, driver)
}
