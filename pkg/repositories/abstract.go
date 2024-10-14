package repositories

import (
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories/implrepo"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=abstract.go -destination=mocks/mock.go

type ProjectRepository interface {
	Create(p *entity.Project) (int64, error)
	GetById(id int64, userId int64) (entity.Project, error)
	GetAll(userId int64) ([]entity.Project, error)
	UpdateById(id int64, input dto.UpdateProjectDTO, userId int64) error
	DeleteById(id int64, userId int64) error
}

type AuthRepository interface {
	SignUp(u *entity.User) (int64, error)
	SignIn(username, passwordHash string) (int64, error)
	CreateRefreshToken(userId int64, token string, expiresAt time.Time) error
	FindRefreshToken(token string) (int64, time.Time, error)
	DeleteRefreshToken(token string) error
}

type AbstractRepository struct {
	ProjectRepository
	AuthRepository
}

func NewRepository(db *sqlx.DB) *AbstractRepository {
	return &AbstractRepository{
		ProjectRepository: implrepo.NewProjectRepository(db),
		AuthRepository:    implrepo.NewUserRepository(db),
	}
}
