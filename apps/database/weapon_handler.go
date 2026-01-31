// 武器相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// WeaponAttribute 武器属性结构
type WeaponAttribute struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

// WeaponSkill 武器技能结构
type WeaponSkill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WeaponInfo 武器信息结构
type WeaponInfo struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Holder     string            `json:"holder"`
	Level      int               `json:"level"`
	Attributes []WeaponAttribute `json:"attributes"`
	Skills     []WeaponSkill     `json:"skills"`
}

// GetAllWeapons 获取所有武器列表
func (d *Database) GetAllWeapons() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name, holder, level
	FROM wuqi
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询武器列表失败: %v", err)
	}
	defer rows.Close()

	var weapons []map[string]interface{}
	for rows.Next() {
		var id int
		var name, holder string
		var level int

		if err := rows.Scan(&id, &name, &holder, &level); err != nil {
			return nil, fmt.Errorf("扫描武器数据失败: %v", err)
		}

		weapon := map[string]interface{}{
			"id":     id,
			"name":   name,
			"holder": holder,
			"level":  level,
		}
		weapons = append(weapons, weapon)
	}

	return weapons, nil
}

// GetWeaponInfo 获取武器详细信息（包含属性和技能）
func (d *Database) GetWeaponInfo(weaponID int) (*WeaponInfo, error) {
	// 查询武器基本信息
	var info WeaponInfo
	query := `
	SELECT id, name, holder, level
	FROM wuqi
	WHERE id = ?`

	err := d.db.QueryRow(query, weaponID).Scan(&info.ID, &info.Name, &info.Holder, &info.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("武器不存在")
		}
		return nil, fmt.Errorf("查询武器信息失败: %v", err)
	}

	// 查询武器属性
	attributes, err := d.GetWeaponAttributes(weaponID)
	if err != nil {
		return nil, fmt.Errorf("查询武器属性失败: %v", err)
	}
	info.Attributes = attributes

	// 查询武器技能
	skills, err := d.GetWeaponSkills(weaponID)
	if err != nil {
		return nil, fmt.Errorf("查询武器技能失败: %v", err)
	}
	info.Skills = skills

	return &info, nil
}

// GetWeaponAttributes 获取武器属性列表
func (d *Database) GetWeaponAttributes(weaponID int) ([]WeaponAttribute, error) {
	query := `
	SELECT id, name, description, value
	FROM wuqi_attributes
	WHERE wuqi_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, weaponID)
	if err != nil {
		return nil, fmt.Errorf("查询属性列表失败: %v", err)
	}
	defer rows.Close()

	var attributes []WeaponAttribute
	for rows.Next() {
		var attr WeaponAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Description, &attr.Value); err != nil {
			return nil, fmt.Errorf("扫描属性数据失败: %v", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetWeaponSkills 获取武器技能列表
func (d *Database) GetWeaponSkills(weaponID int) ([]WeaponSkill, error) {
	query := `
	SELECT id, name, description
	FROM wuqi_skills
	WHERE wuqi_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, weaponID)
	if err != nil {
		return nil, fmt.Errorf("查询技能列表失败: %v", err)
	}
	defer rows.Close()

	var skills []WeaponSkill
	for rows.Next() {
		var skill WeaponSkill
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Description); err != nil {
			return nil, fmt.Errorf("扫描技能数据失败: %v", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// AddWeaponAttribute 添加武器属性
func (d *Database) AddWeaponAttribute(weaponID int, name, description string, value int) error {
	query := `
	INSERT INTO wuqi_attributes (wuqi_id, name, description, value)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, weaponID, name, description, value)
	if err != nil {
		return fmt.Errorf("添加属性失败: %v", err)
	}

	return nil
}

// DeleteWeaponAttribute 删除武器属性
func (d *Database) DeleteWeaponAttribute(attributeID int) error {
	query := `DELETE FROM wuqi_attributes WHERE id = ?`

	result, err := d.db.Exec(query, attributeID)
	if err != nil {
		return fmt.Errorf("删除属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("属性不存在")
	}

	return nil
}

// UpdateWeaponAttribute 更新武器属性
func (d *Database) UpdateWeaponAttribute(attributeID int, name, description string, value int) error {
	query := `
	UPDATE wuqi_attributes
	SET name = ?, description = ?, value = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, value, attributeID)
	if err != nil {
		return fmt.Errorf("更新属性失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("属性不存在")
	}

	return nil
}

// AddWeaponSkill 添加武器技能
func (d *Database) AddWeaponSkill(weaponID int, name, description string) error {
	query := `
	INSERT INTO wuqi_skills (wuqi_id, name, description)
	VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, weaponID, name, description)
	if err != nil {
		return fmt.Errorf("添加技能失败: %v", err)
	}

	return nil
}

// DeleteWeaponSkill 删除武器技能
func (d *Database) DeleteWeaponSkill(skillID int) error {
	query := `DELETE FROM wuqi_skills WHERE id = ?`

	result, err := d.db.Exec(query, skillID)
	if err != nil {
		return fmt.Errorf("删除技能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("技能不存在")
	}

	return nil
}

// UpdateWeaponSkill 更新武器技能
func (d *Database) UpdateWeaponSkill(skillID int, name, description string) error {
	query := `
	UPDATE wuqi_skills
	SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, description, skillID)
	if err != nil {
		return fmt.Errorf("更新技能失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("技能不存在")
	}

	return nil
}

// CreateWeapon 创建新武器
func (d *Database) CreateWeapon(name, holder string, level int) (int, error) {
	query := `
	INSERT INTO wuqi (name, holder, level)
	VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, name, holder, level)
	if err != nil {
		return 0, fmt.Errorf("创建武器失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取武器ID失败: %v", err)
	}

	// 为新武器添加默认属性：耐久100
	err = d.AddWeaponAttribute(int(id), "耐久", "武器耐久度", 100)
	if err != nil {
		return 0, fmt.Errorf("添加默认属性失败: %v", err)
	}

	return int(id), nil
}

// UpdateWeaponBasicInfo 更新武器基本信息
func (d *Database) UpdateWeaponBasicInfo(weaponID int, name, holder string, level int) error {
	query := `
	UPDATE wuqi
	SET name = ?, holder = ?, level = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, holder, level, weaponID)
	if err != nil {
		return fmt.Errorf("更新武器信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("武器不存在")
	}

	return nil
}

// DeleteWeapon 删除武器
func (d *Database) DeleteWeapon(weaponID int) error {
	query := `DELETE FROM wuqi WHERE id = ?`

	result, err := d.db.Exec(query, weaponID)
	if err != nil {
		return fmt.Errorf("删除武器失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("武器不存在")
	}

	return nil
}
