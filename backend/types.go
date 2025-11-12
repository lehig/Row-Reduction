package main

// Matrix represents a 3x3 matrix
type Matrix struct {
	Data [3][3]float64 `json:"data"`
}

// RREFRequest represents the request body
type RREFRequest struct {
	Matrix Matrix `json:"matrix"`
}

// RREFResponse represents the response body
type RREFResponse struct {
	Original Matrix `json:"original"`
	RREF     Matrix `json:"rref"`
}

// RREF calculates the Reduced Row Echelon Form of a 3x3 matrix
func RREF(matrix [3][3]float64) [3][3]float64 {
	// Create a copy to avoid modifying the original
	result := matrix

	// Number of rows and columns
	rows := 3
	cols := 3

	// Lead (pivot) position
	lead := 0

	// Process each row
	for r := 0; r < rows; r++ {
		if lead >= cols {
			break
		}

		// Find the pivot row
		i := r
		for i < rows && result[i][lead] == 0 {
			i++
		}

		if i == rows {
			// No pivot in this column, move to next column
			lead++
			r-- // Stay on the same row
			continue
		}

		// Swap rows if necessary
		if i != r {
			for j := 0; j < cols; j++ {
				result[r][j], result[i][j] = result[i][j], result[r][j]
			}
		}

		// Normalize the pivot row (make leading entry 1)
		pivot := result[r][lead]
		if pivot != 0 {
			for j := 0; j < cols; j++ {
				result[r][j] /= pivot
			}
		}

		// Eliminate entries above and below the pivot
		for i = 0; i < rows; i++ {
			if i != r {
				factor := result[i][lead]
				for j := 0; j < cols; j++ {
					result[i][j] -= factor * result[r][j]
				}
			}
		}

		lead++
	}

	return result
}

