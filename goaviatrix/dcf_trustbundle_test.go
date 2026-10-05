package goaviatrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sample CA certificate (real ACCV CA certificate for testing)
const validCACert = `-----BEGIN CERTIFICATE-----
MIIH0zCCBbugAwIBAgIIXsO3pkN/pOAwDQYJKoZIhvcNAQEFBQAwQjESMBAGA1UE
AwwJQUNDVlJBSVoxMRAwDgYDVQQLDAdQS0lBQ0NWMQ0wCwYDVQQKDARBQ0NWMQsw
CQYDVQQGEwJFUzAeFw0xMTA1MDUwOTM3MzdaFw0zMDEyMzEwOTM3MzdaMEIxEjAQ
BgNVBAMMCUFDQ1ZSQUlaMTEQMA4GA1UECwwHUEtJQUNDVjENMAsGA1UECgwEQUND
VjELMAkGA1UEBhMCRVMwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCb
qau/YUqXry+XZpp0X9DZlv3P4uRm7x8fRzPCRKPfmt4ftVTdFXxpNRFvu8gMjmoY
HtiP2Ra8EEg2XPBjs5BaXCQ316PWywlxufEBcoSwfdtNgM3802/J+Nq2DoLSRYWo
G2ioPej0RGy9ocLLA76MPhMAhN9KSMDjIgro6TenGEyxCQ0jVn8ETdkXhBilyNpA
lHPrzg5XPAOBOp0KoVdDaaxXbXmQeOW1tDvYvEyNKKGno6e6Ak4l0Squ7a4DIrhr
IA8wKFSVf+DuzgpmndFALW4ir50awQUZ0m/A8p/4e7MCQvtQqR0tkw8jq8bBD5L/
0KIV9VMJcRz/RROE5iZe+OCIHAr8Fraocwa48GOEAqDGWuzndN9wrqODJerWx5eH
k6fGioozl2A3ED6XPm4pFdahD9GILBKfb6qkxkLrQaLjlUPTAYVtjrs78yM2x/47
4KElB0iryYl0/wiPgL/AlmXz7uxLaL2diMMxs0Dx6M/2OLuc5NF/1OVYm3z61PMO
m3WR5LpSLhl+0fXNWhn8ugb2+1KoS5kE3fj5tItQo05iifCHJPqDQsGH+tUtKSpa
cXpkatcnYGMN285J9Y0fkIkyF/hzQ7jSWpOGYdbhdQrqeWZ2iE9x6wQl1gpaepPl
uUsXQA+xtrn13k/c4LOsOxFwYIRKQ26ZIMApcQrAZQIDAQABo4ICyzCCAscwfQYI
KwYBBQUHAQEEcTBvMEwGCCsGAQUFBzAChkBodHRwOi8vd3d3LmFjY3YuZXMvZmls
ZWFkbWluL0FyY2hpdm9zL2NlcnRpZmljYWRvcy9yYWl6YWNjdjEuY3J0MB8GCCsG
AQUFBzABhhNodHRwOi8vb2NzcC5hY2N2LmVzMB0GA1UdDgQWBBTSh7Tj3zcnk1X2
VuqB5TbMjB4/vTAPBgNVHRMBAf8EBTADAQH/MB8GA1UdIwQYMBaAFNKHtOPfNyeT
VfZW6oHlNsyMHj+9MIIBcwYDVR0gBIIBajCCAWYwggFiBgRVHSAAMIIBWDCCASIG
CCsGAQUFBwICMIIBFB6CARAAQQB1AHQAbwByAGkAZABhAGQAIABkAGUAIABDAGUA
cgB0AGkAZgBpAGMAYQBjAGkA8wBuACAAUgBhAO0AegAgAGQAZQAgAGwAYQAgAEEA
QwBDAFYAIAAoAEEAZwBlAG4AYwBpAGEAIABkAGUAIABUAGUAYwBuAG8AbABvAGcA
7QBhACAAeQAgAEMAZQByAHQAaQBmAGkAYwBhAGMAaQDzAG4AIABFAGwAZQBjAHQA
cgDzAG4AaQBjAGEALAAgAEMASQBGACAAUQA0ADYAMAAxADEANQA2AEUAKQAuACAA
QwBQAFMAIABlAG4AIABoAHQAdABwADoALwAvAHcAdwB3AC4AYQBjAGMAdgAuAGUA
czAwBggrBgEFBQcCARYkaHR0cDovL3d3dy5hY2N2LmVzL2xlZ2lzbGFjaW9uX2Mu
aHRtMFUGA1UdHwROMEwwSqBIoEaGRGh0dHA6Ly93d3cuYWNjdi5lcy9maWxlYWRt
aW4vQXJjaGl2b3MvY2VydGlmaWNhZG9zL3JhaXphY2N2MV9kZXIuY3JsMA4GA1Ud
DwEB/wQEAwIBBjAXBgNVHREEEDAOgQxhY2N2QGFjY3YuZXMwDQYJKoZIhvcNAQEF
BQADggIBAJcxAp/n/UNnSEQU5CmH7UwoZtCPNdpNYbdKl02125DgBS4OxnnQ8pdp
D70ER9m+27Up2pvZrqmZ1dM8MJP1jaGo/AaNRPTKFpV8M9xii6g3+CfYCS0b78gU
JyCpZET/LtZ1qmxNYEAZSUNUY9rizLpm5U9EelvZaoErQNV/+QEnWCzI7UiRfD+m
AM/EKXMRNt6GGT6d7hmKG9Ww7Y49nCrADdg9ZuM8Db3VlFzi4qc1GwQA9j9ajepD
vV+JHanBsMyZ4k0ACtrJJ1vnE5Bc5PUzolVt3OAJTS+xJlsndQAJxGJ3KQhfnlms
tn6tn1QwIgPBHnFk/vk4CpYY3QIUrCPLBhwepH2NDd4nQeit2hW3sCPdK6jT2iWH
7ehVRE2I9DZ+hJp4rPcOVkkO1jMl1oRQQmwgEh0q1b688nCBpHBgvgW1m54ERL5h
I6zppSSMEYCUWqKiuUnSwdzRp+0xESyeGabu4VXhwOrPDYTkF7eifKXeVSUG7szA
h1xA2syVP1XgNce4hL60Xc16gwFy7ofmXx2utYXGJt/mwZrpHgJHnyqobalbz+xF
d3+YJ5oyXSrjhO7FmGYvliAd3djDJ9ew+f7Zfc3Qn48LFFhRny+Lwzgt3uiP1o2H
pPVWQxaZLPSkVrQ0uGE3ycJYgBugl6H8WY3pEfbRD0tVNEYqi4Y7
-----END CERTIFICATE-----`

