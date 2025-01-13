package db

func (db *DataBase) ListFarms(filter FarmFilter, limit Limit) (farms []Farm, err error) {
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

func (db *DataBase) GetFarm(farmID uint64) (farm Farm, err error) {
	if result := db.gormDB.First(&farm, farmID); result.Error != nil {
		return farm, result.Error
	}

	return
}

func (db *DataBase) CreateFarm(farm Farm) (err error) {
	return db.gormDB.Create(&farm).Error
}

func (db *DataBase) UpdateFarm(farm Farm) (err error) {
	result := db.gormDB.Save(&farm)
	return result.Error
}
