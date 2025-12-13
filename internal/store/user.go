package store

// User represents a user entity
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// UserStore defines the interface for user data operations
type UserStore interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user User) (*User, error)
	Update(id int, user User) (*User, error)
	Delete(id int) error
}
