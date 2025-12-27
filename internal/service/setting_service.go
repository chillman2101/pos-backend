package service

import (
	"encoding/json"
	"pos-backend/internal/domain"
)

type SettingService struct {
	repo domain.SettingRepository
}

func NewSettingService(repo domain.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// GetSettings returns all settings grouped by category
func (s *SettingService) GetSettings() (map[string]map[string]interface{}, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Group by category
	result := make(map[string]map[string]interface{})

	for _, setting := range settings {
		if result[setting.Category] == nil {
			result[setting.Category] = make(map[string]interface{})
		}

		// Try to parse JSON value, otherwise use as string
		var value interface{}
		if err := json.Unmarshal([]byte(setting.Value), &value); err != nil {
			value = setting.Value
		}

		result[setting.Category][setting.Key] = value
	}

	return result, nil
}

// GetSettingByKey returns a single setting by key
func (s *SettingService) GetSettingByKey(key string) (*domain.Setting, error) {
	return s.repo.GetByKey(key)
}

// UpdateSettings updates multiple settings
func (s *SettingService) UpdateSettings(settingsMap map[string]map[string]interface{}) error {
	var settings []domain.Setting

	for category, values := range settingsMap {
		for key, value := range values {
			// Convert value to JSON string
			valueBytes, err := json.Marshal(value)
			if err != nil {
				return err
			}

			settings = append(settings, domain.Setting{
				Key:      key,
				Value:    string(valueBytes),
				Category: category,
			})
		}
	}

	return s.repo.BulkUpsert(settings)
}

// InitializeDefaultSettings creates default settings if they don't exist
func (s *SettingService) InitializeDefaultSettings() error {
	// Check if settings already exist
	existing, _ := s.repo.GetAll()
	if len(existing) > 0 {
		return nil // Settings already initialized
	}

	defaultSettings := []domain.Setting{
		// Store settings
		{Key: "store_name", Value: `"My POS Store"`, Category: "store"},
		{Key: "store_address", Value: `"Jl. Example No. 123, Jakarta"`, Category: "store"},
		{Key: "store_phone", Value: `"021-12345678"`, Category: "store"},
		{Key: "store_email", Value: `"store@example.com"`, Category: "store"},

		// Tax settings
		{Key: "tax_enabled", Value: `true`, Category: "tax"},
		{Key: "tax_rate", Value: `10`, Category: "tax"},
		{Key: "tax_label", Value: `"PPN"`, Category: "tax"},
		{Key: "currency", Value: `"IDR"`, Category: "tax"},
		{Key: "currency_symbol", Value: `"Rp"`, Category: "tax"},

		// Receipt settings
		{Key: "show_logo", Value: `true`, Category: "receipt"},
		{Key: "show_address", Value: `true`, Category: "receipt"},
		{Key: "show_phone", Value: `true`, Category: "receipt"},
		{Key: "footer_text", Value: `"Terima kasih atas kunjungan Anda!"`, Category: "receipt"},

		// System settings
		{Key: "auto_sync_enabled", Value: `true`, Category: "system"},
		{Key: "sync_interval", Value: `5`, Category: "system"},
		{Key: "offline_mode_enabled", Value: `true`, Category: "system"},
		{Key: "theme", Value: `"light"`, Category: "system"},
	}

	return s.repo.BulkUpsert(defaultSettings)
}
