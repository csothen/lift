package db

import (
	"context"
	"fmt"
)

func (q *querier) ValidateAPIKey(ctx context.Context, key string) error {
	var apiKey APIKey
	res := q.db.First(&apiKey, "value = ?", key)
	if res.Error != nil {
		return fmt.Errorf("failed to retrieve the api key: %w", res.Error)
	}

	if apiKey.Value != key {
		return fmt.Errorf("key does not match")
	}
	return nil
}

func (q *querier) RefreshAPIKey(ctx context.Context, newKey string) error {
	var apiKey APIKey
	res := q.db.First(&apiKey)
	if res.Error != nil {
		return fmt.Errorf("failed to retrieve the api key: %w", res.Error)
	}

	apiKey.Value = newKey

	q.db.Save(&apiKey)
	return nil
}
