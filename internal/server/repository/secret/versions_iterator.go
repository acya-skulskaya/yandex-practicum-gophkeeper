package secret

import (
	"iter"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgx/v5"
)

const versionsIteratorSelectColumns = "id, secret_id, created_at"

// versionsIterator creates an iter.Seq[SecretVersion] from sql.Rows.
// It uses iter.Seq2 to also return a potential error at the end of the iteration.
func versionsIterator(rows pgx.Rows) iter.Seq2[models.SecretVersion, error] {
	return func(yield func(models.SecretVersion, error) bool) {
		defer rows.Close() // Ensure rows are closed when the iterator finishes

		for rows.Next() {
			var secretVersion models.SecretVersion
			err := rows.Scan(&secretVersion.ID, &secretVersion.SecretID, &secretVersion.CreatedAt)
			if err != nil {
				// If scan fails, yield the error and stop iteration.
				if !yield(models.SecretVersion{}, err) {
					return
				}
			}

			// Yield the SecretVersion data. If the loop breaks early, yield returns false.
			if !yield(secretVersion, nil) {
				return
			}
		}

		// Check for any error that occurred during iteration.
		if err := rows.Err(); err != nil {
			yield(models.SecretVersion{}, err)
		}
	}
}