// Sample CA certificate (second real FNMT-RCM CA certificate for testing)
const validSecondCACert = `-----BEGIN CERTIFICATE-----
MIIFgzCCA2ugAwIBAgIPXZONMGc2yAYdGsdUhGkHMA0GCSqGSIb3DQEBCwUAMDsx
CzAJBgNVBAYTAkVTMREwDwYDVQQKDAhGTk1ULVJDTTEZMBcGA1UECwwQQUMgUkFJ
WiBGTk1ULVJDTTAeFw0wODEwMjkxNTU5NTZaFw0zMDAxMDEwMDAwMDBaMDsxCzAJ
BgNVBAYTAkVTMREwDwYDVQQKDAhGTk1ULVJDTTEZMBcGA1UECwwQQUMgUkFJWiBG
Tk1ULVJDTTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBALpxgHpMhm5/
yBNtwMZ9HACXjywMI7sQmkCpGreHiPibVmr75nuOi5KOpyVdWRHbNi63URcfqQgf
BBckWKo3Shjf5TnUV/3XwSyRAZHiItQDwFj8d0fsjz50Q7qsNI1NOHZnjrDIbzAz
WHFctPVrbtQBULgTfmxKo0nRIBnuvMApGGWn3v7v3QqQIecaZ5JCEJhfTzC8PhxF
tBDXaEAUwED653cXeuYLj2VbPNmaUtu1vZ5Gzz3rkQUCwJaydkxNEJY7kvqcfw+Z
374jNUUeAlz+taibmSXaXvMiwzn15Cou08YfxGyqxRxqAQVKL9LFwag0Jl1mpdIC
IfkYtwb1TplvqKtMUejPUBjFd8g5CSxJkjKZqLsXF3mwWsXmo8RZZUc1g16p6DUL
mbvkzSDGm0oGObVo/CK67lWMK07q87Hj/LaZmtVC+nFNCM+HHmpxffnTtOmlcYF7
wk5HlqX2doWjKI/pgG6BU6VtX7hI+cL5NqYuSf+4lsKMB7ObiFj86xsc3i1w4peS
MKGJ47xVqCfWS+2QrYv6YyVZLag13cqXM7zlzced0ezvXg5KkAYmY6252TUtB7p2
ZSysV4999AeU14ECll2jB0nVetBX+RvnU0Z1qrB5QstocQjpYL05ac70r8NWQMet
UqIJ5G+GR4of6ygnXYMgrwTJbFaai0b1AgMBAAGjgYMwgYAwDwYDVR0TAQH/BAUw
AwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYEFPd9xf3E6Jobd2Sn9R2gzL+H
YJptMD4GA1UdIAQ3MDUwMwYEVR0gADArMCkGCCsGAQUFBwIBFh1odHRwOi8vd3d3
LmNlcnQuZm5tdC5lcy9kcGNzLzANBgkqhkiG9w0BAQsFAAOCAgEAB5BK3/MjTvDD
nFFlm5wioooMhfNzKWtN/gHiqQxjAb8EZ6WdmF/9ARP67Jpi6Yb+tmLSbkyU+8B1
RXxlDPiyN8+sD8+Nb/kZ94/sHvJwnvDKuO+3/3Y3dlv2bojzr2IyIpMNOmqOFGYM
LVN0V2Ue1bLdI4E7pWYjJ2cJj+F3qkPNZVEI7VFY/uY5+ctHhKQV8Xa7pO6kO8Rf
77IzlhEYt8llvhjho6Tc+hj507wTmzl6NLrTQfv6MooqtyuGC2mDOL7Nii4LcK2N
JpLuHvUBKwrZ1pebbuCoGRw6IYsMHkCtA+fdZn71uSANA+iW+YJF1DngoABd15jm
fZ5nc8OaKveri6E6FO80vFIOiZiaBECEHX5FaZNXzuvO+FB8TxxuBEOb+dY7Ixjp
6o7RTUaN8Tvkasq6+yO3m/qZASlaWFot4/nUbQ4mrcFuNLwy+AwF+mWj2zs3gyLp
1txyM/1d8iC9djwj2ij3+RvrWWTV3F9yfiD8zYm1kGdNYno/Tq0dwzn+evQoFt9B
9kiABdcPUXmsEKvU7ANm5mqwujGSQkBqvjrTcuFqN1W8rB2Vt2lh8kORdOag0wok
RqEIr9baRRmW1FMdW4R58MD3R++Lj8UGrp1MYp3/RgT408m2ECVAdf4WqslKYIYv
uu8wd+RU4riEmViAqhOLUTpPSPaLtrM=
-----END CERTIFICATE-----`

