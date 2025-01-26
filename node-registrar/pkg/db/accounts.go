package db

import (
	"errors"

	"gorm.io/gorm"
)

// CreateAccount creates a new account in the database
func (db *Database) CreateAccount(account *Account) error {
	err := db.gormDB.Create(account).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrRecordAlreadyExists
	}
	return err
}

// GetAccountByTwinID retrieves an account by its twin ID
func (db *Database) GetAccount(twinID uint64) (*Account, error) {
	var account Account
	if err := db.gormDB.First(&account, twinID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &account, nil
}
