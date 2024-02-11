package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessRows(t *testing.T) {
	tests := []struct {
		name         string
		inputRows    []CSVRow
		wantOutput   UploadInventoryInput
		wantErrs     []Error
		wantErrCount int // Count of errors, helpful if expecting multiple errors
	}{
		// Test Case definitions here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotErrs := processRows(tt.inputRows)

			// Assert on the returned output
			assert.Equal(t, tt.wantOutput, gotOutput)

			// Assert on the returned errors
			assert.Equal(t, tt.wantErrCount, len(gotErrs)) // Check the number of errors

			// Optionally, assert specific error messages
			if len(gotErrs) > 0 {
				for i, err := range gotErrs {
					assert.Equal(t, tt.wantErrs[i].Err.Error(), err.Err.Error(), "Error messages should match")
				}
			}
		})
	}
}
