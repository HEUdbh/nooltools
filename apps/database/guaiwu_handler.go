package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GuaiwuInfo 怪物基本信息
type GuaiwuInfo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Level     int       `json:"level"`
	Health    int       `json:"health"`
	Attack    int       `json:"attack"`
	Defense   int       `json:"defense"`
	Rewards   string    `json:"rewards"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GuaiwuAttribute 怪物属性
type GuaiwuAttribute struct {
	ID          int       `json:"id"`
	GuaiwuID    int       `json:"guaiwu_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       int       `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GuaiwuSkill 怪物技能
type GuaiwuSkill struct {
	ID          int       `json:"id"`
	GuaiwuID    int       `json:"guaiwu_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAllGuaiwu 获取所有怪物
func (d *Database) GetAllGuaiwu() ([]GuaiwuInfo, error) {
	query := `
	SELECT id, name, type, level, health, attack, defense, rewards, created_at, updated_at
	FROM guaiwu
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询怪物列表失败: %v", err)
	}
	defer rows.Close()

	var guaiwuList []GuaiwuInfo
	for rows.Next() {
		var g GuaiwuInfo
		err := rows.Scan(&g.ID, &g.Name, &g.Type, &g.Level, &g.Health, &g.Attack, &g.Defense, &g.Rewards, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描怪物数据失败: %v", err)
		}
		guaiwuList = append(guaiwuList, g)
	}

	return guaiwuList, nil
}

// GetGuaiwuInfo 获取怪物基本信息
func (d *Database) GetGuaiwuInfo(guaiwuID int) (*GuaiwuInfo, error) {
	query := `
	SELECT id, name, type, level, health, attack, defense, rewards, created_at, updated_at
	FROM guaiwu
	WHERE id = ?`

	var g GuaiwuInfo
	err := d.db.QueryRow(query, guaiwuID).Scan(&g.ID, &g.Name, &g.Type, &g.Level, &g.Health, &g.Attack, &g.Defense, &g.Rewards, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("怪物不存在")
		}
		return nil, fmt.Errorf("查询怪物信息失败: %v", err)
	}

	return &g, nil
}

