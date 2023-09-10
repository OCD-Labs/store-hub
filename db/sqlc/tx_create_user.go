package db

import "context"

// A CreateUserTxParams contains the input parameters of
// the create user transaction.
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTx creates a user row and schedules a verify email task on redis.
func (dbTx *SQLTx) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error) {
	var result User

	err := dbTx.execTx(ctx, func(q *Queries) error {
		var err error

		result, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		if err := q.CreateCartForUser(ctx, result.ID); err != nil {
			return err
		}

		return arg.AfterCreate(result)
	})

	return result, err
}
