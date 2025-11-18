package repository

import (
	"fmt"
	"log/slog"
)

var UserRepo UsersRepo

func Init() {
	usrRepo, err := NewUserRepo()
	slog.Info("init user repository")
	roflix := User{
		Login:        "roflixd",
		PasswordHash: "$2a$10$mhrsuffUN17whEp/wi3hpeSceEa1F/HJvspVnX.Py9MU94WQ5KzRq",
	}

	kolfiild := User{
		Login:        "kolfiild",
		PasswordHash: "$2a$10$mhrsuffUN17whEp/wi3hpeSceEa1F/HJvspVnX.Py9MU94WQ5KzRq",
	}

	if err != nil {
		panic(err)
	}

	usrRepo.Add(roflix)
	usrRepo.Add(kolfiild)
	UserRepo = usrRepo
}

type User struct {
	Login        string
	PasswordHash string
}

type UsersRepo interface {
	IsExists(string) (bool, error)
	Get(string) (User, error)
	Add(User) error
}

type InMemoryUsersRepo struct {
	users map[string]User
}

func NewUserRepo() (UsersRepo, error) {
	return &InMemoryUsersRepo{
		users: make(map[string]User),
	}, nil
}

func (repo *InMemoryUsersRepo) IsExists(login string) (bool, error) {
	_, result := repo.users[login]
	return result, nil
}

func (repo *InMemoryUsersRepo) Get(login string) (User, error) {
	user, exists := repo.users[login]
	if !exists {
		return User{}, fmt.Errorf("user %s doesn't exitst", login)
	}
	return user, nil
}

func (repo *InMemoryUsersRepo) Add(new User) error {
	_, exists := repo.users[new.Login]
	if exists {
		return fmt.Errorf("user %s already exitst", new.Login)
	}
	repo.users[new.Login] = new
	return nil
}
