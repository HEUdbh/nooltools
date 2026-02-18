// 宠物相关的后端接口处理
package database

import (
	"database/sql"
	"fmt"
)

// PetAttribute 宠物属性结构
type PetAttribute struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

// PetSkill 宠物技能结构
type PetSkill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PetInfo 宠物信息结构
type PetInfo struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Level      int            `json:"level"`
	Owner      string         `json:"owner"`
	Attributes []PetAttribute `json:"attributes"`
	Skills     []PetSkill     `json:"skills"`
}

// GetAllPets 获取所有宠物列表
func (d *Database) GetAllPets() ([]map[string]interface{}, error) {
	query := `
	SELECT id, name, level, owner
	FROM chongwu
	ORDER BY id ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询宠物列表失败: %v", err)
	}
	defer rows.Close()

	var pets []map[string]interface{}
	for rows.Next() {
		var id int
		var name, owner string
		var level int

		if err := rows.Scan(&id, &name, &level, &owner); err != nil {
			return nil, fmt.Errorf("扫描宠物数据失败: %v", err)
		}

		pet := map[string]interface{}{
			"id":    id,
			"name":  name,
			"level": level,
			"owner": owner,
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

// GetPetInfo 获取宠物详细信息（包含属性和技能）
func (d *Database) GetPetInfo(petID int) (*PetInfo, error) {
	// 查询宠物基本信息
	var info PetInfo
	query := `
	SELECT id, name, level, owner
	FROM chongwu
	WHERE id = ?`

	err := d.db.QueryRow(query, petID).Scan(&info.ID, &info.Name, &info.Level, &info.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("宠物不存在")
		}
		return nil, fmt.Errorf("查询宠物信息失败: %v", err)
	}

	// 查询宠物属性
	attributes, err := d.GetPetAttributes(petID)
	if err != nil {
		return nil, fmt.Errorf("查询宠物属性失败: %v", err)
	}
	info.Attributes = attributes

	// 查询宠物技能
	skills, err := d.GetPetSkills(petID)
	if err != nil {
		return nil, fmt.Errorf("查询宠物技能失败: %v", err)
	}
	info.Skills = skills

	return &info, nil
}

// GetPetAttributes 获取宠物属性列表
func (d *Database) GetPetAttributes(petID int) ([]PetAttribute, error) {
	query := `
	SELECT id, name, description, value
	FROM chongwu_attributes
	WHERE chongwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, petID)
	if err != nil {
		return nil, fmt.Errorf("查询属性列表失败: %v", err)
	}
	defer rows.Close()

	var attributes []PetAttribute
	for rows.Next() {
		var attr PetAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Description, &attr.Value); err != nil {
			return nil, fmt.Errorf("扫描属性数据失败: %v", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetPetSkills 获取宠物技能列表
func (d *Database) GetPetSkills(petID int) ([]PetSkill, error) {
	query := `
	SELECT id, name, description
	FROM chongwu_skills
	WHERE chongwu_id = ?
	ORDER BY id ASC`

	rows, err := d.db.Query(query, petID)
	if err != nil {
		return nil, fmt.Errorf("查询技能列表失败: %v", err)
	}
	defer rows.Close()

	var skills []PetSkill
	for rows.Next() {
		var skill PetSkill
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Description); err != nil {
			return nil, fmt.Errorf("扫描技能数据失败: %v", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// AddPetAttribute 添加宠物属性
func (d *Database) AddPetAttribute(petID int, name, description string, value int) error {
	query := `
	INSERT INTO chongwu_attributes (chongwu_id, name, description, value)
	VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, petID, name, description, value)
	if err != nil {
		return fmt.Errorf("添加属性失败: %v", err)
	}

	return nil
}

// DeletePetAttribute 删除宠物属性
func (d *Database) DeletePetAttribute(attributeID int) error {
	query := `DELETE FROM chongwu_attributes WHERE id = ?`

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

// UpdatePetAttribute 更新宠物属性
func (d *Database) UpdatePetAttribute(attributeID int, name, description string, value int) error {
	query := `
	UPDATE chongwu_attributes
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

// AddPetSkill 添加宠物技能
func (d *Database) AddPetSkill(petID int, name, description string) error {
	query := `
	INSERT INTO chongwu_skills (chongwu_id, name, description)
	VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, petID, name, description)
	if err != nil {
		return fmt.Errorf("添加技能失败: %v", err)
	}

	return nil
}

// DeletePetSkill 删除宠物技能
func (d *Database) DeletePetSkill(skillID int) error {
	query := `DELETE FROM chongwu_skills WHERE id = ?`

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

// UpdatePetSkill 更新宠物技能
func (d *Database) UpdatePetSkill(skillID int, name, description string) error {
	query := `
	UPDATE chongwu_skills
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

// CreatePet 创建新宠物
func (d *Database) CreatePet(name, owner string, level int) (int, error) {
	query := `
	INSERT INTO chongwu (name, owner, level)
	VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, name, owner, level)
	if err != nil {
		return 0, fmt.Errorf("创建宠物失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取宠物ID失败: %v", err)
	}

	// 为新宠物添加默认属性
	defaultAttributes := []struct {
		name        string
		description string
		value       int
	}{
		{"攻击", "攻击力", 100},
		{"防御", "防御力", 100},
		{"忠诚", "忠诚度", 100},
	}

	for _, attr := range defaultAttributes {
		err = d.AddPetAttribute(int(id), attr.name, attr.description, attr.value)
		if err != nil {
			return 0, fmt.Errorf("添加默认属性失败: %v", err)
		}
	}

	return int(id), nil
}

// UpdatePetBasicInfo 更新宠物基本信息（名称不可变更）
func (d *Database) UpdatePetBasicInfo(petID int, owner string, level int) error {
	query := `
	UPDATE chongwu
	SET owner = ?, level = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := d.db.Exec(query, owner, level, petID)
	if err != nil {
		return fmt.Errorf("更新宠物信息失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("宠物不存在")
	}

	return nil
}

// DeletePet 删除宠物
func (d *Database) DeletePet(petID int) error {
	query := `DELETE FROM chongwu WHERE id = ?`

	result, err := d.db.Exec(query, petID)
	if err != nil {
		return fmt.Errorf("删除宠物失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("宠物不存在")
	}

	return nil
}
