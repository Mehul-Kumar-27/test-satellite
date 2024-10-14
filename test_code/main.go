package main

import "fmt"

type Artifact struct {
    Name string
}

func main() {
    artifacts := []Artifact{
        {"Artifact 1"},
        {"Artifact 2"},
        {"Artifact 3"},
    }

    var incorrectArtifactPointers []*Artifact
    var artifact Artifact // declare outside the loop
    for _, a := range artifacts {
        artifact = a
        incorrectArtifactPointers = append(incorrectArtifactPointers, &artifact)
    }

    var correctArtifactPointers []*Artifact
    for i := range artifacts {
        correctArtifactPointers = append(correctArtifactPointers, &artifacts[i])
    }

    fmt.Println("Incorrect Pointers:")
    for _, pointer := range incorrectArtifactPointers {
        fmt.Printf("%p: %s\n", pointer, pointer.Name)
    }

    fmt.Println("\nCorrect Pointers:")
    for _, pointer := range correctArtifactPointers {
        fmt.Printf("%p: %s\n", pointer, pointer.Name)
    }
}
