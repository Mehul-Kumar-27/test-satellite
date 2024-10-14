package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func Fetch() {
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

	// Connect to Harbor repository
	ctx := context.Background()
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

	// Copy from the remote repository to the file store
	tag := "latest"
	manifestDescriptor, err := oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		fmt.Printf("Error copying from repository: %v\n", err)
		return
	}

	fmt.Println("Manifest descriptor:", manifestDescriptor)

	// Find and print the content of the downloaded state.json
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "state.json" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Println("Contents of downloaded state.json:")
			fmt.Println(string(content))
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		return
	}

	fmt.Println("Successfully fetched state.json from Harbor")
}
