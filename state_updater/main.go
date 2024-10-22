package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
)

type Artifact struct {
	Repository string   `json:"repository"`
	Tags       []string `json:"tag"`
	Labels     []string `json:"labels"`
	Type       string   `json:"type"`
	Digest     string   `json:"digest"`
	Deleted    bool     `json:"deleted"`
}

type Result struct {
	Registry  string     `json:"registry"`
	Artifacts []Artifact `json:"artifacts"`
}

func main() {
	registry := "registry.bupd.xyz"
	groupName := "satellite-test-group-state"
	stateName := "state"


	var result Result
	// Path to the json file to upload
	jsonPath := "updated_state.json"
	// read this json file
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("Failed to open json file: %v", err)
	}
	defer jsonFile.Close()
	// Read the json file and marshal it ito the result
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read json file: %v", err)
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal json data: %v", err)
	}

	data, err := json.Marshal(result)

	if err != nil {
		log.Fatalf("Failed to marshal json data: %v", err)
	}

	img, err := crane.Image(map[string][]byte{
		"artifacts.json": data,
	})

	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	destination := fmt.Sprintf("%s/%s/%s", registry, groupName, stateName)

	err = crane.Push(img, destination)
	if err != nil {
		log.Fatalf("Failed to push image: %v", err)
	}

	err = crane.Tag(destination, "latest")
	if err != nil {
		log.Fatalf("Failed to tag image: %v", err)
	}
}



