// 任务相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// ShiqingDetail 任务详情结构
type ShiqingDetail struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// ShiqingInfo 任务信息结构
type ShiqingInfo struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Location string          `json:"location"`
	Time     string          `json:"time"`
	Details  []ShiqingDetail `json:"details"`
}

// GetAllShiqing 获取所有任务列表
func (d *Database) GetAllShiqing() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name, location, time
	FROM shiqing
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %v", err)
	}
	defer rows.Close()

	var shiqingList []map[string]interface{}
	for rows.Next() {
		var id int
		var name, location, time string

		if err := rows.Scan(&id, &name, &location, &time); err != nil {
			return nil, fmt.Errorf("扫描任务数据失败: %v", err)
		}

		shiqing := map[string]interface{}{
			"id":       id,
			"name":     name,
			"location": location,
			"time":     time,
		}
		shiqingList = append(shiqingList, shiqing)
	}

	return shiqingList, nil
}

// GetShiqingInfo 获取任务详细信息（包含详情）
func (d *Database) GetShiqingInfo(shiqingID int) (*ShiqingInfo, error) {
	// 查询任务基本信息
	var info ShiqingInfo
	query := `
	SELECT id, name, location, time
	FROM shiqing
	WHERE id = ?`

	err := d.db.QueryRow(query, shiqingID).Scan(&info.ID, &info.Name, &info.Location, &info.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("任务不存在")
		}
		return nil, fmt.Errorf("查询任务信息失败: %v", err)
	}

	// 查询任务详情
	details, err := d.GetShiqingDetails(shiqingID)
	if err != nil {
		return nil, fmt.Errorf("查询任务详情失败: %v", err)
	}
	info.Details = details

	return &info, nil
}

// CreateShiqing 创建任务
func (d *Database) CreateShiqing(name, location, time string) (int64, error) {
	query := `
	INSERT INTO shiqing (name, location, time)
	VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, name, location, time)
	if err != nil {
		return 0, fmt.Errorf("创建任务失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取任务ID失败: %v", err)
	}

	return id, nil
}

// DeleteShiqing 删除任务
func (d *Database) DeleteShiqing(shiqingID int) error {
	query := `DELETE FROM shiqing WHERE id = ?`

	result, err := d.db.Exec(query, shiqingID)
	if err != nil {
		return fmt.Errorf("删除任务失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("任务不存在")
	}

	return nil
}

// UpdateShiqingBasicInfo 更新任务基本信息（地点、时间）
func (d *Database) UpdateShiqingBasicInfo(shiqingID int, location, time string) error {
	query := `
	UPDATE shiqing
	SET location = ?, time = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, location, time, shiqingID)
	if err != nil {
		return fmt.Errorf("更新任务基本信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("任务不存在")
	}

	return nil
}

// GetShiqingDetails 获取任务详情列表
func (d *Database) GetShiqingDetails(shiqingID int) ([]ShiqingDetail, error) {
	query := `
	SELECT id, description
	FROM shiqing_details
	WHERE shiqing_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, shiqingID)
	if err != nil {
		return nil, fmt.Errorf("查询详情列表失败: %v", err)
	}
	defer rows.Close()

	var details []ShiqingDetail
	for rows.Next() {
		var detail ShiqingDetail
		if err := rows.Scan(&detail.ID, &detail.Description); err != nil {
			return nil, fmt.Errorf("扫描详情数据失败: %v", err)
		}
		details = append(details, detail)
	}

	return details, nil
}

// AddShiqingDetail 添加任务详情
func (d *Database) AddShiqingDetail(shiqingID int, description string) error {
	query := `
	INSERT INTO shiqing_details (shiqing_id, description)
	VALUES (?, ?)`

	_, err := d.db.Exec(query, shiqingID, description)
	if err != nil {
		return fmt.Errorf("添加详情失败: %v", err)
	}

	return nil
}

// DeleteShiqingDetail 删除任务详情
func (d *Database) DeleteShiqingDetail(detailID int) error {
	query := `DELETE FROM shiqing_details WHERE id = ?`

	result, err := d.db.Exec(query, detailID)
	if err != nil {
		return fmt.Errorf("删除详情失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("详情不存在")
	}

	return nil
}

// UpdateShiqingDetail 更新任务详情
func (d *Database) UpdateShiqingDetail(detailID int, description string) error {
	query := `
	UPDATE shiqing_details
	SET description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, description, detailID)
	if err != nil {
		return fmt.Errorf("更新详情失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("详情不存在")
	}

	return nil
}
