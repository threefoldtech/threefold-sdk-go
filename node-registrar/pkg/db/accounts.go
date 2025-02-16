package db

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// CreateAccount creates a new account in the database
func (db *Database) CreateAccount(account *Account) error {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"create", "accounts"})
	err := db.gormDB.Create(account).Error
	stop()
	if err != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"create", "accounts"})
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrRecordAlreadyExists
	}
	return err
}

// UpdateAccount updates an account's relays and RMB encryption key
func (db *Database) UpdateAccount(twinID uint64, relays pq.StringArray, rmbEncKey string) error {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"update", "accounts"})

	result := db.gormDB.Model(&Account{}).Where("twin_id = ?", twinID).Updates(map[string]interface{}{
		"relays":      relays,
		"rmb_enc_key": rmbEncKey,
	})
	stop()
	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"update", "accounts"})
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
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select_by_twinID", "accounts"})
	err := db.gormDB.First(&account, twinID).Error
	stop()
	if err != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select_by_twinID", "accounts"})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Account{}, ErrRecordNotFound
		}
		return Account{}, err
	}
	return account, nil
}

func (db *Database) GetAccountByPublicKey(publicKey string) (Account, error) {
	var account Account
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select_by_public_key", "accounts"})

	result := db.gormDB.Where("public_key = ?", publicKey).First(&account)
	stop()
	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select_by_public_key", "accounts"})
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, ErrRecordNotFound
		}
		return Account{}, result.Error
	}
	return account, nil
}
