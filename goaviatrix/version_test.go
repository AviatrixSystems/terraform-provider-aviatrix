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
				Major:        6,
				Minor:        3,
				Build:        2309,
				MinorBuildID: "patch",
				HasBuild:     true,
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

func TestCompareSoftwareVersions(t *testing.T) {
	tt := []struct {
		name    string
		a       string
		b       string
		want    int
		wantErr bool
	}{
		{
			name: "test01",
			a:    "6.5.100",
			b:    "6.5.101",
			want: -1,
		},
		{
			name: "test02",
			a:    "6.5.102",
			b:    "6.5.101",
			want: 1,
		},
		{
			name: "test03",
			a:    "6.5.100",
			b:    "6.5.100",
			want: 0,
		},
		{
			name: "test04",
			a:    "6.5.100",
			b:    "6.5",
			want: -1,
		},
		{
			name: "test05",
			a:    "6.6.100",
			b:    "6.5",
			want: 1,
		},
		{
			name: "test06",
			a:    "6.6",
			b:    "6.5",
			want: 1,
		},
		{
			name: "test07",
			a:    "6.6",
			b:    "6.5.100",
			want: 1,
		},
		{
			name: "test08",
			a:    "6.5.1232",
			b:    "6.5-patch.100",
			want: -1,
		},
		{
			name: "test09",
			a:    "6.5-patch.1232",
			b:    "6.5-patch.100",
			want: 1132,
		},
		{
			name:    "test10",
			a:       "6",
			b:       "6.5-patch.100",
			wantErr: true,
		},
		{
			name:    "test11",
			a:       "6.",
			b:       "6.5-patch.100",
			wantErr: true,
		},
		{
			name:    "test12",
			a:       "6.a",
			b:       "6.5-patch.100",
			wantErr: true,
		},
		{
			name: "test13",
			a:    "6.5-patch.1232",
			b:    "6.5-patch.100",
			want: 1132,
		},
		{
			name: "test14",
			a:    "6.5-patch.1232",
			b:    "6.6-patch.100",
			want: -1,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CompareSoftwareVersions(tc.a, tc.b)
			if tc.wantErr && err == nil {
				t.Fatal("wantErr but got nil err")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("don't wantErr but got err (%v)", err)
			}
			if tc.want != got {
				t.Fatalf("want (%d) does not equal got (%d)", tc.want, got)
			}
		})
	}
}
