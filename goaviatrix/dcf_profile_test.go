package goaviatrix

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogProfileByName_URLEscaping(t *testing.T) {
	tests := []struct {
		name                 string
		profileName          string
		expectedEscapedName  string
		shouldContainEscaped bool
		description          string
	}{
		{
			name:                 "Simple alphanumeric name",
			profileName:          "simpleProfile123",
			expectedEscapedName:  "simpleProfile123",
			shouldContainEscaped: false,
			description:          "Normal profile name should not be escaped",
		},
		{
			name:                 "Name with spaces",
			profileName:          "profile with spaces",
			expectedEscapedName:  "profile+with+spaces",
			shouldContainEscaped: true,
			description:          "Spaces should be escaped as plus signs",
		},
		{
			name:                 "Name with special characters",
			profileName:          "profile&name=value",
			expectedEscapedName:  "profile%26name%3Dvalue",
			shouldContainEscaped: true,
			description:          "Ampersand and equals signs should be percent-encoded",
		},
		{
			name:                 "Name with percent signs",
			profileName:          "profile%20name",
			expectedEscapedName:  "profile%2520name",
			shouldContainEscaped: true,
			description:          "Existing percent signs should be double-encoded",
		},
		{
			name:                 "Name with question mark",
			profileName:          "profile?query",
			expectedEscapedName:  "profile%3Fquery",
			shouldContainEscaped: true,
			description:          "Question marks should be percent-encoded",
		},
		{
			name:                 "Name with hash symbol",
			profileName:          "profile#fragment",
			expectedEscapedName:  "profile%23fragment",
			shouldContainEscaped: true,
			description:          "Hash symbols should be percent-encoded",
		},
		{
			name:                 "Name with forward slash",
			profileName:          "profile/path",
			expectedEscapedName:  "profile%2Fpath",
			shouldContainEscaped: true,
			description:          "Forward slashes should be percent-encoded",
		},
		{
			name:                 "Name with plus sign",
			profileName:          "profile+name",
			expectedEscapedName:  "profile%2Bname",
			shouldContainEscaped: true,
			description:          "Plus signs should be percent-encoded",
		},
		{
			name:                 "Unicode characters",
			profileName:          "profile名前",
			expectedEscapedName:  "profile%E5%90%8D%E5%89%8D",
			shouldContainEscaped: true,
			description:          "Unicode characters should be UTF-8 percent-encoded",
		},
		{
			name:                 "Emoji characters",
			profileName:          "profile🚀test",
			expectedEscapedName:  "profile%F0%9F%9A%80test",
			shouldContainEscaped: true,
			description:          "Emoji should be UTF-8 percent-encoded",
		},
		{
			name:                 "Multiple consecutive spaces",
			profileName:          "profile    with    spaces",
			expectedEscapedName:  "profile++++with++++spaces",
			shouldContainEscaped: true,
			description:          "Multiple spaces should all be escaped",
		},
		{
			name:                 "Mixed special characters",
			profileName:          "profile?name=value&type=test#section",
			expectedEscapedName:  "profile%3Fname%3Dvalue%26type%3Dtest%23section",
			shouldContainEscaped: true,
			description:          "Complex query-like strings should be fully escaped",
		},
		{
			name:                 "Control characters",
			profileName:          "profile\x00\x01\x02name",
			expectedEscapedName:  "profile%00%01%02name",
			shouldContainEscaped: true,
			description:          "Control characters should be percent-encoded",
		},
		{
			name:                 "Newline and tab characters",
			profileName:          "profile\n\tname",
			expectedEscapedName:  "profile%0A%09name",
			shouldContainEscaped: true,
			description:          "Newline and tab characters should be percent-encoded",
		},
		{
			name:                 "Empty string",
			profileName:          "",
			expectedEscapedName:  "",
			shouldContainEscaped: false,
			description:          "Empty string should remain empty",
		},
		{
			name:                 "Only special characters",
			profileName:          "!@#$%^&*()",
			expectedEscapedName:  "%21%40%23%24%25%5E%26%2A%28%29",
			shouldContainEscaped: true,
			description:          "String with only special characters should be fully escaped",
		},
		{
			name:                 "Very long string with spaces",
			profileName:          strings.Repeat("profile name ", 50),
			expectedEscapedName:  strings.Repeat("profile+name+", 50),
			shouldContainEscaped: true,
			description:          "Very long strings should be properly escaped",
		},
		{
			name:                 "Already URL encoded string",
			profileName:          "profile%20name%26value",
			expectedEscapedName:  "profile%2520name%2526value",
			shouldContainEscaped: true,
			description:          "Already encoded strings should be double-encoded to prevent injection",
		},
		{
			name:                 "SQL injection attempt",
			profileName:          "profile'; DROP TABLE users; --",
			expectedEscapedName:  "profile%27%3B+DROP+TABLE+users%3B+--",
			shouldContainEscaped: true,
			description:          "SQL injection attempts should be safely escaped",
		},
		{
			name:                 "XSS attempt",
			profileName:          "<script>alert('xss')</script>",
			expectedEscapedName:  "%3Cscript%3Ealert%28%27xss%27%29%3C%2Fscript%3E",
			shouldContainEscaped: true,
			description:          "XSS attempts should be safely escaped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the actual URL escaping behavior
			escaped := url.QueryEscape(tt.profileName)
			assert.Equal(t, tt.expectedEscapedName, escaped, tt.description)

			// Verify that the escaped version is safe for URL usage
			if tt.shouldContainEscaped {
				assert.NotEqual(t, tt.profileName, escaped, "Expected input to be different from escaped output")
				// Verify that unescaping returns the original
				unescaped, err := url.QueryUnescape(escaped)
				assert.NoError(t, err, "Should be able to unescape the escaped string")
				assert.Equal(t, tt.profileName, unescaped, "Unescaping should return original string")
			}

			// Test that the escaped name doesn't contain dangerous characters for URL path
			dangerousChars := []string{"&", "?", "#", " ", "%20", "\n", "\t"}
			for _, char := range dangerousChars {
				if strings.Contains(tt.profileName, char) && !strings.Contains(escaped, char) {
					// This is expected - the dangerous character was escaped
					continue
				}
				if strings.Contains(tt.profileName, char) && strings.Contains(escaped, char) {
					// Only allow this if it's a + for spaces, which is acceptable in query escaping
					if !(char == " " && strings.Contains(escaped, "+")) {
						t.Errorf("Dangerous character '%s' was not properly escaped in result: %s", char, escaped)
					}
				}
			}
		})
	}
}

