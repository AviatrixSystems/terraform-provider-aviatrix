package goaviatrix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVersionSupported(t *testing.T) {
	tests := []struct {
		name              string
		currentVersion    string
		supportedVersions []string
		expectedError     error
	}{
		{
			name:              "Supported version matches",
			currentVersion:    "8.0.1",
			supportedVersions: []string{"8.0"},
			expectedError:     nil,
		},
		{
			name:              "Supported version matches with pre-release",
			currentVersion:    "8.0.1-1000.1000",
			supportedVersions: []string{"8.0.2-2000.2000"},
			expectedError:     nil,
		},
		{
			name:              "Unsupported version",
			currentVersion:    "8.0.1",
			supportedVersions: []string{"8.1"},
			expectedError:     fmt.Errorf("version %s is not supported", "8.0.1"),
		},
		{
			name:              "Multiple supported versions with match",
			currentVersion:    "8.0.1",
			supportedVersions: []string{"7.5", "8.0"},
			expectedError:     nil,
		},
		{
			name:              "Multiple supported versions with no match",
			currentVersion:    "8.0.1",
			supportedVersions: []string{"7.5", "8.1"},
			expectedError:     fmt.Errorf("version %s is not supported", "8.0.1"),
		},
		{
			name:              "Old UserConnect versions",
			currentVersion:    "UserConnect-7.2-1804.4665",
			supportedVersions: []string{"8.0.0"},
			expectedError:     fmt.Errorf("version %s is not supported", "UserConnect-7.2-1804.4665"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isVersionSupported(tt.currentVersion, tt.supportedVersions)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
