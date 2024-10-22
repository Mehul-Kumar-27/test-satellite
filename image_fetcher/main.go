package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
)

type Result struct {
	Registry  string     `json:"registry"`
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Repository string   `json:"repository"`
	Tag        []string `json:"tag"`
	Labels     []string `json:"labels"`
	Type       string   `json:"type"`
	Digest     string   `json:"digest"`
	Deleted    bool     `json:"deleted"`
}

func main() {
	imgRef := "registry.bupd.xyz/satellite-test-group-state/state:latest"

	// Pull the image from the registry
	img, err := crane.Pull(imgRef)
	if err != nil {
		log.Fatalf("Failed to pull image: %v", err)
	}

	// Export the content of the image (as a tar file)
	tarContent := new(bytes.Buffer)
	if err := crane.Export(img, tarContent); err != nil {
		log.Fatalf("Failed to export image: %v", err)
	}

	// Parse the tar content to extract `artifacts.json`
	tr := tar.NewReader(tarContent)
	var artifactsJSON []byte
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			log.Fatalf("Error reading tar archive: %v", err)
		}

		if hdr.Name == "artifacts.json" {
			// Found `artifacts.json`, read the content
			artifactsJSON, err = io.ReadAll(tr)
			if err != nil {
				log.Fatalf("Failed to read artifacts.json: %v", err)
			}
			break
		}
	}

	if artifactsJSON == nil {
		log.Fatal("artifacts.json not found in image")
	}

	// Unmarshal the JSON content into the Result struct
	var result Result
	err = json.Unmarshal(artifactsJSON, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal artifacts.json: %v", err)
	}

	file, err := os.Create("state_now.json")
	jsn, err := json.MarshalIndent(result, "", "  ")
	file.Write(jsn)

	// Print the extracted data
	fmt.Printf("Registry: %s\n", result.Registry)
	for _, artifact := range result.Artifacts {
		fmt.Printf("Artifact: %s, Digest: %s\n", artifact.Repository, artifact.Digest)
	}
}
