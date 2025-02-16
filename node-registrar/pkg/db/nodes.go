package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const ZOS4VersionKey = "zos_4"

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

	start := time.Now()

	result := query.Find(&nodes)
	db.metrics.DBOperationsDuration.WithLabelValues("list", "nodes").Observe(time.Since(start).Seconds())
	if result.Error != nil {
		db.metrics.DBOperationsErrors.WithLabelValues("list", "nodes").Inc()
		return nil, result.Error
	}

	return nodes, nil
}

// GetNode retrieves a specific node by its nodeID
func (db *Database) GetNode(nodeID uint64) (node Node, err error) {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select_by_nodeID", "nodes"})
	result := db.gormDB.First(&node, nodeID)

	stop()

	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select_by_nodeID", "nodes"})

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return node, ErrRecordNotFound
		}
		return node, result.Error
	}
	return node, nil
}

// RegisterNode registers a new node in the database
func (db *Database) RegisterNode(node Node) (uint64, error) {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"create", "nodes"})
	result := db.gormDB.Create(&node)

	stop()

	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"create", "nodes"})
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, ErrRecordAlreadyExists
		}
		return 0, result.Error
	}

	return node.NodeID, nil
}

func (db *Database) UpdateNode(nodeID uint64, node Node) error {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"update", "nodes"})

	result := db.gormDB.Model(&Node{}).
		Where("node_id = ?", nodeID).
		Updates(node)
	stop()

	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"update", "nodes"})
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Uptime updates the uptime for a specific node
func (db *Database) GetUptimeReports(nodeID uint64, start, end time.Time) ([]UptimeReport, error) {
	var reports []UptimeReport
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select_by_nodeID", "uptime_reports"})
	result := db.gormDB.Where("node_id = ? AND timestamp BETWEEN ? AND ?",
		nodeID, start, end).Order("timestamp asc").Find(&reports)

	stop()

	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select_by_nodeID", "uptime_reports"})
	}
	return reports, result.Error
}

func (db *Database) CreateUptimeReport(report *UptimeReport) error {
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"create", "uptime_reports"})
	err := db.gormDB.Create(report).Error
	stop()
	if err != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"create", "uptime_reports"})
	}
	return db.gormDB.Create(report).Error
}

func (db *Database) SetZOSVersion(version string) error {
	var current ZosVersion
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"create", "zos_version"})

	result := db.gormDB.Where(ZosVersion{Key: ZOS4VersionKey}).Attrs(ZosVersion{Version: version}).FirstOrCreate(&current)
	stop()

	if result.Error != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"create", "zos_version"})
		return result.Error
	}

	if result.RowsAffected == 0 {
		if current.Version == version {
			return errors.New("version already set")
		}
		return db.gormDB.Model(&current).
			Select("version").
			Update("version", version).Error
	}
	return nil
}

func (db *Database) GetZOSVersion() (string, error) {
	var setting ZosVersion
	stop := db.metrics.RecordDuration(db.metrics.DBOperationsDuration, []string{"select", "zos_version"})

	err := db.gormDB.Where("key = ?", "zos_4").First(&setting).Error
	stop()

	if err != nil {
		db.metrics.RecordCount(db.metrics.DBOperationsErrors, []string{"select", "zos_version"})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrRecordNotFound
		}
		return "", err
	}

	return setting.Version, nil
}
