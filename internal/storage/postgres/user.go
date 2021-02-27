package postgres

import "context"

func findUserByID(ctx context.Context, tx *Tx, id string) (*User, error) {
	sqlStmt := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := tx.QueryRow(ctx, sqlStmt, id)

	var user User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func findUserExternalLoginBySourceID(ctx context.Context, tx *Tx, source, sourceID string) (*UserExternalLogin, error) {
	sqlStmt := `
		SELECT id, user_id, source, source_id, created_at, updated_at
		FROM user_external_logins
		WHERE source = $1 AND source_id = $2
	`

	row := tx.QueryRow(ctx, sqlStmt, source, sourceID)

	var externalLogin UserExternalLogin
	if err := row.Scan(
		&externalLogin.ID,
		&externalLogin.UserID,
		&externalLogin.Source,
		&externalLogin.SourceID,
		&externalLogin.CreatedAt,
		&externalLogin.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &externalLogin, nil
}

func createUser(ctx context.Context, tx *Tx, user *User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO users(id, name, email, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func createUserExternalLogin(ctx context.Context, tx *Tx, userExternalLogin *UserExternalLogin) error {
	userExternalLogin.CreatedAt = tx.now
	userExternalLogin.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO user_external_logins(id, user_id, source, source_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		userExternalLogin.ID,
		userExternalLogin.UserID,
		userExternalLogin.Source,
		userExternalLogin.SourceID,
		userExternalLogin.CreatedAt,
		userExternalLogin.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
