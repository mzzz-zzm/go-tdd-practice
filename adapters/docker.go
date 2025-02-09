package adapters

import (
	"bytes"
	"context"
	"os"
	"testing"
	"text/template"

	"github.com/alecthomas/assert/v2"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

const DOCKER_COMPOSE_TEMPLATE = `version: '3.9'

services:
  {{.ServiceName}}:
    container_name: {{.ContainerName}}
    build:
      context: ../../
      dockerfile: {{.DockerFileName}}
      args:
        - BIN_TO_BUILD={{.BinToBuild}}
        - PORT_TO_EXPOSE={{.Port}}
    ports:
      - "{{.Port}}:{{.Port}}"
    
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    
    extra_hosts:
      - 'host.docker.internal:host-gateway'
    
    environment:
      - TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal
`

type DockerConfig struct {
	DockerFileName string
	ServiceName    string
	ContainerName  string
	Port           uint
	Protocol       string
	BinToBuild     string
}

func StartDockerServer(
	t testing.TB,
	dockerConfig DockerConfig,
) (endPoint string) {
	t.Helper()
	templ, err := template.New("docker-compose.yml").Parse(DOCKER_COMPOSE_TEMPLATE)
	assert.NoError(t, err)
	buf := bytes.Buffer{}
	err = templ.Execute(&buf, dockerConfig)
	assert.NoError(t, err)

	// create temporary docker-compose.yml file
	tmpFile, err := os.CreateTemp(".", "docker-compose-*.yml")
	assert.NoError(t, err)

	defer func() {
		err := os.Remove(tmpFile.Name())
		assert.NoError(t, err)
	}()

	defer tmpFile.Close()
	_, err = tmpFile.Write(buf.Bytes())
	assert.NoError(t, err)

	compose, err := tc.NewDockerCompose(tmpFile.Name())
	assert.NoError(t, err)
	t.Cleanup(func() {
		err := compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
		assert.NoError(t, err)
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// TODO test start
	// err = os.WriteFile("generated-docker-compose.yml", buf.Bytes(), 0644)
	// assert.NoError(t, err)
	// TODO test ends

	err = compose.Up(ctx, tc.Wait(true))
	assert.NoError(t, err)

	container, err := compose.ServiceContainer(ctx, dockerConfig.ContainerName)
	assert.NoError(t, err)
	endPoint, err = container.Endpoint(ctx, dockerConfig.Protocol)
	assert.NoError(t, err)

	return
}
