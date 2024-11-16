package service

import (
	"context"

	"ticket-booking-app-backend/internal/application/types/requests"
	"ticket-booking-app-backend/internal/application/types/responses"
	"ticket-booking-app-backend/internal/domain/entities"
	"ticket-booking-app-backend/internal/domain/repository"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/pkg/values"

	"github.com/sirupsen/logrus"
)

const (
	BEARER_TOKEN_TYPE = "Bearer"
)

type Users interface {
	UserSignIn(ctx context.Context, input *requests.UserSignInRequest) (*responses.TokenResponse, error)
	UserSignUp(ctx context.Context, input *requests.UserSignUpRequest) error
	AdminSignIn(ctx context.Context, input *requests.AdminSignInRequest) (*responses.TokenResponse, error)
	AdminSignUp(ctx context.Context, input *requests.AdminSignUpRequest) error
	OrganizerSignIn(ctx context.Context, input *requests.OrganizerSignInRequest) (*responses.TokenResponse, error)
	OrganizerSignUp(ctx context.Context, input *requests.OrganizerSignUpRequest) error
}

type usersService struct {
	repo       repository.UsersRepository
	commonRepo repository.CommonRepository
	jwt        helpers.Jwt
}

func NewUsersService(repo repository.UsersRepository, commonRepo repository.CommonRepository, jwt helpers.Jwt) *usersService {
	return &usersService{
		repo:       repo,
		commonRepo: commonRepo,
		jwt:        jwt,
	}
}

func (s *usersService) UserSignUp(ctx context.Context, input *requests.UserSignUpRequest) error {
	if err := s.commonRepo.CheckIfUserExistsByEmail(ctx, input.Email); err == nil {
		return domainErrors.ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		logrus.Errorf("Error hashing password: %s", err)
		return nil
	}

	// Create the user
	user := entities.User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	}

	// Save the user to the repository
	err = s.repo.Create(ctx, values.UserRole, &user)
	if err != nil {
		logrus.Errorf("Error creating user: %s", err)
		return err
	}

	return nil
}
func (s *usersService) UserSignIn(ctx context.Context, input *requests.UserSignInRequest) (*responses.TokenResponse, error) {
	if err := s.commonRepo.CheckIfUserExistsByEmail(ctx, input.Email); err != nil {
		return nil, err
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if !helpers.CheckPasswordHash(input.Password, user.Password) {
		return nil, domainErrors.ErrUserPasswordIncorrect
	}

	token, err := s.jwt.CreateAccessToken(helpers.UserAccessTokenClaims{
		UserId: user.ID,
		Role:   values.UserRole,
	})
	if err != nil {
		logrus.Errorf("Error creating access token: %s", err)
		return nil, err
	}
	return &responses.TokenResponse{
		Success:               true,
		Token:                 token.AccessToken,
		TokenType:             BEARER_TOKEN_TYPE,
		ExpiresAt:             token.AccessTokenExpiresAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenType:      BEARER_TOKEN_TYPE,
		RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
	}, nil
}

// AdminSignIn handles the sign-in process for admin users.
func (s *usersService) AdminSignIn(ctx context.Context, input *requests.AdminSignInRequest) (*responses.TokenResponse, error) {
	// Check if admin user exists
	admin, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domainErrors.ErrAdminNotFound
	}

	// Verify the user's role is admin
	if admin.Role != values.AdminRole {
		return nil, domainErrors.ErrInsufficientPermissions
	}

	// Verify the password
	if !helpers.CheckPasswordHash(input.Password, admin.Password) {
		return nil, domainErrors.ErrUserPasswordIncorrect
	}

	// Generate access token
	token, err := s.jwt.CreateAccessToken(helpers.UserAccessTokenClaims{
		UserId: admin.ID,
		Role:   values.AdminRole,
	})
	if err != nil {
		logrus.Errorf("Error creating access token: %s", err)
		return nil, err
	}

	return &responses.TokenResponse{
		Success:               true,
		Token:                 token.AccessToken,
		TokenType:             BEARER_TOKEN_TYPE,
		ExpiresAt:             token.AccessTokenExpiresAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenType:      BEARER_TOKEN_TYPE,
		RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
	}, nil
}

// AdminSignUp handles the sign-up process for admin users.
func (s *usersService) AdminSignUp(ctx context.Context, input *requests.AdminSignUpRequest) error {
	// Check if an admin with the same email already exists
	if err := s.commonRepo.CheckIfUserExistsByEmail(ctx, input.Email); err == nil {
		logrus.Warn("User already exists")
		return nil
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		logrus.Errorf("Error hashing password: %s", err)
		return err
	}

	// Create the admin user
	admin := entities.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Save the admin to the repository
	err = s.repo.Create(ctx, values.AdminRole, &admin)
	if err != nil {
		logrus.Errorf("Error creating admin: %s", err)
		return err
	}

	return nil
}

// OrganizerSignIn handles the sign-in process for organizer users.
func (s *usersService) OrganizerSignIn(ctx context.Context, input *requests.OrganizerSignInRequest) (*responses.TokenResponse, error) {
	// Check if organizer exists
	organizer, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domainErrors.ErrOrganizerNotFound
	}

	// Verify the user's role is organizer
	if organizer.Role != values.OrganizerRole {
		return nil, domainErrors.ErrInsufficientPermissions
	}

	// Verify the password
	if !helpers.CheckPasswordHash(input.Password, organizer.Password) {
		return nil, domainErrors.ErrUserPasswordIncorrect
	}

	// Generate access token
	token, err := s.jwt.CreateAccessToken(helpers.UserAccessTokenClaims{
		UserId: organizer.ID,
		Role:   values.OrganizerRole,
	})
	if err != nil {
		logrus.Errorf("Error creating access token: %s", err)
		return nil, err
	}

	return &responses.TokenResponse{
		Success:               true,
		Token:                 token.AccessToken,
		TokenType:             BEARER_TOKEN_TYPE,
		ExpiresAt:             token.AccessTokenExpiresAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenType:      BEARER_TOKEN_TYPE,
		RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
	}, nil
}

// OrganizerSignUp handles the sign-up process for organizer users.
func (s *usersService) OrganizerSignUp(ctx context.Context, input *requests.OrganizerSignUpRequest) error {
	// Check if an organizer with the same email already exists
	if err := s.commonRepo.CheckIfUserExistsByEmail(ctx, input.Email); err == nil {
		return domainErrors.ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		logrus.Errorf("Error hashing password: %s", err)
		return err
	}

	// Create the organizer user
	organizer := entities.User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
		Address:  input.Address,
		Phone:    input.Phone,
	}

	// Save the organizer to the repository
	err = s.repo.Create(ctx, values.OrganizerRole, &organizer)
	if err != nil {
		logrus.Errorf("Error creating organizer: %s", err)
		return err
	}

	return nil
}
