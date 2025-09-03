package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/tuor4eg/ip_accounting_bot/internal/crypto"
	"github.com/tuor4eg/ip_accounting_bot/internal/domain"
	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

// InsertIncome inserts a single income record.
// 'amount' is in minor currency units (e.g., kopecks), must be >= 0.
// 'at' is the income date; only the date part is stored (cast to DATE in SQL).
func (s *Store) UpsertIdentity(ctx context.Context, transport, externalID string, chatID int64) (int64, error) {
	const op = "postgres.UpsertIdentity"

	if err := validate.ValidateTransport(transport); err != nil {
		return 0, validate.Wrap(op, err)
	}
	if err := validate.ValidateExternalID(externalID); err != nil {
		return 0, validate.Wrap(op, err)
	}

	// Prepare HMAC(external_id) as per new schema
	extHash := s.ExternalHash(transport, externalID)

	// AAD ties ciphertext to transport (prevents cross-transport reuse)
	aad := []byte(strings.ToLower(transport))

	var uid int64

	err := s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// 1) Check if mapping already exists.
		if err := tx.QueryRow(ctx,
			`SELECT user_id
			   FROM user_identities
			  WHERE transport = $1 AND external_hash = $2 AND hmac_kid = $3`,
			transport, extHash, s.GetHMACKid(),
		).Scan(&uid); err == nil { // Optional: refresh pii.telegram if chatID provided (e.g., chat migrated)
			if chatID != 0 {
				if err := upsertPIITelegram(ctx, tx, s.GetAEADBox(), s.GetAEADKid(), uid, chatID, aad); err != nil {
					return validate.Wrap(op, err)
				}
			}
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
			INSERT INTO user_identities (user_id, transport, external_hash, hmac_kid)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (transport, external_hash, hmac_kid) DO NOTHING
			RETURNING user_id
		`, newUserID, transport, extHash, s.GetHMACKid()).Scan(&insertedUserID)

		switch {
		case err == nil:
			uid = insertedUserID

		case errors.Is(err, pgx.ErrNoRows):
			// Conflict path: identity was inserted concurrently. Clean up the orphan user and fetch the right uid.
			if _, derr := tx.Exec(ctx, `DELETE FROM users WHERE id=$1`, newUserID); derr != nil {
				// best-effort cleanup; ignore error to proceed fetching the correct uid
			}

			if err := tx.QueryRow(ctx,
				`SELECT user_id
				   FROM user_identities
				  WHERE transport = $1 AND external_hash = $2 AND hmac_kid = $3`,
				transport, extHash, s.GetHMACKid(),
			).Scan(&uid); err != nil {
				return validate.Wrap(op, err)
			}

		default:
			return validate.Wrap(op, err)
		}

		// 4) Upsert encrypted chat_id into pii.telegram (if provided)
		if chatID != 0 {
			if err := upsertPIITelegram(ctx, tx, s.GetAEADBox(), s.GetAEADKid(), uid, chatID, aad); err != nil {
				return validate.Wrap(op, err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, validate.Wrap(op, err)
	}
	return uid, nil
}

func (s *Store) GetUserScheme(ctx context.Context, userID int64) (domain.TaxScheme, error) {
	const op = "postgres.GetUserScheme"

	if err := validate.ValidateUserID(userID); err != nil {
		return "", validate.Wrap(op, err)
	}

	var scheme domain.TaxScheme

	if err := s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT tax_scheme FROM users WHERE id=$1 LIMIT 1`, userID).Scan(&scheme)
	}); err != nil {
		return "", validate.Wrap(op, err)
	}

	return scheme, nil
}

// upsertPIITelegram encrypts and upserts chatID into pii.telegram for given user.
func upsertPIITelegram(ctx context.Context, tx pgx.Tx, box *crypto.AEADBox, encKid int16, userID int64, chatID int64, aad []byte) error {
	chatEnc, err := crypto.EncryptInt64(box, chatID, aad)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO pii.telegram (user_id, chat_enc, enc_kid)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		    SET chat_enc = EXCLUDED.chat_enc,
		        enc_kid  = EXCLUDED.enc_kid,
		        updated_at = now()
	`, userID, chatEnc, encKid)
	return err
}
