package goaviatrix

import (
	"reflect"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
		want1   *AviatrixVersion
		wantErr bool
	}{
		{
			"test1",
			"UserConnect-6.3-patch.2309",
			"6.3",
			&AviatrixVersion{
				Major: 6,
				Minor: 3,
				Build: 2309,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseVersion() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
