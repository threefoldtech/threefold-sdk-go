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

	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"list", "farms"})

	result := query.Find(&farms)
	stop()
	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"list", "farms"})
		return nil, result.Error
	}

	return farms, nil
}

func (db *Database) GetFarm(farmID uint64) (farm Farm, err error) {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select_by_farmID", "farms"})

	result := db.gormDB.First(&farm, farmID)
	stop()
	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select_by_farmID", "farms"})

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return farm, ErrRecordNotFound
		}
		return farm, result.Error
	}

	return
}

func (db *Database) CreateFarm(farm Farm) (uint64, error) {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"create", "farms"})
	err := db.gormDB.Create(&farm).Error

	stop()
	if err != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"create", "farms"})
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, ErrRecordAlreadyExists
		}
		return 0, err
	}
	return farm.FarmID, nil
}

func (db *Database) UpdateFarm(farmID uint64, name string) (err error) {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"update", "farms"})

	result := db.gormDB.Model(&Farm{}).Where("farm_id = ?", farmID).Updates(map[string]interface{}{
		"farm_name": name,
	})
	stop()
	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"update", "farms"})
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
