package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

func getFilePath(filename string) string {
	return path.Join(".", "cache_data", filename)
}

// readEmbeddingsFromFile reads vector embeddings from a JSON file
func readEmbeddingsFromFile(filename string) (EmbeddingsIndex, error) {
	// Read the contents of the JSON file
	data, err := ioutil.ReadFile(getFilePath(filename))
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a slice of slices of float64
	var embeddings EmbeddingsIndex
	if err := json.Unmarshal(data, &embeddings); err != nil {
		return nil, err
	}

	return embeddings, nil
}

// writeJSONToFile writes data as a JSON string to a file
func writeJSONToFile(data interface{}, filename string) error {
	// Marshal the data into a JSON string
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON string to the file
	if err := ioutil.WriteFile(getFilePath(filename), jsonData, 0644); err != nil {
		return err
	}

	return nil
}
