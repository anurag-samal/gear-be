package auth

import (
	"context"
	"github/anurag/altar-be/system/config"
	"github/anurag/altar-be/system/packages"
	"github/anurag/altar-be/modules/auth/utils"
	"github/anurag/altar-be/system/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo   *AuthRepository
	db     *pgxpool.Pool
	redis  *redis.Client
	jwt    *pkg.JWTManager
	google *pkg.GoogleProvider
	config *config.Config
}

func NewService(repo *AuthRepository, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config, jwt *pkg.JWTManager, google *pkg.GoogleProvider) *AuthService {
	return &AuthService{
		repo:   repo,
		db:     db,
		redis:  redis,
		jwt:    jwt,
		google: google,
		config: cfg,
	}
}

// ----------HELPERS-----------------

func (s *AuthService) createUser(ctx context.Context, db database.DBTX, user *User) error {
	if user == nil {
		return auth.ErrUserNotFound
	}

	user.ID = uuid.New()
	user.LastLoginAt = nil

	return s.repo.CreateUser(ctx, db, user)
}

func (s *AuthService) createCredential(ctx context.Context, db database.DBTX, userID uuid.UUID, password string) error {

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	credential := &UserCredential{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	return s.repo.CreateUserCredential(ctx, db, credential)
}

func (s *AuthService) createRefreshSession(ctx context.Context, db database.DBTX, userID uuid.UUID, tokenPair *pkg.TokenPair) error {

	refreshToken := &RefreshToken{
		UserID:    userID,
		TokenHash: tokenPair.RefreshTokenHash,
		ExpiresAt: time.Now().UTC().Add(auth.JWT_REFRESH_TOKEN_TTL),
		Revoked:   false,
	}

	return s.repo.CreateRefreshToken(ctx, db, refreshToken)
}

func (s *AuthService) authenticateUser(ctx context.Context, db database.DBTX, user *User) (*auth.AuthResponse, string, error) {

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, "", err
	}

	if err := s.createRefreshSession(ctx, db, user.ID, tokenPair); err != nil {
		return nil, "", err
	}

	response := &auth.AuthResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(auth.JWT_ACCESS_TOKEN_TTL),
		User: auth.ProfileResponse{
			ID:              user.ID,
			OrganizationID:  user.OrganizationID,
			Email:           user.Email,
			FullName:        user.FullName,
			AvatarURL:       user.AvatarURL,
			IsEmailVerified: user.IsEmailVerified,
		},
	}

	return response, tokenPair.RefreshToken, nil
}

func (s *AuthService) buildProfileResponse(user *User) auth.ProfileResponse {
	return auth.ProfileResponse{
		ID:              user.ID,
		OrganizationID:  user.OrganizationID,
		Email:           user.Email,
		FullName:        user.FullName,
		AvatarURL:       user.AvatarURL,
		IsEmailVerified: user.IsEmailVerified,
	}
}

func (s *AuthService) buildAuthResponse(user *User, tokenPair *pkg.TokenPair) *auth.AuthResponse {

	return &auth.AuthResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(auth.JWT_ACCESS_TOKEN_TTL),
		User:        s.buildProfileResponse(user),
	}
}



//-----------SERVICES----------------

func (s *AuthService) Register(ctx context.Context, req auth.RegisterRequest) (*auth.AuthResponse, string, error) {

	// Normalize request
	auth.NormalizeRegisterRequest(&req)

	// Validate request
	if err := auth.ValidateRegisterRequest(req); err != nil {
		return nil, "", err
	}

	// Check if email already exists
	existingUser, err := s.repo.FindUserByEmail(ctx, s.db, req.Email)
	if err != nil {
		return nil, "", err
	}

	if existingUser != nil {
		return nil, "", auth.ErrEmailAlreadyExists
	}

	// Prepare user
	user := &User{
		Email:           req.Email,
		FullName:        req.FullName,
		IsEmailVerified: false,
	}

	// Generate JWTs
	tokenPair, err := s.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, "", err
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Create user
	if err := s.createUser(ctx, tx, user); err != nil {
		return nil, "", err
	}

	// Create credentials
	if err := s.createCredential(ctx, tx, user.ID, req.Password); err != nil {
		return nil, "", err
	}

	// Create refresh session
	if err := s.createRefreshSession(ctx, tx, user.ID, tokenPair); err != nil {
		return nil, "", err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return s.buildAuthResponse(user, tokenPair), tokenPair.RefreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, req auth.LoginRequest) (*auth.AuthResponse, string, error) {

	// Normalize request
	auth.NormalizeLoginRequest(&req)

	// Validate request
	if err := auth.ValidateLoginRequest(req); err != nil {
		return nil, "", err
	}

	// Find user
	user, err := s.repo.FindUserByEmail(ctx, s.db, req.Email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", auth.ErrInvalidCredentials
	}

	// Find credentials
	credential, err := s.repo.FindUserCredential(ctx, s.db, user.ID)
	if err != nil {
		return nil, "", err
	}

	if credential == nil {
		return nil, "", auth.ErrInvalidCredentials
	}

	// Verify password
	if err := auth.VerifyPassword(credential.PasswordHash, req.Password); err != nil {
		return nil, "", auth.ErrInvalidCredentials
	}

	// Generate JWTs
	tokenPair, err := s.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, "", err
	}

	// Store refresh session
	if err := s.createRefreshSession(ctx, s.db, user.ID, tokenPair); err != nil {
		return nil, "", err
	}

	if err := s.repo.UpdateLastLoginAt(ctx, s.db, user.ID); err != nil {
		return nil, "", err
	}

	return s.buildAuthResponse(user, tokenPair), tokenPair.RefreshToken, nil
}