// GetGuaiwuAttributes 获取怪物属性
func (d *Database) GetGuaiwuAttributes(guaiwuID int) ([]GuaiwuAttribute, error) {
	query := `
	SELECT id, guaiwu_id, name, description, value, created_at, updated_at
	FROM guaiwu_attributes
	WHERE guaiwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, guaiwuID)
	if err != nil {
		return nil, fmt.Errorf("查询怪物属性失败: %v", err)
	}
	defer rows.Close()

	var attributes []GuaiwuAttribute
	for rows.Next() {
		var attr GuaiwuAttribute
		err := rows.Scan(&attr.ID, &attr.GuaiwuID, &attr.Name, &attr.Description, &attr.Value, &attr.CreatedAt, &attr.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描怪物属性数据失败: %v", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetGuaiwuSkills 获取怪物技能
func (d *Database) GetGuaiwuSkills(guaiwuID int) ([]GuaiwuSkill, error) {
	query := `
	SELECT id, guaiwu_id, name, description, created_at, updated_at
	FROM guaiwu_skills
	WHERE guaiwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, guaiwuID)
	if err != nil {
		return nil, fmt.Errorf("查询怪物技能失败: %v", err)
	}
	defer rows.Close()

	var skills []GuaiwuSkill
	for rows.Next() {
		var skill GuaiwuSkill
		err := rows.Scan(&skill.ID, &skill.GuaiwuID, &skill.Name, &skill.Description, &skill.CreatedAt, &skill.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描怪物技能数据失败: %v", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// AddGuaiwuAttribute 添加怪物属性
func (d *Database) AddGuaiwuAttribute(guaiwuID int, name string, description string, value int) error {
	query := `
	INSERT INTO guaiwu_attributes (guaiwu_id, name, description, value)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, guaiwuID, name, description, value)
	if err != nil {
		return fmt.Errorf("添加怪物属性失败: %v", err)
	}

	return nil
}

// DeleteGuaiwuAttribute 删除怪物属性
func (d *Database) DeleteGuaiwuAttribute(attributeID int) error {
	query := `DELETE FROM guaiwu_attributes WHERE id = ?`

	result, err := d.db.Exec(query, attributeID)
	if err != nil {
		return fmt.Errorf("删除怪物属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物属性不存在")
	}

	return nil
}

// UpdateGuaiwuAttribute 更新怪物属性
func (d *Database) UpdateGuaiwuAttribute(attributeID int, name string, description string, value int) error {
	query := `
	UPDATE guaiwu_attributes
	SET name = ?, description = ?, value = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, value, attributeID)
	if err != nil {
		return fmt.Errorf("更新怪物属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物属性不存在")
	}

	return nil
}

// AddGuaiwuSkill 添加怪物技能
func (d *Database) AddGuaiwuSkill(guaiwuID int, name string, description string) error {
	query := `
	INSERT INTO guaiwu_skills (guaiwu_id, name, description)
	VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, guaiwuID, name, description)
	if err != nil {
		return fmt.Errorf("添加怪物技能失败: %v", err)
	}

	return nil
}

// DeleteGuaiwuSkill 删除怪物技能
func (d *Database) DeleteGuaiwuSkill(skillID int) error {
	query := `DELETE FROM guaiwu_skills WHERE id = ?`

	result, err := d.db.Exec(query, skillID)
	if err != nil {
		return fmt.Errorf("删除怪物技能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物技能不存在")
	}

	return nil
}

// UpdateGuaiwuSkill 更新怪物技能
func (d *Database) UpdateGuaiwuSkill(skillID int, name string, description string) error {
	query := `
	UPDATE guaiwu_skills
	SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, skillID)
	if err != nil {
		return fmt.Errorf("更新怪物技能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物技能不存在")
	}

	return nil
}

// CreateGuaiwu 创建怪物
func (d *Database) CreateGuaiwu(name string, guaiwuType string, level int, health int, attack int, defense int, rewards string) (int, error) {
	query := `
	INSERT INTO guaiwu (name, type, level, health, attack, defense, rewards)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, name, guaiwuType, level, health, attack, defense, rewards)
	if err != nil {
		return 0, fmt.Errorf("创建怪物失败: %v", err)
	}

	guaiwuID, _ := result.LastInsertId()

	// 自动添加默认属性：攻击100，防御100
	if err := d.AddGuaiwuAttribute(int(guaiwuID), "攻击", "基础攻击力", 100); err != nil {
		return 0, fmt.Errorf("添加默认攻击属性失败: %v", err)
	}

	if err := d.AddGuaiwuAttribute(int(guaiwuID), "防御", "基础防御力", 100); err != nil {
		return 0, fmt.Errorf("添加默认防御属性失败: %v", err)
	}

	return int(guaiwuID), nil
}

// UpdateGuaiwuBasicInfo 更新怪物基本信息
func (d *Database) UpdateGuaiwuBasicInfo(guaiwuID int, level int, health int, attack int, defense int, rewards string) error {
	query := `
	UPDATE guaiwu
	SET level = ?, health = ?, attack = ?, defense = ?, rewards = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, level, health, attack, defense, rewards, guaiwuID)
	if err != nil {
		return fmt.Errorf("更新怪物基本信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物不存在")
	}

	return nil
}

// DeleteGuaiwu 删除怪物
func (d *Database) DeleteGuaiwu(guaiwuID int) error {
	query := `DELETE FROM guaiwu WHERE id = ?`

	result, err := d.db.Exec(query, guaiwuID)
	if err != nil {
		return fmt.Errorf("删除怪物失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("怪物不存在")
	}

	return nil
}
