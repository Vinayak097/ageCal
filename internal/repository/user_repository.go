package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	db "github.com/user/dob-api/db/sqlc"
)

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, params db.ListUsersParams) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, params)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, params db.ListUsersParams) ([]db.User, error) {
	return r.q.ListUsers(ctx, params)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
