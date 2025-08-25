package memstore

import (
	"context"
)

func (s *Store) UpsertIdentity(ctx context.Context, transport, externalID string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return getUserID(s, transport, externalID)
}

func getUserID(s *Store, transport, externalID string) (int64, error) {
	key := transport + ":" + externalID

	userID, exists := s.identities[key]

	if !exists {
		userID = s.nextUserID
		s.nextUserID++
		s.identities[key] = userID
	}

	return userID, nil
}
