package repos

import (
	"clean-rest-arch/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.UserEntity) (uint, error)
	GetUserById(id uint) (*models.UserEntity, error)
	GetUserByUsername(username string) (*models.UserEntity, error)
}

type userRepo struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{database: db}
}

func (r *userRepo) CreateUser(user *models.UserEntity) (uint, error) {
	const op = "storage.repos.CreateUser"

	result := r.database.Create(user)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user.Id, nil
}

func (r *userRepo) GetUserById(id uint) (*models.UserEntity, error) {
	const op = "storage.repos.GetUserById"

	var userFromDb models.UserEntity
	result := r.database.Where("id = ?", id).First(&userFromDb)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &userFromDb, nil
}

func (r *userRepo) GetUserByUsername(username string) (*models.UserEntity, error) {
	const op = "storage.repos.GetUserByUsername"

	var userFromDb models.UserEntity
	result := r.database.Where("username = ?", username).First(&userFromDb)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &userFromDb, nil
}
