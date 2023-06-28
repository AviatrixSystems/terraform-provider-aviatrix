package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateAwsAccountNumber(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr []error
	}{
		{
			name:        "Valid AWS account number",
			input:       "123456789012",
			expectedErr: nil,
		},
		{
			name:        "Invalid AWS account number (too short)",
			input:       "12345678901",
			expectedErr: []error{fmt.Errorf("\"unitTestKey\" must be 12 digits, got: 12345678901")},
		},
		{
			name:        "Invalid AWS account number (too long)",
			input:       "1234567890123",
			expectedErr: []error{fmt.Errorf("\"unitTestKey\" must be 12 digits, got: 1234567890123")},
		},
		{
			name:        "Invalid AWS account number (letters)",
			input:       "12345678901A",
			expectedErr: []error{fmt.Errorf("\"unitTestKey\" must be 12 digits, got: 12345678901A")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			warns, errs := validateAwsAccountNumber(test.input, "unitTestKey")

			assert.Empty(t, warns)
			assert.Equal(t, test.expectedErr, errs)
		})
	}
}

func TestResourceAviatrixAccountDelete_CallsDeleteAccount(t *testing.T) {
	client := &goaviatrix.ClientInterfaceMock{
		DeleteAccountFunc: func(account *goaviatrix.Account) error {
			assert.Equal(t, "unit_test_account", account.AccountName)
			return nil
		},
	}
	defer func() { assert.Len(t, client.DeleteAccountCalls(), 1) }()

	d := schema.TestResourceDataRaw(t, resourceAviatrixAccount().Schema, map[string]interface{}{
		"account_name": "unit_test_account",
	})
	res := resourceAviatrixAccountDelete(nil, d, client)

	assert.Empty(t, res)
}

func TestResourceAviatrixAccountDelete_WhenDeleteAccountFails(t *testing.T) {
	client := &goaviatrix.ClientInterfaceMock{
		DeleteAccountFunc: func(account *goaviatrix.Account) error {
			return errors.New("controller API failure")
		},
	}
	defer func() { assert.Len(t, client.DeleteAccountCalls(), 1) }()

	d := schema.TestResourceDataRaw(t, resourceAviatrixAccount().Schema, map[string]interface{}{
		"account_name": "unit_test_account",
	})
	res := resourceAviatrixAccountDelete(nil, d, client)

	assert.Equal(t, diag.Errorf("failed to delete Aviatrix Account: controller API failure"), res)
}

func TestResourceAviatrixAccountRead_AccountWithAudit(t *testing.T) {
	client := &goaviatrix.ClientInterfaceMock{
		GetAccountFunc: func(account *goaviatrix.Account) (*goaviatrix.Account, error) {
			assert.Equal(t, "unit_test_account", account.AccountName)
			return &goaviatrix.Account{
				AccountName:      "unit_test_account",
				CloudType:        goaviatrix.AWS,
				AwsAccountNumber: "123456789012",
			}, nil
		},
		AuditAccountFunc: func(ctx context.Context, account *goaviatrix.Account) error {
			assert.Equal(t, "unit_test_account", account.AccountName)
			return nil
		},
	}
	defer func() { assert.Len(t, client.GetAccountCalls(), 1) }()
	defer func() { assert.Len(t, client.AuditAccountCalls(), 1) }()

	d := schema.TestResourceDataRaw(t, resourceAviatrixAccount().Schema, map[string]interface{}{
		"account_name":  "unit_test_account",
		"audit_account": true,
	})
	res := resourceAviatrixAccountRead(nil, d, client)

	assert.Equal(t, "unit_test_account", d.Get("account_name"))
	assert.Equal(t, 1, d.Get("cloud_type"))
	assert.Equal(t, false, d.Get("aws_iam"))

	assert.Empty(t, res)
}
