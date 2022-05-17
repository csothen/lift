package services

import (
	"context"
	"fmt"
)

func (s *Service) ValidateAuth(ctx context.Context, auth string) error {
	err := s.repo.ValidateAPIKey(ctx, auth)
	if err != nil {
		return fmt.Errorf("failed to retrieve API Key: %w", err)
	}
	return nil
}
