package db

import (
	"errors"
	"time"

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
func (db *Database) RegisterNode(node Node) (uint64, error) {
	if result := db.gormDB.Create(&node); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, ErrRecordAlreadyExists
		}
		return 0, result.Error
	}
	return node.NodeID, nil
}

// Uptime updates the uptime for a specific node
func (db *Database) GetUptimeReports(nodeID uint64, start, end time.Time) ([]UptimeReport, error) {
	var reports []UptimeReport
	result := db.gormDB.Where("node_id = ? AND timestamp BETWEEN ? AND ?",
		nodeID, start, end).Order("timestamp asc").Find(&reports)
	return reports, result.Error
}

func (db *Database) CreateUptimeReport(report *UptimeReport) error {
	return db.gormDB.Create(report).Error
}
