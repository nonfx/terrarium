package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func cacheFileName() string {
	return fmt.Sprintf("%s-prompt-%d.json", OPENAI_MODEL, OPENAI_PROMPT)
}

// generateEmbeddings generates vector embeddings for the given text strings using OpenAI API
func generateEmbeddings(tfTokens []TFToken) (EmbeddingsIndex, error) {
	embeddings, _ := readEmbeddingsFromFile(cacheFileName())
	if len(embeddings) > 0 {
		return embeddings, nil
	}

	log.Default().Println("Making API call")
	// Set up the OpenAI API request
	apiKey := os.Getenv("OPENAI_API_KEY")
	orgId := os.Getenv("OPENAI_ORG_ID")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY env var not set")
	}
	url := "https://api.openai.com/v1/embeddings"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}
	if orgId != "" {
		headers["OpenAI-Organization"] = orgId
	}

	prompts := make([]string, len(tfTokens))
	for i, t := range tfTokens {
		prompts[i] = t.Prompt()
	}

	data := map[string]interface{}{
		"model": OPENAI_MODEL,
		"input": prompts,
	}

	// Send the request to the OpenAI API
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Extract the embeddings from the API response
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Write the response to the file
	ioutil.WriteFile(getFilePath("api_resp_raw.json"), respData, 0644)

	var respJSON OpenAIEmbeddingResp
	if err := json.Unmarshal(respData, &respJSON); err != nil {
		return nil, err
	}

	if len(respJSON.Data) == 0 {
		return nil, fmt.Errorf("no embeddings found in response")
	}

	embeddings = EmbeddingsIndex{}
	for i, data := range respJSON.Data {
		embeddings.SetByTFToken(tfTokens[i], data.Embedding)
	}

	writeJSONToFile(embeddings, cacheFileName())

	return embeddings, nil
}
