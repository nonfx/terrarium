package main

import "math"

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) == 0 || len(vec2) == 0 {
		return 0
	}

	dotProduct := 0.0
	for i, val := range vec1 {
		dotProduct += val * vec2[i]
	}
	magnitude1 := math.Sqrt(sumSquares(vec1))
	magnitude2 := math.Sqrt(sumSquares(vec2))
	return dotProduct / (magnitude1 * magnitude2)
}

// sumSquares calculates the sum of squares of a vector
func sumSquares(vec []float64) float64 {
	sum := 0.0
	for _, val := range vec {
		sum += val * val
	}
	return sum
}
