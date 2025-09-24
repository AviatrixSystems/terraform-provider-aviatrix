package goaviatrix

import (
	"context"
	"fmt"
)

const (
	ipsRuleFeedsEndpoint = "dcf/ips-rule-feeds"
	ipsProfilesEndpoint  = "dcf/ips-profiles"
)

// IpsRuleFeed represents an IPS rule feed
type IpsRuleFeed struct {
	UUID        string   `json:"uuid,omitempty"`
	FeedName    string   `json:"feed_name"`
	FileContent string   `json:"file_content,omitempty"`
	IpsRules    []string `json:"ips_rules,omitempty"`
	ContentHash string   `json:"content_hash,omitempty"`
}

// IpsRuleFeedsList represents the response for listing IPS rule feeds
type IpsRuleFeedsList struct {
	IpsRuleFeeds []IpsRuleFeed `json:"ips_rule_feeds"`
}

// IpsRuleFeedUploadResponse represents the response after uploading a rule feed
type IpsRuleFeedUploadResponse struct {
	UUID        string `json:"uuid"`
	FeedName    string `json:"feed_name"`
	ContentHash string `json:"content_hash"`
}

// IpsRuleFeeds represents the rule feeds configuration in an IPS profile
type IpsRuleFeeds struct {
	CustomFeedsIds   []string `json:"custom_feeds_ids"`
	ExternalFeedsIds []string `json:"external_feeds_ids"`
	IgnoredSids      []int    `json:"ignored_sids"`
	NeverDropSids    []int    `json:"never_drop_sids"`
}

// IpsProfile represents an IPS profile
type IpsProfile struct {
	UUID             string            `json:"uuid,omitempty"`
	ProfileName      string            `json:"profile_name"`
	RuleFeeds        IpsRuleFeeds      `json:"rule_feeds"`
	IntrusionActions map[string]string `json:"intrusion_actions"`
}

// IpsProfilesList represents the response for listing IPS profiles
type IpsProfilesList struct {
	IpsProfiles []IpsProfile `json:"ips_profiles"`
}

// IpsProfileCreateResponse represents the response after creating an IPS profile
type IpsProfileCreateResponse struct {
	UUID        string `json:"uuid"`
	ProfileName string `json:"profile_name"`
}

// IPS Rule Feed methods

func (c *Client) CreateIpsRuleFeed(ctx context.Context, ruleFeed *IpsRuleFeed) (*IpsRuleFeedUploadResponse, error) {
	var response IpsRuleFeedUploadResponse

	err := c.PostAPIContext25(ctx, &response, ipsRuleFeedsEndpoint, ruleFeed)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPS rule feed: %w", err)
	}

	return &response, nil
}

func (c *Client) GetIpsRuleFeed(ctx context.Context, uuid string) (*IpsRuleFeed, error) {
	var response IpsRuleFeed

	endpoint := fmt.Sprintf("%s/%s", ipsRuleFeedsEndpoint, uuid)
	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPS rule feed: %w", err)
	}

	return &response, nil
}

func (c *Client) UpdateIpsRuleFeed(ctx context.Context, uuid string, ruleFeed *IpsRuleFeed) (*IpsRuleFeed, error) {
	endpoint := fmt.Sprintf("%s/%s", ipsRuleFeedsEndpoint, uuid)
	err := c.PutAPIContext25(ctx, endpoint, ruleFeed)
	if err != nil {
		return nil, fmt.Errorf("failed to update IPS rule feed: %w", err)
	}

	// Get updated rule feed
	return c.GetIpsRuleFeed(ctx, uuid)
}

func (c *Client) DeleteIpsRuleFeed(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("%s/%s", ipsRuleFeedsEndpoint, uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (c *Client) ListIpsRuleFeeds(ctx context.Context) (*IpsRuleFeedsList, error) {
	var response IpsRuleFeedsList

	err := c.GetAPIContext25(ctx, &response, ipsRuleFeedsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list IPS rule feeds: %w", err)
	}

	return &response, nil
}

// IPS Profile methods

func (c *Client) CreateIpsProfile(ctx context.Context, profile *IpsProfile) (*IpsProfileCreateResponse, error) {
	var response IpsProfileCreateResponse

	err := c.PostAPIContext25(ctx, &response, ipsProfilesEndpoint, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPS profile: %w", err)
	}

	return &response, nil
}

func (c *Client) GetIpsProfile(ctx context.Context, uuid string) (*IpsProfile, error) {
	var response IpsProfile

	endpoint := fmt.Sprintf("%s/%s", ipsProfilesEndpoint, uuid)
	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPS profile: %w", err)
	}

	return &response, nil
}

func (c *Client) UpdateIpsProfile(ctx context.Context, uuid string, profile *IpsProfile) (*IpsProfile, error) {
	endpoint := fmt.Sprintf("%s/%s", ipsProfilesEndpoint, uuid)
	err := c.PutAPIContext25(ctx, endpoint, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update IPS profile: %w", err)
	}

	return c.GetIpsProfile(ctx, uuid)
}

func (c *Client) DeleteIpsProfile(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("%s/%s", ipsProfilesEndpoint, uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (c *Client) ListIpsProfiles(ctx context.Context) (*IpsProfilesList, error) {
	var response IpsProfilesList

	err := c.GetAPIContext25(ctx, &response, ipsProfilesEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list IPS profiles: %w", err)
	}

	return &response, nil
}
