package auth

// TODO: add other fields like phone number, firstname, secondname, etc.
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,lte=32"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Password  string `json:"password" validate:"required,gt=6"`
}
