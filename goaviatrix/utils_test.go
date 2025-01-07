package goaviatrix

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateASN(t *testing.T) {
	tt := []struct {
		Name        string
		Input       interface{}
		ExpectedErr string
	}{
		{
			"too small",
			"0",
			`"test" must be an integer in 1-4294967294, got: 0`,
		},
		{
			"too large",
			"4294967295",
			`"test" must be an integer in 1-4294967294, got: 4294967295`,
		},
		{
			"wrong type",
			65001,
			`"test" must be of type string`,
		},
		{
			"passing",
			"4294967294",
			"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, errs := ValidateASN(tc.Input, "test")
			if tc.ExpectedErr != "" {
				if len(errs) < 1 {
					t.Fatalf("test case %q expected an error: %q, got: none", tc.Name, tc.ExpectedErr)
				}
				if errs[0].Error() != tc.ExpectedErr {
					t.Fatalf("test case %q expected an error: %q, got: %q", tc.Name, tc.ExpectedErr, errs[0].Error())
				}
			} else {
				if len(errs) > 0 {
					t.Fatalf("test case %q expected no error, got %q", tc.Name, errs[0].Error())
				}
			}
		})
	}
}

func TestMapContains(t *testing.T) {
	testMap := make(map[string]interface{})
	testKeys := []string{"one", "two", "three"}
	for _, key := range testKeys {
		testMap[key] = key
	}
	assert.True(t, MapContains(testMap, "one"))
	assert.False(t, MapContains(testMap, "Random"))
}

func TestCreateZtpFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		content     string
		expectedErr string
	}{
		{
			name:        "Successful file creation and writing",
			filePath:    "test-file.txt",
			content:     "This is a test file.",
			expectedErr: "",
		},
		{
			name:        "Error creating file (invalid path)",
			filePath:    "/invalid/path/test-file.txt",
			content:     "This content should not be written.",
			expectedErr: "failed to create the file",
		},
		{
			name:        "Error writing to file (empty content)",
			filePath:    "test-file.txt",
			content:     "",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the createZtpFile function
			err := createZtpFile(tt.filePath, tt.content)

			if tt.expectedErr != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)

				// Check if the file is created and the content is written (if no error occurred)
				if _, err := os.Stat(tt.filePath); err == nil {
					// File exists, check the content
					data, err := os.ReadFile(tt.filePath)
					assert.NoError(t, err)
					assert.Equal(t, tt.content, string(data))
					// Clean up the file after test
					os.Remove(tt.filePath)
				}
			}
		})
	}
}