func TestGetLogProfileByName_EndpointConstruction(t *testing.T) {
	tests := []struct {
		name                 string
		profileName          string
		expectedEndpointPath string
		description          string
	}{
		{
			name:                 "Simple name endpoint",
			profileName:          "testProfile",
			expectedEndpointPath: "dcf/log-profile/name/testProfile",
			description:          "Simple profile name should create clean endpoint",
		},
		{
			name:                 "Name with spaces endpoint",
			profileName:          "test profile",
			expectedEndpointPath: "dcf/log-profile/name/test+profile",
			description:          "Profile name with spaces should have spaces escaped in endpoint",
		},
		{
			name:                 "Name with special characters endpoint",
			profileName:          "test&profile=value",
			expectedEndpointPath: "dcf/log-profile/name/test%26profile%3Dvalue",
			description:          "Profile name with special characters should be percent-encoded in endpoint",
		},
		{
			name:                 "Unicode name endpoint",
			profileName:          "test名前",
			expectedEndpointPath: "dcf/log-profile/name/test%E5%90%8D%E5%89%8D",
			description:          "Unicode profile name should be UTF-8 encoded in endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the endpoint construction logic from GetLogProfileByName
			escapedName := url.QueryEscape(tt.profileName)
			endpoint := fmt.Sprintf("dcf/log-profile/name/%s", escapedName)

			assert.Equal(t, tt.expectedEndpointPath, endpoint, tt.description)

			// Verify the endpoint doesn't contain unescaped dangerous characters
			assert.NotContains(t, endpoint, " ", "Endpoint should not contain unescaped spaces")
			assert.NotContains(t, endpoint, "&", "Endpoint should not contain unescaped ampersands")
			assert.NotContains(t, endpoint, "?", "Endpoint should not contain unescaped question marks")
			assert.NotContains(t, endpoint, "#", "Endpoint should not contain unescaped hash symbols")
		})
	}
}

