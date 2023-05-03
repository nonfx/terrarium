package main

import (
	"fmt"
	"strings"
)

func main() {
	// Define the text strings to compare
	tfTokens := []TFToken{
		{"aws_vpc", "security_group_id"},
		{"aws_security_group", "id"},
		{"aws_security_group", "vpc_id"},
		{"aws_vpc", "id"},
	}

	// for i := 1; i <= 7; i++ {
	// 	OPENAI_PROMPT = i

	// Generate vector embeddings for each text string
	embeddings, err := generateEmbeddings(tfTokens)
	if err != nil {
		panic(err)
	}

	// Calculate cosine similarity between each pair of text strings
	similarityMatrix := make([][]float64, len(tfTokens))
	for i := range tfTokens {
		similarityMatrix[i] = make([]float64, len(tfTokens))
		for j := range tfTokens {
			similarityMatrix[i][j] = cosineSimilarity(
				embeddings.GetByTFToken(tfTokens[i]),
				embeddings.GetByTFToken(tfTokens[j]),
			)
		}
	}

	// Print the similarity matrix
	fmt.Printf("\nPrompt eg: %s\n", tfTokens[0].Prompt())
	fmt.Println("Cosine similarity matrix:")
	printTable(tfTokens, similarityMatrix)
	// }
}

// printTable prints a table given a header row and a 2D array of data
func printTable(header []TFToken, data [][]float64) {
	// Determine the maximum length of each column
	colWidths := make([]int, len(header))
	maxWidth := 0
	for i, col := range header {
		colWidths[i] = len(col.String())
		if maxWidth < colWidths[i] {
			maxWidth = colWidths[i]
		}
	}

	// Print the header row
	headerLine := fmt.Sprintf("| %-*s ", maxWidth, "")
	for i, col := range header {
		headerLine += fmt.Sprintf("| %-*s ", colWidths[i], col)
	}
	headerLine += "|"
	fmt.Println(headerLine)
	dividerLine := strings.Repeat("-", len(headerLine))
	fmt.Println(dividerLine)

	// Print the data rows
	for i, row := range data {
		dataLine := fmt.Sprintf("| %-*s ", maxWidth, header[i])
		for j, col := range row {
			dataLine += fmt.Sprintf("| %-*f ", colWidths[j], col)
		}
		dataLine += "|"
		fmt.Println(dataLine)
	}
}
