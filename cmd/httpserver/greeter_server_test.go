package main_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"

	"github.com/mzzz-zzm/go-tdd-practice/adapters"
	"github.com/mzzz-zzm/go-tdd-practice/adapters/httpserver"
	"github.com/mzzz-zzm/go-tdd-practice/specifications"
)

// first time go mod download takes a long time, avoid timeout by increasing -timeout
// ex) go test -v -timeout 300s -run ^TestGreeterServer$ github.com/mzzz-zzm/go-tdd-practice/cmd/httpserver
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

	testsvrContainer, err := compose.ServiceContainer(ctx, "testsvr")
	assert.NoError(t, err)

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	// --- case 1: use MappedPort and Host ---
	// port, err := testsvrContainer.MappedPort(ctx, "8080/tcp")
	// assert.NoError(t, err)
	// host, err := testsvrContainer.Host(ctx)
	// assert.NoError(t, err)
	// url := fmt.Sprintf("http://%s:%s", host, port.Port())
	// driver := httpserver.Driver{BaseURL: url, Client: &client}

	// --- case 2: use Endpoint ---
	endPt, err := testsvrContainer.Endpoint(ctx, "http")
	assert.NoError(t, err)

	driver := httpserver.Driver{
		BaseURL: endPt,
		Client:  &client,
	}

	specifications.GreetSpecifications(t, driver)
}

// first time go mod download takes a long time, avoid timeout by increasing -timeout
// ex) go test -v -timeout 300s -run ^TestGreeterServerWithTemplateConfig$ github.com/mzzz-zzm/go-tdd-practice/cmd/httpserver
func TestGreeterServerWithTemplateConfig(t *testing.T) {
	// Get the absolute path to the project root (where the Dockerfile is)
	// projectRoot, err := filepath.Abs("../..") // Adjust this path if needed
	// assert.NoError(t, err)
	// dcfpath := filepath.Join(projectRoot, "Dockerfile")

	dockerConfig := adapters.DockerConfig{
		DockerFileName: "Dockerfile",
		ServiceName:    "testsvr",
		ContainerName:  "testsvr",
		Port:           8080,
		Protocol:       "http",
		BinToBuild:     "httpserver",
	}

	endPt := adapters.StartDockerServer(t, dockerConfig)

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := httpserver.Driver{
		BaseURL: endPt,
		Client:  &client,
	}

	specifications.GreetSpecifications(t, driver)
	specifications.CurseSpecifications(t, driver)
}
