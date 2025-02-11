package db

import (
	"errors"

	"github.com/lib/pq"
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

// UpdateAccount updates an account's relays and RMB encryption key
func (db *Database) UpdateAccount(twinID uint64, relays pq.StringArray, rmbEncKey string) error {
	result := db.gormDB.Model(&Account{}).Where("twin_id = ?", twinID).Updates(map[string]interface{}{
		"relays":      relays,
		"rmb_enc_key": rmbEncKey,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// GetAccountByTwinID retrieves an account by its twin ID
func (db *Database) GetAccount(twinID uint64) (Account, error) {
	var account Account
	if err := db.gormDB.First(&account, twinID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Account{}, ErrRecordNotFound
		}
		return Account{}, err
	}
	return account, nil
}

func (db *Database) GetAccountByPublicKey(publicKey string) (Account, error) {
	var account Account
	result := db.gormDB.Where("public_key = ?", publicKey).First(&account)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, ErrRecordNotFound
		}
		return Account{}, result.Error
	}
	return account, nil
}
