package goaviatrix

import (
	"errors"
	"fmt"

	"golang.org/x/mod/semver"
)

const (
	helpURL = "https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-compatibility"
)

func (c *Client) ControllerVersionValidation(supportedVersions []string) error {
	if len(supportedVersions) == 0 {
		return errors.New("supportedVersions is not provided")
	}

	currentVersion, err := c.GetCurrentVersion()
	if err != nil {
		return err
	}
	if err := isVersionSupported(currentVersion, supportedVersions); err != nil {
		return fmt.Errorf(
			"current Terraform version does not support controller version: %s, Please see %s for a list of compatible controller versions.",
			currentVersion, helpURL,
		)
	}

	return nil
}

// isVersionSupported compares the current version against a list of supported versions.
// It returns nil if the current version is supported, otherwise an error.
func isVersionSupported(currentVersion string, supportedVersions []string) error {
	version := "v" + currentVersion
	for _, supported := range supportedVersions {
		supported = "v" + supported
		if semver.Compare(semver.MajorMinor(version), semver.MajorMinor(supported)) == 0 {
			return nil
		}
	}
	return fmt.Errorf("version %s is not supported", currentVersion)
}
