package db

import (
	"errors"

	"gorm.io/gorm"
)

// ListNodes retrieves all nodes from the database with applied filters and pagination
func (db *Database) ListNodes(filter NodeFilter, limit Limit) (nodes []Node, err error) {
	query := db.gormDB.Model(&Node{})

	if filter.NodeID != nil {
		query = query.Where("node_id= ?", *filter.NodeID)
	}
	if filter.FarmID != nil {
		query = query.Where("farm_id = ?", *filter.FarmID)
	}
	if filter.TwinID != nil {
		query = query.Where("twin_id = ?", *filter.TwinID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	offset := (limit.Page - 1) * limit.Size
	query = query.Offset(int(offset)).Limit(int(limit.Size))

	if result := query.Find(&nodes); result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

// GetNode retrieves a specific node by its nodeID
func (db *Database) GetNode(nodeID uint64) (node Node, err error) {
	if result := db.gormDB.First(&node, nodeID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return node, ErrRecordNotFound
		}
		return node, result.Error
	}
	return node, nil
}

// RegisterNode registers a new node in the database
func (db *Database) RegisterNode(node Node) (err error) {
	if result := db.gormDB.First(&node, nodeID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return result.Error
	}

	if result := db.gormDB.Create(&node); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrRecordAlreadyExists
		}
		return result.Error
	}
	return nil
}

// Uptime updates the uptime for a specific node
func (db *Database) Uptime(nodeID uint64, report Uptime) (err error) {
	var node Node
	if result := db.gormDB.First(&node, nodeID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return result.Error
	}

	node.Uptime = report
	if result := db.gormDB.Save(&node); result.Error != nil {
		return result.Error
	}
	return nil
}

// Consumption updates the consumption report for a specific node
func (db *Database) Consumption(nodeID uint64, report Consumption) (err error) {
	var node Node

	if err := db.gormDB.First(&node, nodeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}

	node.Consumption = report
	if result := db.gormDB.Save(&node); result.Error != nil {
		return result.Error
	}
	return nil
}
