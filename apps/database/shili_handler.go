package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ShiliInfo 势力基本信息
type ShiliInfo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Level       int       `json:"level"`
	Founder     string    `json:"founder"`
	Wealth      int       `json:"wealth"`
	MemberCount int       `json:"member_count"`
	MaxMembers  int       `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShiliPosition 势力职务
type ShiliPosition struct {
	ID           int       `json:"id"`
	ShiliID      int       `json:"shili_id"`
	PositionName string    `json:"position_name"`
	PersonName   string    `json:"person_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ShiliAttribute 势力属性
type ShiliAttribute struct {
	ID          int       `json:"id"`
	ShiliID     int       `json:"shili_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       int       `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAllShili 获取所有势力
func (d *Database) GetAllShili() ([]ShiliInfo, error) {
	query := `
	SELECT id, name, level, founder, wealth, member_count, max_members, created_at, updated_at
	FROM shili
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询势力列表失败: %v", err)
	}
	defer rows.Close()

	var shiliList []ShiliInfo
	for rows.Next() {
		var s ShiliInfo
		err := rows.Scan(&s.ID, &s.Name, &s.Level, &s.Founder, &s.Wealth, &s.MemberCount, &s.MaxMembers, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描势力数据失败: %v", err)
		}
		shiliList = append(shiliList, s)
	}

	return shiliList, nil
}

// GetShiliInfo 获取势力基本信息
func (d *Database) GetShiliInfo(shiliID int) (*ShiliInfo, error) {
	query := `
	SELECT id, name, level, founder, wealth, member_count, max_members, created_at, updated_at
	FROM shili
	WHERE id = ?`

	var s ShiliInfo
	err := d.db.QueryRow(query, shiliID).Scan(&s.ID, &s.Name, &s.Level, &s.Founder, &s.Wealth, &s.MemberCount, &s.MaxMembers, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("势力不存在")
		}
		return nil, fmt.Errorf("查询势力信息失败: %v", err)
	}

	return &s, nil
}

// GetShiliPositions 获取势力职务
func (d *Database) GetShiliPositions(shiliID int) ([]ShiliPosition, error) {
	query := `
	SELECT id, shili_id, position_name, person_name, description, created_at, updated_at
	FROM shili_positions
	WHERE shili_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, shiliID)
	if err != nil {
		return nil, fmt.Errorf("查询势力职务失败: %v", err)
	}
	defer rows.Close()

	var positions []ShiliPosition
	for rows.Next() {
		var pos ShiliPosition
		err := rows.Scan(&pos.ID, &pos.ShiliID, &pos.PositionName, &pos.PersonName, &pos.Description, &pos.CreatedAt, &pos.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描势力职务数据失败: %v", err)
		}
		positions = append(positions, pos)
	}

	return positions, nil
}

// GetShiliAttributes 获取势力属性
func (d *Database) GetShiliAttributes(shiliID int) ([]ShiliAttribute, error) {
	query := `
	SELECT id, shili_id, name, description, value, created_at, updated_at
	FROM shili_attributes
	WHERE shili_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, shiliID)
	if err != nil {
		return nil, fmt.Errorf("查询势力属性失败: %v", err)
	}
	defer rows.Close()

	var attributes []ShiliAttribute
	for rows.Next() {
		var attr ShiliAttribute
		err := rows.Scan(&attr.ID, &attr.ShiliID, &attr.Name, &attr.Description, &attr.Value, &attr.CreatedAt, &attr.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描势力属性数据失败: %v", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// CreateShili 创建势力
func (d *Database) CreateShili(name string, founder string, level int, wealth int, maxMembers int) (int64, error) {
	query := `
	INSERT INTO shili (name, founder, level, wealth, member_count, max_members)
	VALUES (?, ?, ?, ?, 0, ?)`

	result, err := d.db.Exec(query, name, founder, level, wealth, maxMembers)
	if err != nil {
		return 0, fmt.Errorf("创建势力失败: %v", err)
	}

	shiliID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取势力ID失败: %v", err)
	}

	return shiliID, nil
}

// UpdateShiliBasicInfo 更新势力基本信息
func (d *Database) UpdateShiliBasicInfo(shiliID int, level int, founder string, wealth int, memberCount int, maxMembers int) error {
	// 验证现有人数不能超过容纳人数
	if memberCount > maxMembers {
		return fmt.Errorf("现有人数不能超过容纳人数")
	}

	query := `
	UPDATE shili
	SET level = ?, founder = ?, wealth = ?, member_count = ?, max_members = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, level, founder, wealth, memberCount, maxMembers, shiliID)
	if err != nil {
		return fmt.Errorf("更新势力基本信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力不存在")
	}

	return nil
}

// DeleteShili 删除势力
func (d *Database) DeleteShili(shiliID int) error {
	query := `DELETE FROM shili WHERE id = ?`

	result, err := d.db.Exec(query, shiliID)
	if err != nil {
		return fmt.Errorf("删除势力失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力不存在")
	}

	return nil
}

// AddShiliPosition 添加势力职务
func (d *Database) AddShiliPosition(shiliID int, positionName string, personName string, description string) error {
	query := `
	INSERT INTO shili_positions (shili_id, position_name, person_name, description)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, shiliID, positionName, personName, description)
	if err != nil {
		return fmt.Errorf("添加势力职务失败: %v", err)
	}

	return nil
}

// DeleteShiliPosition 删除势力职务
func (d *Database) DeleteShiliPosition(positionID int) error {
	query := `DELETE FROM shili_positions WHERE id = ?`

	result, err := d.db.Exec(query, positionID)
	if err != nil {
		return fmt.Errorf("删除势力职务失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力职务不存在")
	}

	return nil
}

// UpdateShiliPosition 更新势力职务
func (d *Database) UpdateShiliPosition(positionID int, positionName string, personName string, description string) error {
	query := `
	UPDATE shili_positions
	SET position_name = ?, person_name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, positionName, personName, description, positionID)
	if err != nil {
		return fmt.Errorf("更新势力职务失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力职务不存在")
	}

	return nil
}

// AddShiliAttribute 添加势力属性
func (d *Database) AddShiliAttribute(shiliID int, name string, description string, value int) error {
	query := `
	INSERT INTO shili_attributes (shili_id, name, description, value)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, shiliID, name, description, value)
	if err != nil {
		return fmt.Errorf("添加势力属性失败: %v", err)
	}

	return nil
}

// DeleteShiliAttribute 删除势力属性
func (d *Database) DeleteShiliAttribute(attributeID int) error {
	query := `DELETE FROM shili_attributes WHERE id = ?`

	result, err := d.db.Exec(query, attributeID)
	if err != nil {
		return fmt.Errorf("删除势力属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力属性不存在")
	}

	return nil
}

// UpdateShiliAttribute 更新势力属性
func (d *Database) UpdateShiliAttribute(attributeID int, name string, description string, value int) error {
	query := `
	UPDATE shili_attributes
	SET name = ?, description = ?, value = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, value, attributeID)
	if err != nil {
		return fmt.Errorf("更新势力属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("势力属性不存在")
	}

	return nil
}
