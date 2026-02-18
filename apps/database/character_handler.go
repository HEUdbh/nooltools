// 人物相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// CharacterAttribute 人物属性结构
type CharacterAttribute struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

// CharacterSkill 人物技能结构
type CharacterSkill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CharacterInfo 人物信息结构
type CharacterInfo struct {
	ID         int                  `json:"id"`
	Name       string               `json:"name"`
	Shili      string               `json:"shili"`
	Property   int                  `json:"property"`
	Level      int                  `json:"level"`
	Attributes []CharacterAttribute `json:"attributes"`
	Skills     []CharacterSkill     `json:"skills"`
}

// GetAllCharacters 获取所有人物列表
func (d *Database) GetAllCharacters() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name, shili, property, level
	FROM renwu
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询人物列表失败: %v", err)
	}
	defer rows.Close()

	var characters []map[string]interface{}
	for rows.Next() {
		var id int
		var name, shili string
		var property, level int

		if err := rows.Scan(&id, &name, &shili, &property, &level); err != nil {
			return nil, fmt.Errorf("扫描人物数据失败: %v", err)
		}

		character := map[string]interface{}{
			"id":       id,
			"name":     name,
			"shili":    shili,
			"property": property,
			"level":    level,
		}
		characters = append(characters, character)
	}

	return characters, nil
}

// GetCharacterInfo 获取人物详细信息（包含属性和技能）
func (d *Database) GetCharacterInfo(characterID int) (*CharacterInfo, error) {
	// 查询人物基本信息
	var info CharacterInfo
	query := `
	SELECT id, name, shili, property, level
	FROM renwu
	WHERE id = ?`

	err := d.db.QueryRow(query, characterID).Scan(&info.ID, &info.Name, &info.Shili, &info.Property, &info.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("人物不存在")
		}
		return nil, fmt.Errorf("查询人物信息失败: %v", err)
	}

	// 查询人物属性
	attributes, err := d.GetCharacterAttributes(characterID)
	if err != nil {
		return nil, fmt.Errorf("查询人物属性失败: %v", err)
	}
	info.Attributes = attributes

	// 查询人物技能
	skills, err := d.GetCharacterSkills(characterID)
	if err != nil {
		return nil, fmt.Errorf("查询人物技能失败: %v", err)
	}
	info.Skills = skills

	return &info, nil
}

// GetCharacterAttributes 获取人物属性列表
func (d *Database) GetCharacterAttributes(characterID int) ([]CharacterAttribute, error) {
	query := `
	SELECT id, name, description, value
	FROM renwu_attributes
	WHERE renwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, characterID)
	if err != nil {
		return nil, fmt.Errorf("查询属性列表失败: %v", err)
	}
	defer rows.Close()

	var attributes []CharacterAttribute
	for rows.Next() {
		var attr CharacterAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Description, &attr.Value); err != nil {
			return nil, fmt.Errorf("扫描属性数据失败: %v", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetCharacterSkills 获取人物技能列表
func (d *Database) GetCharacterSkills(characterID int) ([]CharacterSkill, error) {
	query := `
	SELECT id, name, description
	FROM renwu_skills
	WHERE renwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, characterID)
	if err != nil {
		return nil, fmt.Errorf("查询技能列表失败: %v", err)
	}
	defer rows.Close()

	var skills []CharacterSkill
	for rows.Next() {
		var skill CharacterSkill
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Description); err != nil {
			return nil, fmt.Errorf("扫描技能数据失败: %v", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// AddCharacterAttribute 添加人物属性
func (d *Database) AddCharacterAttribute(characterID int, name, description string, value int) error {
	query := `
	INSERT INTO renwu_attributes (renwu_id, name, description, value)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, characterID, name, description, value)
	if err != nil {
		return fmt.Errorf("添加属性失败: %v", err)
	}

	return nil
}

// DeleteCharacterAttribute 删除人物属性
func (d *Database) DeleteCharacterAttribute(attributeID int) error {
	query := `DELETE FROM renwu_attributes WHERE id = ?`

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

// UpdateCharacterAttribute 更新人物属性
func (d *Database) UpdateCharacterAttribute(attributeID int, name, description string, value int) error {
	query := `
	UPDATE renwu_attributes
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

// AddCharacterSkill 添加人物技能
func (d *Database) AddCharacterSkill(characterID int, name, description string) error {
	query := `
	INSERT INTO renwu_skills (renwu_id, name, description)
	VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, characterID, name, description)
	if err != nil {
		return fmt.Errorf("添加技能失败: %v", err)
	}

	return nil
}

// DeleteCharacterSkill 删除人物技能
func (d *Database) DeleteCharacterSkill(skillID int) error {
	query := `DELETE FROM renwu_skills WHERE id = ?`

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

// UpdateCharacterSkill 更新人物技能
func (d *Database) UpdateCharacterSkill(skillID int, name, description string) error {
	query := `
	UPDATE renwu_skills
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

// CreateCharacter 创建新人物
func (d *Database) CreateCharacter(name, shili string, property, level int) (int, error) {
	query := `
	INSERT INTO renwu (name, shili, property, level)
	VALUES (?, ?, ?, ?)`

	result, err := d.db.Exec(query, name, shili, property, level)
	if err != nil {
		return 0, fmt.Errorf("创建人物失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取人物ID失败: %v", err)
	}

	// 为新人物添加默认属性
	defaultAttributes := []struct {
		name        string
		description string
		value       int
	}{
		{"血量", "当前生命值", 100},
		{"蓝量", "当前魔法值", 100},
		{"攻击", "攻击力", 100},
		{"防御", "防御力", 100},
	}

	for _, attr := range defaultAttributes {
		err = d.AddCharacterAttribute(int(id), attr.name, attr.description, attr.value)
		if err != nil {
			return 0, fmt.Errorf("添加默认属性失败: %v", err)
		}
	}

	return int(id), nil
}

// UpdateCharacterBasicInfo 更新人物基本信息
func (d *Database) UpdateCharacterBasicInfo(characterID int, name, shili string, property, level int) error {
	query := `
	UPDATE renwu
	SET name = ?, shili = ?, property = ?, level = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, name, shili, property, level, characterID)
	if err != nil {
		return fmt.Errorf("更新人物信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("人物不存在")
	}

	return nil
}

// DeleteCharacter 删除人物
func (d *Database) DeleteCharacter(characterID int) error {
	query := `DELETE FROM renwu WHERE id = ?`

	result, err := d.db.Exec(query, characterID)
	if err != nil {
		return fmt.Errorf("删除人物失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("人物不存在")
	}

	return nil
}
