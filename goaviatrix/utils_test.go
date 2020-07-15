package goaviatrix

import "testing"

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
