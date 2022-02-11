package goaviatrix

import (
	"errors"
	"strings"
)

func (c *Client) ControllerVersionValidation(supportedVersions []string) error {
	if len(supportedVersions) == 0 {
		return errors.New("supportedVersions is not provided")
	}

	currentVersion, _, err := c.GetCurrentVersion()
	if err != nil {
		return err
	}
	currVersion := strings.Split(currentVersion, ".")
	if len(currVersion) < 2 {
		return errors.New("couldn't get current version correctly")
	}

	for i := 0; i < len(supportedVersions); i++ {
		suppVersion := strings.Split(supportedVersions[i], ".")
		if len(suppVersion) < 2 {
			return errors.New("" + supportedVersions[i] + " is not set correctly, correct example: '5.1'")
		}
		if suppVersion[0] == currVersion[0] && suppVersion[1] == currVersion[1] {
			return nil
		}
	}

	return errors.New("current Terraform branch does not support controller version: UserConnect-" + currentVersion +
		". Please go to 'https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-compatibility' for version construct instructions")
}
