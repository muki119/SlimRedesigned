package dtos

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1,max=15"`
	Password string `json:"password" validate:"required,min=8,max=96"`
}

type RegisterRequest struct {
	Forename    string `json:"forename" db:"forename" validate:"required,min=1,max=30"`
	Surname     string `json:"surname"  db:"surname"  validate:"required,min=1,max=30"`
	Username    string `json:"username"   db:"username"  validate:"required,min=1,max=15"`
	Email       string `json:"email" db:"email"  validate:"required,email,max=320"`
	Password    string `json:"password"  db:"password"   validate:"required,min=8,max=96"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"  validate:"required"`
}
