package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
)

type UserServiceConfig struct {
	PostgresURL string
	Port        string
}

func (s UserServiceConfig) StartContainer(ctx context.Context, networkName string) (internalURL, mappedURL string) {
	dir, _ := os.Getwd()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{Context: filepath.Join(dir, "userservice")},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"user-service"},
			},
			Env:          s.env(),
			ExposedPorts: []string{s.Port},
			WaitingFor:   wait.ForListeningPort(nat.Port(s.Port)),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, nat.Port(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("http://%s:%s", "user-service", s.Port), fmt.Sprintf("http://localhost:%s", mappedPort.Port())
}

func (s UserServiceConfig) env() map[string]string {
	return map[string]string{"POSTGRES_URL": s.PostgresURL, "PORT": s.Port}
}
