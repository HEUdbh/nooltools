// 背包相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// BeibaoInfo 背包信息结构
type BeibaoInfo struct {
	ID    int          `json:"id"`
	Name  string       `json:"name"`
	Items []BeibaoItem `json:"items"`
}

// BeibaoItem 背包物品结构
type BeibaoItem struct {
	ID          int    `json:"id"`
	BeibaoID    int    `json:"beibao_id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
}

// GetAllBeibao 获取所有背包列表
func (d *Database) GetAllBeibao() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name
	FROM beibao
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询背包列表失败: %v", err)
	}
	defer rows.Close()

	var beibaoList []map[string]interface{}
	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("扫描背包数据失败: %v", err)
		}

		beibao := map[string]interface{}{
			"id":   id,
			"name": name,
		}
		beibaoList = append(beibaoList, beibao)
	}

	return beibaoList, nil
}

// GetBeibaoInfo 获取背包详细信息（包含物品）
func (d *Database) GetBeibaoInfo(beibaoID int) (*BeibaoInfo, error) {
	// 查询背包基本信息
	var info BeibaoInfo
	query := `
	SELECT id, name
	FROM beibao
	WHERE id = ?`

	err := d.db.QueryRow(query, beibaoID).Scan(&info.ID, &info.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("背包不存在")
		}
		return nil, fmt.Errorf("查询背包信息失败: %v", err)
	}

	// 查询背包物品
	items, err := d.GetBeibaoItems(beibaoID)
	if err != nil {
		return nil, fmt.Errorf("查询背包物品失败: %v", err)
	}
	info.Items = items

	return &info, nil
}

// CreateBeibao 创建背包
func (d *Database) CreateBeibao(name string) (int64, error) {
	query := `
	INSERT INTO beibao (name)
	VALUES (?)`

	result, err := d.db.Exec(query, name)
	if err != nil {
		return 0, fmt.Errorf("创建背包失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取背包ID失败: %v", err)
	}

	return id, nil
}

// DeleteBeibao 删除背包
func (d *Database) DeleteBeibao(beibaoID int) error {
	// 先删除背包中的所有物品
	_, err := d.db.Exec("DELETE FROM beibao_items WHERE beibao_id = ?", beibaoID)
	if err != nil {
		return fmt.Errorf("删除背包物品失败: %v", err)
	}

	// 删除背包
	query := `DELETE FROM beibao WHERE id = ?`

	result, err := d.db.Exec(query, beibaoID)
	if err != nil {
		return fmt.Errorf("删除背包失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("背包不存在")
	}

	return nil
}

// UpdateBeibao 更新背包名称
func (d *Database) UpdateBeibao(beibaoID int, name string) error {
	query := `
	UPDATE beibao
	SET name = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, beibaoID)
	if err != nil {
		return fmt.Errorf("更新背包失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("背包不存在")
	}

	return nil
}

// GetBeibaoItems 获取背包物品列表
func (d *Database) GetBeibaoItems(beibaoID int) ([]BeibaoItem, error) {
	query := `
	SELECT id, beibao_id, name, quantity, description
	FROM beibao_items
	WHERE beibao_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, beibaoID)
	if err != nil {
		return nil, fmt.Errorf("查询物品列表失败: %v", err)
	}
	defer rows.Close()

	var items []BeibaoItem
	for rows.Next() {
		var item BeibaoItem
		if err := rows.Scan(&item.ID, &item.BeibaoID, &item.Name, &item.Quantity, &item.Description); err != nil {
			return nil, fmt.Errorf("扫描物品数据失败: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// AddBeibaoItem 添加背包物品
func (d *Database) AddBeibaoItem(beibaoID int, name string, quantity int, description string) error {
	query := `
	INSERT INTO beibao_items (beibao_id, name, quantity, description)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, beibaoID, name, quantity, description)
	if err != nil {
		return fmt.Errorf("添加物品失败: %v", err)
	}

	return nil
}

// DeleteBeibaoItem 删除背包物品
func (d *Database) DeleteBeibaoItem(itemID int) error {
	query := `DELETE FROM beibao_items WHERE id = ?`

	result, err := d.db.Exec(query, itemID)
	if err != nil {
		return fmt.Errorf("删除物品失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("物品不存在")
	}

	return nil
}

// UpdateBeibaoItem 更新背包物品
func (d *Database) UpdateBeibaoItem(itemID int, name string, quantity int, description string) error {
	query := `
	UPDATE beibao_items
	SET name = ?, quantity = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, quantity, description, itemID)
	if err != nil {
		return fmt.Errorf("更新物品失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("物品不存在")
	}

	return nil
}
