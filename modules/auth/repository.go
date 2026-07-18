package auth

import (
	"context"
	"github/anurag/altar-be/system/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct{}

func NewRepository() *AuthRepository {
	return &AuthRepository{}
}

// Users
func (r *AuthRepository) CreateUser(ctx context.Context, db database.DBTX, user *User) error {
	query := `INSERT INTO users
			(
				id,
				organization_id,
				email, 
				full_name,
				avatar_url,
				is_email_verified,
				last_login_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			`

	_, err := db.Exec(ctx, query,
		user.ID,
		user.OrganizationID,
		user.Email,
		user.FullName,
		user.AvatarURL,
		user.IsEmailVerified,
		user.LastLoginAt)

	return err
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, db database.DBTX, email string) (*User, error) {
	query := `SELECT
				id,
				organization_id,
				email,
				full_name,
				avatar_url,
				is_email_verified,
				last_login_at,
				created_at,
				updated_at
			FROM users WHERE email = $1
			`

	user := &User{}

	err := db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.FullName,
		&user.AvatarURL,
		&user.IsEmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *AuthRepository) FindUserByID(ctx context.Context, db database.DBTX, userID uuid.UUID) (*User, error) {

	query := `
		SELECT
			id,
			organization_id,
			email,
			full_name,
			avatar_url,
			is_email_verified,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}

	err := db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.FullName,
		&user.AvatarURL,
		&user.IsEmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateLastLoginAt(ctx context.Context, db database.DBTX, userID uuid.UUID) error {
	query := `UPDATE users
				SET
				last_login_at = NOW(),
				updated_At = NOW()
			WHERE id = $1`

	_, err := db.Exec(ctx, query, userID)
	return err
}

// User Credentials
func (r *AuthRepository) CreateUserCredential(ctx context.Context, db database.DBTX, credential *UserCredential) error {
	query := `INSERT INTO user_credentials
			(
				user_id,
				password_hash
			)
			VALUES ($1, $2)
			`

	_, err := db.Exec(ctx, query,
		credential.UserID,
		credential.PasswordHash,
	)

	return err
}

func (r *AuthRepository) FindUserCredential(ctx context.Context, db database.DBTX, userID uuid.UUID) (*UserCredential, error) {
	query := `SELECT
				user_id,
				password_hash,
				created_at
			FROM user_credentials
			WHERE user_id = $1
			`

	credential := &UserCredential{}

	err := db.QueryRow(ctx, query, userID).Scan(
		&credential.UserID,
		&credential.PasswordHash,
		&credential.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return credential, nil
}

// OAuth Accounts
func (r *AuthRepository) CreateAuthAccount(ctx context.Context, db database.DBTX, account *AuthAccount) error {
	query := `INSERT INTO auth_accounts
			(
				id,
				user_id,
				provider,
				provider_user_id
			)
			VALUES ($1, $2, $3, $4)
			`

	_, err := db.Exec(ctx, query,
		account.ID,
		account.UserID,
		account.Provider,
		account.ProviderUserID,
	)

	return err
}

func (r *AuthRepository) FindAuthAccount(ctx context.Context, db database.DBTX, provider string, providerUserID string) (*AuthAccount, error) {
	query := `SELECT
				id,
				user_id,
				provider,
				provider_user_id,
				created_at,
				updated_at
			FROM auth_accounts
			WHERE provider = $1 AND provider_user_id = $2
			`

	authAccount := &AuthAccount{}

	err := db.QueryRow(ctx, query,
		provider,
		providerUserID,
	).Scan(
		&authAccount.ID,
		&authAccount.UserID,
		&authAccount.Provider,
		&authAccount.ProviderUserID,
		&authAccount.CreatedAt,
		&authAccount.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return authAccount, nil
}

// Refresh Tokens
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, db database.DBTX, token *RefreshToken) error {
	query := `INSERT INTO refresh_tokens
			(
				id,
				user_id,
				token_hash,
				expires_at,
				revoked
			)
			VALUES ($1, $2, $3, $4, $5)
			`

	_, err := db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.Revoked,
	)

	return err
}

func (r *AuthRepository) FindRefreshTokenByHash(ctx context.Context, db database.DBTX, hash []byte) (*RefreshToken, error) {
	query := `SELECT
				id,
				user_id,
				token_hash,
				expires_at,
				revoked,
				created_at,
				updated_at
			FROM refresh_tokens
			WHERE token_hash = $1
			`

	refreshToken := &RefreshToken{}

	err := db.QueryRow(ctx, query,
		hash).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, db database.DBTX, id uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET
			revoked = TRUE,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)

	return err
}

func (r *AuthRepository) RevokeAllRefreshTokens(ctx context.Context, db database.DBTX, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens
				SET
				revoked = TRUE,
				updated_at = NOW()
				WHERE user_id = $1 AND revoked = FALSE
			`

	_, err := db.Exec(ctx, query, userID)

	return err
}
