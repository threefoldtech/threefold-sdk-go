package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (db *Database) ListFarms(filter FarmFilter, limit Limit) (farms []Farm, err error) {
	query := db.gormDB.Model(&Farm{})

	if filter.FarmName != nil {
		query = query.Where("farm_name ILIKE ?", "%"+*filter.FarmName+"%") // Case-insensitive search
	}
	if filter.FarmID != nil {
		query = query.Where("farm_id = ?", *filter.FarmID)
	}
	if filter.TwinID != nil {
		query = query.Where("twin_id = ?", *filter.TwinID)
	}

	offset := (limit.Page - 1) * limit.Size
	query = query.Offset(int(offset)).Limit(int(limit.Size))

	if result := query.Find(&farms); result.Error != nil {
		return nil, result.Error
	}

	return farms, nil
}

func (db *Database) GetFarm(farmID uint64) (farm Farm, err error) {
	if result := db.gormDB.First(&farm, farmID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return farm, ErrRecordNotFound
		}
		return farm, result.Error
	}

	return
}

func (db *Database) CreateFarm(farm Farm) (err error) {
	return db.gormDB.Create(&farm).Error
}

func (db *Database) UpdateFarm(farmID uint64, farm Farm) (err error) {
	if result := db.gormDB.First(&farm, farmID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return result.Error
	}

	result := db.gormDB.Save(&farm)
	return result.Error
}
