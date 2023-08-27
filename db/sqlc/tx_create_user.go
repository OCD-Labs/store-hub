package db

import "context"

// A CreateUserTxParams contains the input parameters of
// the create user transaction.
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// A CreateUserTxResult contains the result of the create user transaction.
type CreateUserTxResult struct {
	User User
}

// CreateUserTx creates a user row and schedules a verify email task on redis.
func (dbTx *SQLTx) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := dbTx.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
