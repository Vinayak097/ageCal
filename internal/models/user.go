package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

type ListUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int32                 `json:"page"`
	Limit      int32                 `json:"limit"`
	TotalPages int32                 `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
