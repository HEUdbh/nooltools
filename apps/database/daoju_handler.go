// 道具相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// DaojuFunction 道具功能结构
type DaojuFunction struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DaojuInfo 道具信息结构
type DaojuInfo struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Level     int             `json:"level"`
	Holder    string          `json:"holder"`
	Functions []DaojuFunction `json:"functions"`
}

// GetAllDaoju 获取所有道具列表
func (d *Database) GetAllDaoju() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name, level, holder
	FROM daoju
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询道具列表失败: %v", err)
	}
	defer rows.Close()

	var daojuList []map[string]interface{}
	for rows.Next() {
		var id int
		var name, holder string
		var level int

		if err := rows.Scan(&id, &name, &level, &holder); err != nil {
			return nil, fmt.Errorf("扫描道具数据失败: %v", err)
		}

		daoju := map[string]interface{}{
			"id":     id,
			"name":   name,
			"level":  level,
			"holder": holder,
		}
		daojuList = append(daojuList, daoju)
	}

	return daojuList, nil
}

// GetDaojuInfo 获取道具详细信息（包含功能）
func (d *Database) GetDaojuInfo(daojuID int) (*DaojuInfo, error) {
	// 查询道具基本信息
	var info DaojuInfo
	query := `
	SELECT id, name, level, holder
	FROM daoju
	WHERE id = ?`

	err := d.db.QueryRow(query, daojuID).Scan(&info.ID, &info.Name, &info.Level, &info.Holder)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("道具不存在")
		}
		return nil, fmt.Errorf("查询道具信息失败: %v", err)
	}

	// 查询道具功能
	functions, err := d.GetDaojuFunctions(daojuID)
	if err != nil {
		return nil, fmt.Errorf("查询道具功能失败: %v", err)
	}
	info.Functions = functions

	return &info, nil
}

// CreateDaoju 创建道具
func (d *Database) CreateDaoju(name string, level int, holder string) (int64, error) {
	query := `
	INSERT INTO daoju (name, level, holder)
	VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, name, level, holder)
	if err != nil {
		return 0, fmt.Errorf("创建道具失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取道具ID失败: %v", err)
	}

	return id, nil
}

// DeleteDaoju 删除道具
func (d *Database) DeleteDaoju(daojuID int) error {
	query := `DELETE FROM daoju WHERE id = ?`

	result, err := d.db.Exec(query, daojuID)
	if err != nil {
		return fmt.Errorf("删除道具失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("道具不存在")
	}

	return nil
}

// UpdateDaojuBasicInfo 更新道具基本信息（等级、所有人）
func (d *Database) UpdateDaojuBasicInfo(daojuID int, level int, holder string) error {
	query := `
	UPDATE daoju
	SET level = ?, holder = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, level, holder, daojuID)
	if err != nil {
		return fmt.Errorf("更新道具基本信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("道具不存在")
	}

	return nil
}

// GetDaojuFunctions 获取道具功能列表
func (d *Database) GetDaojuFunctions(daojuID int) ([]DaojuFunction, error) {
	query := `
	SELECT id, name, description
	FROM daoju_functions
	WHERE daoju_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, daojuID)
	if err != nil {
		return nil, fmt.Errorf("查询功能列表失败: %v", err)
	}
	defer rows.Close()

	var functions []DaojuFunction
	for rows.Next() {
		var function DaojuFunction
		if err := rows.Scan(&function.ID, &function.Name, &function.Description); err != nil {
			return nil, fmt.Errorf("扫描功能数据失败: %v", err)
		}
		functions = append(functions, function)
	}

	return functions, nil
}

// AddDaojuFunction 添加道具功能
func (d *Database) AddDaojuFunction(daojuID int, name, description string) error {
	query := `
	INSERT INTO daoju_functions (daoju_id, name, description)
	VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, daojuID, name, description)
	if err != nil {
		return fmt.Errorf("添加功能失败: %v", err)
	}

	return nil
}

// DeleteDaojuFunction 删除道具功能
func (d *Database) DeleteDaojuFunction(functionID int) error {
	query := `DELETE FROM daoju_functions WHERE id = ?`

	result, err := d.db.Exec(query, functionID)
	if err != nil {
		return fmt.Errorf("删除功能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("功能不存在")
	}

	return nil
}

// UpdateDaojuFunction 更新道具功能
func (d *Database) UpdateDaojuFunction(functionID int, name, description string) error {
	query := `
	UPDATE daoju_functions
	SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, functionID)
	if err != nil {
		return fmt.Errorf("更新功能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("功能不存在")
	}

	return nil
}
