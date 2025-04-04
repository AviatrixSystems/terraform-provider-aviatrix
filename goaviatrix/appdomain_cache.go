package goaviatrix

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"
)

type AppdomainGroup struct {
	UUID     string          `json:"uuid"`
	Name     string          `json:"name"`
	Selector json.RawMessage `json:"selector"`
}

type AppdomainCache struct {
	lock      sync.Mutex
	updatedAt time.Time

	cache map[string]AppdomainGroup
}

func (a *AppdomainCache) expired() bool {
	const cacheTime = 5 * time.Second
	return time.Since(a.updatedAt) > cacheTime
}

func (a *AppdomainCache) refresh(ctx context.Context, c *Client) error {
	const endpoint = "app-domains"

	var response struct {
		Groups []AppdomainGroup `json:"app_domains"`
	}

	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return err
	}

	a.updatedAt = time.Now()
	a.cache = make(map[string]AppdomainGroup, len(response.Groups))
	for _, group := range response.Groups {
		a.cache[group.UUID] = group
	}

	return nil
}

func (a *AppdomainCache) Get(ctx context.Context, c *Client, uuid string) (AppdomainGroup, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.expired() || len(a.cache) == 0 {
		if err := a.refresh(ctx, c); err != nil {
			return AppdomainGroup{}, err
		}
	}

	group, ok := a.cache[uuid]
	if !ok {
		return AppdomainGroup{}, ErrNotFound
	}

	return group, nil
}

func (a *AppdomainCache) List(ctx context.Context, c *Client) ([]AppdomainGroup, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.expired() || len(a.cache) == 0 {
		if err := a.refresh(ctx, c); err != nil {
			return nil, err
		}
	}

	return slices.Collect(maps.Values(a.cache)), nil
}

func (a *AppdomainCache) Delete(ctx context.Context, c *Client, uuid string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.cache = nil

	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (a *AppdomainCache) Update(ctx context.Context, c *Client, uuid string, value any) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.cache = nil

	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.PutAPIContext25(ctx, endpoint, value)
}

func (a *AppdomainCache) Create(ctx context.Context, c *Client, value any) (string, error) {
	const endpoint = "app-domains"

	a.lock.Lock()
	defer a.lock.Unlock()

	a.cache = nil

	var response struct {
		UUID string `json:"uuid"`
		// More possible fields, but we don't care.
	}

	if err := c.PostAPIContext25(ctx, &response, endpoint, value); err != nil {
		return "", err
	}

	return response.UUID, nil
}
