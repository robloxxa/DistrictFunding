package auth

// TODO: add other fields like phone number, firstname, secondname, etc.
type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,lte=32"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Password  string `json:"password" validate:"required,gt=6"`
}

type SignInRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password string `json:"password" validate:"required"`
}

type MeRequest struct {
	Id string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
