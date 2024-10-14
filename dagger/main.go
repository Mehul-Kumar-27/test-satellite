// A generated module for TestSatellite functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"

	"github.com/Mehul-Kumar-27/test-satellite/dagger/internal/dagger"
)

const (
	DEFAULT_GO = "golang:1.22"
	MOUNT      = "/app"
	ALPINE     = "alpine:latest"
)

type TestSatellite struct{}

func (t *TestSatellite) Attach(
	ctx context.Context,
	container *dagger.Container,
	// +optional
	// +default="24.0"
	dockerVersion string,
) (*dagger.Container, error) {
	dockerd := t.Service(dockerVersion)

	dockerHost, err := dockerd.Endpoint(ctx, dagger.ServiceEndpointOpts{
		Scheme: "tcp",
	})
	if err != nil {
		return nil, err
	}
	dockerd, err = dockerd.Start(ctx)
	if err != nil {
		return nil, err
	}

	return container.
		WithServiceBinding("docker", dockerd).
		WithEnvVariable("DOCKER_HOST", dockerHost), nil
}

// Get a Service container running dockerd
func (t *TestSatellite) Service(
	// +optional
	// +default="24.0"
	dockerVersion string,
) *dagger.Service {
	port := 2375
	return dag.Container().
		From(fmt.Sprintf("docker:%s-dind", dockerVersion)).
		WithMountedCache(
			"/var/lib/docker",
			dag.CacheVolume(dockerVersion+"-docker-lib"),
			dagger.ContainerWithMountedCacheOpts{
				Sharing: dagger.Private,
			}).
		WithExposedPort(port).
		WithExec([]string{
			"dockerd",
			"--host=tcp://0.0.0.0:2375",
			"--host=unix:///var/run/docker.sock",
			"--tls=false",
		}, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		AsService()
}
func (m *TestSatellite) Publish(ctx context.Context, source *dagger.Directory, name string) (string, error) {
	container := dag.Container().
		From("docker:24.0")

	container, err := m.Attach(ctx, container, "24.0")
	if err != nil {
		return "", err
	}

	image_to_pull := fmt.Sprintf("%s:%s", name, "latest")

	output, err := container.
		WithMountedDirectory(MOUNT, source).
		WithWorkdir(MOUNT).
		WithExec([]string{"docker", "login", "localhost:3080", "-u", "admin", "-p", "Harbor12345"}).
		WithExec([]string{"docker", "pull", image_to_pull}).
		WithExec([]string{"docker", "tag", "alpine:latest", fmt.Sprintf("demo.goharbor.io/satellite-test-%s/%s", name, image_to_pull)}).
		WithExec([]string{"docker", "push", fmt.Sprintf("demo.goharbor.io/satellite-test-%s/%s", name, image_to_pull)}).
		Stderr(ctx)

	if err != nil {
		return output, err
	}
	return output, nil
}
