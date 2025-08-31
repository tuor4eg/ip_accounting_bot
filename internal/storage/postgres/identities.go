package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

// InsertIncome inserts a single income record.
// 'amount' is in minor currency units (e.g., kopecks), must be >= 0.
// 'at' is the income date; only the date part is stored (cast to DATE in SQL).
func (s *Store) UpsertIdentity(ctx context.Context, transport, externalID string) (int64, error) {
	const op = "postgres.UpsertIdentity"

	if err := validate.ValidateTransport(transport); err != nil {
		return 0, validate.Wrap(op, err)
	}
	if err := validate.ValidateExternalID(externalID); err != nil {
		return 0, validate.Wrap(op, err)
	}

	var uid int64

	err := s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// 1) Check if mapping already exists.
		if err := tx.QueryRow(ctx,
			`SELECT user_id FROM user_identities WHERE transport=$1 AND external_id=$2`,
			transport, externalID,
		).Scan(&uid); err == nil {
			return nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return validate.Wrap(op, err)
		}

		// 2) Create new user if not exists.
		var newUserID int64
		if err := tx.QueryRow(ctx, `INSERT INTO users DEFAULT VALUES RETURNING id`).Scan(&newUserID); err != nil {
			return validate.Wrap(op, err)
		}

		// 3) Try to bind identity (could race with a parallel insert).
		var insertedUserID int64
		err := tx.QueryRow(ctx, `
			INSERT INTO user_identities (user_id, transport, external_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (transport, external_id) DO NOTHING
			RETURNING user_id
		`, newUserID, transport, externalID).Scan(&insertedUserID)

		switch {
		case err == nil:
			uid = insertedUserID
			return nil

		case errors.Is(err, pgx.ErrNoRows):
			// Conflict path: identity was inserted concurrently. Clean up the orphan user and fetch the right uid.
			if _, derr := tx.Exec(ctx, `DELETE FROM users WHERE id=$1`, newUserID); derr != nil {
				// best-effort cleanup; ignore error to proceed fetching the correct uid
			}

			if err := tx.QueryRow(ctx,
				`SELECT user_id FROM user_identities WHERE transport=$1 AND external_id=$2`,
				transport, externalID,
			).Scan(&uid); err != nil {
				return validate.Wrap(op, err)
			}
			return nil

		default:
			return validate.Wrap(op, err)
		}
	})

	if err != nil {
		return 0, validate.Wrap(op, err)
	}
	return uid, nil
}
