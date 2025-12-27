package repository

import (
	"pos-backend/internal/domain"

	"gorm.io/gorm"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domain.SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetByKey(key string) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) GetByCategory(category string) ([]domain.Setting, error) {
	var settings []domain.Setting
	err := r.db.Where("category = ?", category).Find(&settings).Error
	return settings, err
}

func (r *settingRepository) GetAll() ([]domain.Setting, error) {
	var settings []domain.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingRepository) Upsert(setting *domain.Setting) error {
	var existing domain.Setting
	err := r.db.Where("key = ?", setting.Key).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		return r.db.Create(setting).Error
	} else if err != nil {
		return err
	}

	// Update existing
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"value":    setting.Value,
		"category": setting.Category,
	}).Error
}

func (r *settingRepository) BulkUpsert(settings []domain.Setting) error {
	for _, setting := range settings {
		if err := r.Upsert(&setting); err != nil {
			return err
		}
	}
	return nil
}

func (r *settingRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&domain.Setting{}).Error
}