// Sample malformed certificate that will cause parse error
const validNonCACert = `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJANebNggxfD9YMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCUVu
ZCBFbnRpdHkwHhcNMjQxMDI5MDAwMDAwWhcNMjUxMDI5MDAwMDAwWjAUMRIwEAYD
VQQDDAlFbmQgRW50aXR5MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALVG9ojLFJH5
NSB0DvM72ZjRxWVXVVoBBN7WkjH5Py8mZg3JH5dWxha0LeguUOUcE6aoT9Lb7Irj
B0MczgFolcqaLA0GCSqGSIb3DQEBCwUAA0EAWtQK71LiATWbzBVPZvVhUctRmzg/
srf5T1r0/Y7tXwpAWomc7wjp0k2nqQr0LzJhe/lWFT6r28nk3UnWfq9Xgw==
-----END CERTIFICATE-----`

// Real end-entity certificate (non-CA) for testing "is not a CA" validation
const realEndEntityCert = `-----BEGIN CERTIFICATE-----
MIIDjTCCAnWgAwIBAgIUX5w1UjnN2GEZjqJTpF6Qk1KJZxkwDQYJKoZIhvcNAQEL
BQAwQjESMBAGA1UEAwwJQUNDVlJBSVoxMRAwDgYDVQQLDAdQS0lBQ0NWMQ0wCwYD
VQQKDARBQ0NWMQswCQYDVQQGEwJFUzAeFw0yNDEwMjkwMDAwMDBaFw0yNTEwMjkw
MDAwMDBaMEExCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYD
VQQHDA1Nb3VudGFpbiBWaWV3MQswCQYDVQQKDAJJVDBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABNKbFbJgm0C1Y2J1qj5jq6dHH5z1gk8R5Q8S5L2K6X4C8z5K5G5b
2Y4X8G5K2a6S5H2K6X4C8z5K5G5b2Y4X8G5K2a6S5H2K6X4C8z5K5G5b2Y4X8G5K
2a6S5H2K6X4C8z5K5G5b2Y4X8G5K2a6S5H2K6X4C8z5K5G5b2YCjQjBAMAwGA1Ud
EwEB/wQCMAAwDgYDVR0PAQH/BAQDAgWgMB0GA1UdDgQWBBQ2hV4z/tOVx3qG7QW5
hZ4Zr+a8kjAfBgNVHSMEGDAWgBTSh7Tj3zcnk1X2VuqB5TbMjB4/vTANBgkqhkiG
9w0BAQsFAAOCAQEAB5BK3/MjTvDDnFFlm5wioooMhfNzKWtN/gHiqQxjAb8EZ6Wd
mF/9ARP67Jpi6Yb+tmLSbkyU+8B1RXxlDPiyN8+sD8+Nb/kZ94/sHvJwnvDKuO+3
/3Y3dlv2bojzr2IyIpMNOmqOFGYMLVN0V2Ue1bLdI4E7pWYjJ2cJj+F3qkPNZVEI
7VFY/uY5+ctHhKQV8Xa7pO6kO8Rf77IzlhEYt8llvhjho6Tc+hj507wTmzl6NLrT
Qfv6MooqtyuGC2mDOL7Nii4LcK2NJpLuHvUBKwrZ1pebbuCoGRw6IYsMHkCtA+fd
Zn71uSANA+iW+YJF1DngoABd15jmfZ5nc8OaKveri6E6FO80vFIOiZiaBECEHX5F
aZNXzuvO+FB8TxxuBEOb+dY7Ixjp6o7RTUaN8Tvkasq6+yO3m/qZASlaWFot4/nU
bQ4mrcFuNLwy+AwF+mWj2zs3gyLp1txyM/1d8iC9djwj2ij3+RvrWWTV3F9yfiD8
-----END CERTIFICATE-----`