func (s *AuthService) GoogleLogin(ctx context.Context) (string, error) {

	state, err := auth.GenerateCodeState()
	if err != nil {
		return "", err
	}

	verifier, err := auth.GenerateCodeVerifier()
	if err != nil {
		return "", err
	}

	challenge := auth.GenerateCodeChallenge(verifier)

	if err := s.redis.Set(ctx, state, verifier, auth.PKCE_TTL).Err(); err != nil {
		return "", err
	}

	return s.google.AuthURL(
		state,
		challenge,
	), nil
}

func (s *AuthService) GoogleCallback(ctx context.Context, code string, state string) (*auth.AuthResponse, string, error) {

	verifier, err := s.redis.Get(ctx, state).Result()
	if err == redis.Nil {
		return nil, "", auth.ErrInvalidState
	}
	if err != nil {
		return nil, "", err
	}

	// One-time use
	_ = s.redis.Del(ctx, state).Err()

	// Exchange authorization code
	token, err := s.google.ExchangeCode(ctx, code, verifier)
	if err != nil {
		return nil, "", err
	}

	// Verify Google ID Token
	googleUser, err := s.google.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, "", err
	}

	var user *User

	// Existing Google account?
	authAccount, err := s.repo.FindAuthAccount(ctx, s.db, "google", googleUser.Subject)
	if err != nil {
		return nil, "", err
	}

	if authAccount != nil {
		user, err = s.repo.FindUserByID(ctx, s.db, authAccount.UserID)
		if err != nil {
			return nil, "", err
		}
	} else {

		// Create or Link Account

		tx, err := s.db.Begin(ctx)
		if err != nil {
			return nil, "", err
		}
		defer tx.Rollback(ctx)

		user, err = s.repo.FindUserByEmail(ctx, tx, googleUser.Email)
		if err != nil {
			return nil, "", err
		}

		if user == nil {

			var avatar *string
			if googleUser.AvatarURL != "" {
				avatar = &googleUser.AvatarURL
			}

			user = &User{
				Email:           googleUser.Email,
				FullName:        googleUser.FullName,
				AvatarURL:       avatar,
				IsEmailVerified: googleUser.EmailVerified,
			}

			err = s.createUser(ctx, tx, user)
			if err != nil {
				return nil, "", err
			}
		}

		err = s.repo.CreateAuthAccount(
			ctx,
			tx,
			&AuthAccount{
				UserID:         user.ID,
				Provider:       "google",
				ProviderUserID: googleUser.Subject,
			},
		)
		if err != nil {
			return nil, "", err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, "", err
		}
	}

	// Update last login

	if err := s.repo.UpdateLastLoginAt(ctx, s.db, user.ID); err != nil {
		return nil, "", err
	}

	// Generate Tokens

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, "", err
	}

	if err := s.createRefreshSession(ctx, s.db, user.ID, tokenPair); err != nil {
		return nil, "", err
	}

	return s.buildAuthResponse(user, tokenPair), tokenPair.RefreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*auth.AuthResponse, string, error) {

	refreshTokenHash := s.jwt.HashRefreshToken(refreshToken)

	session, err := s.repo.FindRefreshTokenByHash(ctx, s.db, refreshTokenHash)
	if err != nil {
		return nil, "", err
	}

	if session == nil {
		return nil, "", auth.ErrRefreshTokenNotFound
	}

	if session.Revoked {
		return nil, "", auth.ErrRefreshTokenRevoked
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, "", auth.ErrRefreshTokenExpired
	}

	//------------------------------------------------------------------
	// Load user
	//------------------------------------------------------------------

	user, err := s.repo.FindUserByID(ctx, s.db, session.UserID)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", auth.ErrUserNotFound
	}

	//------------------------------------------------------------------
	// Rotate refresh token
	//------------------------------------------------------------------

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, "", err
	}

	if err := s.createRefreshSession(ctx, tx, user.ID, tokenPair); err != nil {
		return nil, "", err
	}

	if err := s.repo.RevokeRefreshToken(
		ctx,
		tx,
		session.ID,
	); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return s.buildAuthResponse(user, tokenPair), tokenPair.RefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {

	refreshTokenHash := s.jwt.HashRefreshToken(refreshToken)

	session, err := s.repo.FindRefreshTokenByHash(ctx, s.db, refreshTokenHash)
	if err != nil {
		return err
	}

	if session == nil {
		return auth.ErrRefreshTokenInvalid
	}

	// Already logged out.
	if session.Revoked {
		return nil
	}

	return s.repo.RevokeRefreshToken(ctx, s.db, session.ID)
}

func (s *AuthService) Profile(ctx context.Context,userID uuid.UUID) (*auth.ProfileResponse, error) {

	user, err := s.repo.FindUserByID(ctx,s.db,userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, auth.ErrUserNotFound
	}

	profile := s.buildProfileResponse(user)
	
	return &profile, nil
}