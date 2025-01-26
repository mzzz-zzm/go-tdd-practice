package main_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/mzzz-zzm/go-tdd-practice/adapters/httpserver"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"

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
	// driver := go_specs_greet.Driver{BaseURL: url, Client: &client}

	// --- case 2: use Endpoint ---
	endPt, err := testsvrContainer.Endpoint(ctx, "http")
	assert.NoError(t, err)

	driver := go_specs_greet.Driver{
		BaseURL: endPt,
		Client:  &client,
	}

	specifications.GreetSpecifications(t, driver)
}
