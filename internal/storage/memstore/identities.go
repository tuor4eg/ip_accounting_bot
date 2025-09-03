package memstore

import (
	"context"
	"strconv"

	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
)

func (s *Store) UpsertIdentity(ctx context.Context, transport, externalID string, chatID int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return getUserID(s, transport, externalID)
}

func getUserID(s *Store, transport, externalID string) (int64, error) {
	key := transport + ":" + externalID

	user, exists := s.identities[key]

	if !exists {
		userID := s.nextUserID
		s.nextUserID++
		s.identities[key] = UserRecord{
			UserID: userID,
			Scheme: domain.TaxSchemeUSN6,
		}
		return userID, nil
	}

	return user.UserID, nil
}

func (s *Store) GetUserScheme(ctx context.Context, userID int64) (domain.TaxScheme, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.identities[strconv.FormatInt(userID, 10)].Scheme, nil
}
