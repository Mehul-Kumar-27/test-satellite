package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	// Create a file store in the current working directory
	fs, err := file.New(cwd)
	if err != nil {
		fmt.Printf("Error creating file store: %v\n", err)
		return
	}
	defer fs.Close()

	ctx := context.Background()

	// Add state.json to the file store
	fileName := filepath.Join("sample_app", "state_now.json")
	filePath := filepath.Join(cwd, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", filePath)
		return
	}
	fmt.Printf("File exists: %s\n", filePath)

	mediaType := "application/json"
	fileDescriptor, err := fs.Add(ctx, fileName, mediaType, "")
	if err != nil {
		fmt.Printf("Error adding file to store: %v\n", err)
		return
	}
	fmt.Printf("File descriptor for %s: %v\n", fileName, fileDescriptor)

	// Pack the file and tag the packed manifest
	artifactType := "application/json"
	opts := oras.PackManifestOptions{
		Layers: []v1.Descriptor{fileDescriptor},
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		fmt.Printf("Error packing manifest: %v\n", err)
		return
	}
	fmt.Println("Manifest descriptor:", manifestDescriptor)

	tag := "latest"
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		fmt.Printf("Error tagging manifest: %v\n", err)
		return
	}

	// Connect to Harbor repository
	harborRegistry := "demo.goharbor.io"
	repo, err := remote.NewRepository(harborRegistry + "/test-satellite-group/state-artifact")
	if err != nil {
		fmt.Printf("Error creating repository: %v\n", err)
		return
	}

	// Set up authentication for Harbor
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(harborRegistry, auth.Credential{
			Username: "admin",
			Password: "Harbor12345",
		}),
	}

	// Copy from the file store to Harbor
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		fmt.Printf("Error copying to repository: %v\n", err)
		return
	}

	fmt.Println("Successfully pushed state.json to Harbor")
}
