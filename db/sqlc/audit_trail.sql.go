// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: audit_trail.sql

package db

import (
	"context"

	"github.com/sqlc-dev/pqtype"
)

const logAction = `-- name: LogAction :exec
INSERT INTO store_audit_trail (
  store_id, user_id, action, details
) VALUES (
  $1, $2, $3, $4
)
`

type LogActionParams struct {
	StoreID int64                 `json:"store_id"`
	UserID  int64                 `json:"user_id"`
	Action  string                `json:"action"`
	Details pqtype.NullRawMessage `json:"details"`
}

func (q *Queries) LogAction(ctx context.Context, arg LogActionParams) error {
	_, err := q.db.ExecContext(ctx, logAction,
		arg.StoreID,
		arg.UserID,
		arg.Action,
		arg.Details,
	)
	return err
}