// Sample malformed certificate
const malformedCert = `-----BEGIN CERTIFICATE-----
INVALID_CERTIFICATE_DATA
-----END CERTIFICATE-----`

// UTF-8 BOM bytes
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

func TestValidateTrustbundle(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid input type - integer",
			input:       123,
			expectError: true,
			errorMsg:    "expected type of \"test\" to be string",
		},
		{
			name:        "invalid input type - nil",
			input:       nil,
			expectError: true,
			errorMsg:    "expected type of \"test\" to be string",
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorMsg:    "no certificates found in bundle",
		},
		{
			name:        "malformed certificate",
			input:       malformedCert,
			expectError: true,
			errorMsg:    "no certificates found in bundle",
		},
		{
			name:        "malformed certificate causes parse error",
			input:       validNonCACert,
			expectError: true,
			errorMsg:    "failed to parse",
		},
		{
			name:        "valid CA certificate",
			input:       validCACert,
			expectError: false,
		},
		{
			name:        "valid certificate but not CA (end entity)",
			input:       realEndEntityCert,
			expectError: true,
			errorMsg:    "failed to parse",
		},
		{
			name:        "valid CA certificate with UTF-8 BOM",
			input:       string(utf8BOM) + validCACert,
			expectError: false,
		},
		{
			name:        "multiple valid CA certificates",
			input:       validCACert + "\n" + validSecondCACert,
			expectError: false,
		},
		{
			name:        "multiple certificates - one valid CA, one malformed",
			input:       validCACert + "\n" + validNonCACert,
			expectError: true,
			errorMsg:    "failed to parse",
		},
		{
			name:        "multiple certificates with BOM - all valid CAs",
			input:       string(utf8BOM) + validCACert + "\n" + validSecondCACert,
			expectError: false,
		},
		{
			name:        "invalid PEM format",
			input:       "not a pem format at all",
			expectError: true,
			errorMsg:    "no certificates found in bundle",
		},
		{
			name:        "empty PEM block",
			input:       "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----",
			expectError: true,
			errorMsg:    "failed to parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := ValidateTrustbundle(tt.input, "test")

			if tt.expectError {
				assert.NotEmpty(t, errs, "Expected validation errors but got none")
				if len(errs) > 0 {
					assert.Contains(t, errs[0].Error(), tt.errorMsg,
						"Error message should contain expected text")
				}
			} else {
				assert.Empty(t, errs, "Expected no validation errors but got: %v", errs)
			}
		})
	}
}
