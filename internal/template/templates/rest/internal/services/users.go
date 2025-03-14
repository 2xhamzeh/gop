package services

import (
	"context"
	"errors"

	"example.com/rest/internal/domain"
	"example.com/rest/internal/postgres"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo UserRepo
}

type UserRepo interface {
	Insert(ctx context.Context, email, passwordHash string) (int, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int) error
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, req *domain.UserCredentials) (int, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return 0, err
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, domain.Errorf(domain.INTERNAL_ERROR, "failed to hash password").Wrap(err)
	}

	// insert user
	id, err := s.userRepo.Insert(ctx, req.Email, string(passwordHash))
	if err != nil {
		if errors.Is(err, postgres.ErrConflict) {
			return 0, domain.Errorf(domain.CONFLICT_ERROR, "user already exists")
		}
		return 0, err
	}
	return id, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, req *domain.UserCredentials) (*domain.User, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// find user
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
		}
		return nil, err
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id int, req *domain.UserPatch) (*domain.User, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// get user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return nil, err
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid password")
	}

	// check which fields to update
	if req.Email != nil {
		// compare old and new email
		if *req.Email == user.Email {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "email already in use")
		}
		// update email
		user.Email = *req.Email
	}

	if req.NewPassword != nil {
		// compare old and new password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*req.NewPassword))
		if err == nil {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "new password must be different")
		}
		// hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to hash password").Wrap(err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	// update user
	user, err = s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "update conflict")
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return err
	}

	return nil
}
