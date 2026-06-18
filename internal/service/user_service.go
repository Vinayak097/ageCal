package service

import (
	"context"
	"time"

	db "github.com/user/dob-api/db/sqlc"
	"github.com/user/dob-api/internal/models"
	"github.com/user/dob-api/internal/repository"
)

const dobLayout = "2006-01-02"

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (*models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int32) (*models.ListUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dobTime, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  dobTime,
	})
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserWithAge(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dobTime, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Update(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  dobTime,
	})
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, limit int32) (*models.ListUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]models.UserWithAgeResponse, len(users))
	for i, u := range users {
		data[i] = *toUserWithAge(u)
	}

	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &models.ListUsersResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// CalculateAge returns age in years for dob relative to now.
// Exported for unit testing.
func CalculateAge(dob time.Time, now time.Time) int {
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}

func toUserResponse(u db.User) *models.UserResponse {
	return &models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dobLayout),
	}
}

func toUserWithAge(u db.User) *models.UserWithAgeResponse {
	return &models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dobLayout),
		Age:  CalculateAge(u.Dob, time.Now()),
	}
}
