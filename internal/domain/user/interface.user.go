package user

import (
	"context"
)

type Context interface {
	BindJSON(obj any) error
	BindURI(obj any) error
	BindQuery(obj any) error
	BindHeader(obj any) error
	Bind(obj any) error
	JSON(code int, v any) error
	Param(name string) string
	GetUserID() string
	Get(key string) any
	GetContext() context.Context
}

type UserHandler interface {
	Create(c Context) error
	Get(c Context) error
	List(c Context) error
	Update(c Context) error
	Delete(c Context) error
}

type UserService interface {
	Create(ctx context.Context, req *CreateUserRequest, createdBy string) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, req *UpdateUserRequest, updatedBy string) (*User, error)
	Delete(ctx context.Context, id, deletedBy string) error
}

type UserRepository interface {
	StartContext(ctx context.Context) context.Context
	DeferErrorContext(ctx context.Context, err error)

	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, u *User) error
	SoftDelete(ctx context.Context, id, deletedBy string) error
}
