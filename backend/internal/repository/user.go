package repository

import (
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetAllUsers(page, limit int) ([]*model.User, int64, error)
	GetByEmailOrUsername(login string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByVerificationToken(token string) (*model.User, error)
	Update(user *model.User) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetAllUsers(page, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// Count total users
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Fetch paginated users
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *GormUserRepository) GetByEmailOrUsername(login string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ? OR username = ?", login, login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByVerificationToken(token string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormUserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