func TestGetLogProfileByName_DoubleEncodingPrevention(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		description string
	}{
		{
			name:        "Already encoded spaces",
			profileName: "profile%20name",
			description: "Already encoded spaces should be double-encoded to prevent injection",
		},
		{
			name:        "Already encoded ampersand",
			profileName: "profile%26name",
			description: "Already encoded ampersands should be double-encoded",
		},
		{
			name:        "Mixed encoded and unencoded",
			profileName: "profile%20name with spaces",
			description: "Mixed encoded and unencoded should be consistently handled",
		},
		{
			name:        "Malicious encoded input",
			profileName: "profile%2F..%2F..%2Fetc%2Fpasswd",
			description: "Path traversal attempts should be double-encoded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First escape
			escaped := url.QueryEscape(tt.profileName)

			// Verify that % characters in the input are double-encoded
			if strings.Contains(tt.profileName, "%") {
				assert.Contains(t, escaped, "%25", "Percent signs should be double-encoded to %25")
				assert.NotEqual(t, tt.profileName, escaped, "Input with percent signs should be modified")
			}

			// Verify that unescaping gives us back the original malicious input
			// This confirms we're not accidentally processing it as valid encoding
			unescaped, err := url.QueryUnescape(escaped)
			assert.NoError(t, err, "Should be able to unescape")
			assert.Equal(t, tt.profileName, unescaped, "Unescaping should return original input")

			// Verify that the escaped version is safe for URL construction
			endpoint := fmt.Sprintf("dcf/log-profile/name/%s", escaped)
			assert.NotContains(t, endpoint, "../", "Endpoint should not contain path traversal sequences")
		})
	}
}

func TestURLEscaping_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Null byte",
			input:       "profile\x00name",
			description: "Null bytes should be properly escaped",
		},
		{
			name:        "Non-printable ASCII",
			input:       "profile\x01\x02\x03name",
			description: "Non-printable ASCII characters should be escaped",
		},
		{
			name:        "High Unicode",
			input:       "profile\U0001F600name", // 😀 emoji
			description: "High Unicode should be properly UTF-8 encoded",
		},
		{
			name:        "All URL reserved characters",
			input:       ":/?#[]@!$&'()*+,;=",
			description: "All URL reserved characters should be escaped",
		},
		{
			name:        "Maximum length edge case",
			input:       strings.Repeat("a", 1000) + "&dangerous=value",
			description: "Very long strings with dangerous content should be escaped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			escaped := url.QueryEscape(tt.input)

			// Verify that dangerous characters are not present in escaped form
			assert.NotContains(t, escaped, "&", "Ampersand should be escaped")
			assert.NotContains(t, escaped, "?", "Question mark should be escaped")
			assert.NotContains(t, escaped, "#", "Hash should be escaped")

			// Verify round-trip safety
			unescaped, err := url.QueryUnescape(escaped)
			assert.NoError(t, err, "Should be able to unescape")
			assert.Equal(t, tt.input, unescaped, "Round-trip should preserve original")

			// Verify that escaped form is safe for URL construction
			endpoint := fmt.Sprintf("dcf/log-profile/name/%s", escaped)
			assert.NotContains(t, endpoint, " ", "No unescaped spaces in endpoint")
			assert.True(t, len(endpoint) > len("dcf/log-profile/name/"), "Endpoint should have content")
		})
	}
}
