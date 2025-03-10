package postgres

// package postgres

// import (
// 	"database/sql"

// 	"example.com/rest/internal/domain"
// 	"github.com/jmoiron/sqlx"
// )

// type userRepo struct {
// 	db *sqlx.DB
// }

// type UserRepo interface {
// 	Insert(user *domain.UserCredentials) (int, error)
// 	GetByID(id int) (*domain.User, error)
// 	GetByEmail(email string) (*domain.User, error)
// 	Update(user *domain.User) (*domain.User, error)
// 	Delete(id int) error
// }

// func NewUserRepo(db *sql.DB) UserRepo {
// 	return &userRepo{db: db}
// }

// func (r *userRepo) Insert(email, passwordHash string) (int, error) {
// 	row := r.db.QueryRowx("INSERT INTO users (email, password_hash) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id", email, passwordHash)
// 	var id int
// 	err := row.Scan(&id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return 0, domain.Errorf(domain.CONFLICT_ERROR, "email already exists")
// 		}
// 		return 0, domain.Errorf(domain.INTERNAL_ERROR, "failed to insert user: %v", err)
// 	}
// 	return id, nil
// }
