package main

import (
	"context"
	"fmt"
	"log"
	"nooltools/apps/database"
)

// app 结构体
type app struct {
	database *database.Database
}

// NewApp 创建应用实例
func NewApp() *app {
	return &app{}
}

// startup 应用启动时调用
func (a *app) startup(ctx context.Context) {
	// 初始化数据库
	db, err := database.NewDatabase()
	if err != nil {
		log.Printf("初始化数据库失败: %v", err)
		return
	}

	a.database = db
	log.Println("数据库初始化成功")

	// 检查数据库和表是否存在
	if ok, err := a.database.CheckDatabase(); err != nil {
		log.Printf("检查数据库失败: %v", err)
	} else if ok {
		log.Println("数据库和所有表检查通过")
	}
}

// shutdown 应用关闭时调用
func (a *app) shutdown(ctx context.Context) {
	if a.database != nil {
		a.database.Close()
		log.Println("数据库连接已关闭")
	}
}

// GetDatabase 获取数据库实例
func (a *app) GetDatabase() *database.Database {
	return a.database
}

// CheckDatabaseStatus 检查数据库状态（供前端调用）
func (a *app) CheckDatabaseStatus() (bool, string) {
	if a.database == nil {
		return false, "数据库未初始化"
	}

	ok, err := a.database.CheckDatabase()
	if err != nil {
		return false, fmt.Sprintf("数据库检查失败: %v", err)
	}

	return ok, "数据库状态正常"
}

// GetDatabaseInfo 获取数据库信息（供前端调用）
func (a *app) GetDatabaseInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if a.database == nil {
		info["status"] = "未初始化"
		info["path"] = ""
		return info
	}

	path, err := database.GetDatabasePath()
	if err != nil {
		info["path"] = fmt.Sprintf("获取路径失败: %v", err)
	} else {
		info["path"] = path
	}

	ok, err := a.database.CheckDatabase()
	if err != nil {
		info["status"] = fmt.Sprintf("错误: %v", err)
	} else if ok {
		info["status"] = "正常"
	} else {
		info["status"] = "异常"
	}

	return info
}

// ============ 人物相关接口 ============

// GetAllCharacters 获取所有人物列表
func (a *app) GetAllCharacters() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.database.GetAllCharacters()
}

// GetCharacterInfo 获取人物详细信息
func (a *app) GetCharacterInfo(characterID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetCharacterInfo(characterID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":         info.ID,
		"name":       info.Name,
		"shili":      info.Shili,
		"property":   info.Property,
		"level":      info.Level,
		"attributes": info.Attributes,
		"skills":     info.Skills,
	}

	return result, nil
}

// CreateCharacter 创建新人物
func (a *app) CreateCharacter(name, shili string, property, level int) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.database.CreateCharacter(name, shili, property, level)
}

// UpdateCharacterBasicInfo 更新人物基本信息
func (a *app) UpdateCharacterBasicInfo(characterID int, name, shili string, property, level int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateCharacterBasicInfo(characterID, name, shili, property, level)
}

// DeleteCharacter 删除人物
func (a *app) DeleteCharacter(characterID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteCharacter(characterID)
}

// AddCharacterAttribute 添加人物属性
func (a *app) AddCharacterAttribute(characterID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddCharacterAttribute(characterID, name, description, value)
}

// DeleteCharacterAttribute 删除人物属性
func (a *app) DeleteCharacterAttribute(attributeID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteCharacterAttribute(attributeID)
}

// UpdateCharacterAttribute 更新人物属性
func (a *app) UpdateCharacterAttribute(attributeID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateCharacterAttribute(attributeID, name, description, value)
}

// AddCharacterSkill 添加人物技能
func (a *app) AddCharacterSkill(characterID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddCharacterSkill(characterID, name, description)
}

// DeleteCharacterSkill 删除人物技能
func (a *app) DeleteCharacterSkill(skillID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteCharacterSkill(skillID)
}

// UpdateCharacterSkill 更新人物技能
func (a *app) UpdateCharacterSkill(skillID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateCharacterSkill(skillID, name, description)
}

// ============ 武器相关接口 ============

// GetAllWeapons 获取所有武器列表
func (a *app) GetAllWeapons() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.database.GetAllWeapons()
}

// GetWeaponInfo 获取武器详细信息
func (a *app) GetWeaponInfo(weaponID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetWeaponInfo(weaponID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":         info.ID,
		"name":       info.Name,
		"holder":     info.Holder,
		"level":      info.Level,
		"attributes": info.Attributes,
		"skills":     info.Skills,
	}

	return result, nil
}

// CreateWeapon 创建新武器
func (a *app) CreateWeapon(name, holder string, level int) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.database.CreateWeapon(name, holder, level)
}

// UpdateWeaponBasicInfo 更新武器基本信息
func (a *app) UpdateWeaponBasicInfo(weaponID int, name, holder string, level int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateWeaponBasicInfo(weaponID, name, holder, level)
}

// DeleteWeapon 删除武器
func (a *app) DeleteWeapon(weaponID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteWeapon(weaponID)
}

// AddWeaponAttribute 添加武器属性
func (a *app) AddWeaponAttribute(weaponID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddWeaponAttribute(weaponID, name, description, value)
}

// DeleteWeaponAttribute 删除武器属性
func (a *app) DeleteWeaponAttribute(attributeID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteWeaponAttribute(attributeID)
}

// UpdateWeaponAttribute 更新武器属性
func (a *app) UpdateWeaponAttribute(attributeID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateWeaponAttribute(attributeID, name, description, value)
}

// AddWeaponSkill 添加武器技能
func (a *app) AddWeaponSkill(weaponID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddWeaponSkill(weaponID, name, description)
}

// DeleteWeaponSkill 删除武器技能
func (a *app) DeleteWeaponSkill(skillID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteWeaponSkill(skillID)
}

// UpdateWeaponSkill 更新武器技能
func (a *app) UpdateWeaponSkill(skillID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateWeaponSkill(skillID, name, description)
}
