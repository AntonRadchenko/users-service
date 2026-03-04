package user

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserParams struct {
	Email    string
	Password string
}

type UpdateUserParams struct {
	Email    *string
	Password *string
}

type User struct {
	ID    uint
	Email string
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(r UserRepoInterface) *UserService {
	return &UserService{repo: r}
}

func hashPass(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *UserService) CreateUser(params CreateUserParams) (*User, error) {
	if strings.TrimSpace(params.Email) == "" {
		return nil, errors.New("email is empty")
	}

	if strings.TrimSpace(params.Password) == "" {
		return nil, errors.New("password is empty")
	}

	hashedPassword, err := hashPass(params.Password)
	if err != nil {
		return nil, errors.New("fail to hash password")
	}

	dbUser := &UserStruct{
		Email:    params.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.repo.Create(dbUser)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    createdUser.ID,
		Email: createdUser.Email,
	}, nil
}

func (s *UserService) GetUser(id uint) (*User, error) {
	dbUser, err := s.repo.GetByID(id)
	if err != nil || dbUser.ID == 0 {
		return nil, err
	}

	return &User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	}, nil
}

func (s *UserService) GetUsers() ([]User, error) {
	dbUsers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		})
	}
	return users, nil
}

func (s *UserService) UpdateUser(id uint, params UpdateUserParams) (*User, error) {
	dbUser, err := s.repo.GetByID(id)
	if err != nil || dbUser.ID == 0 {
		return nil, errors.New("user not found")
	}

	updated := false

	if params.Email != nil {
		if strings.TrimSpace(*params.Email) == "" {
			return nil, errors.New("email is empty")
		}
		dbUser.Email = *params.Email
		updated = true
	}

	if params.Password != nil {
		if strings.TrimSpace(*params.Password) == "" {
			return nil, errors.New("password is empty")
		}
		hashedPassword, err := hashPass(*params.Password)
		if err != nil {
			return nil, errors.New("fail to hash password")
		}
		dbUser.Password = hashedPassword
		updated = true
	}

	if !updated {
		return nil, errors.New("no fields to update")
	}

	updatedUser, err := s.repo.Update(&dbUser)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil || user.ID == 0 {
		return errors.New("user not found")
	}

	err = s.repo.Delete(&user)
	if err != nil {
		return err
	}
	return nil
}
