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

// ============ 宠物相关接口 ============

// GetAllPets 获取所有宠物列表
func (a *app) GetAllPets() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.database.GetAllPets()
}

// GetPetInfo 获取宠物详细信息
func (a *app) GetPetInfo(petID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetPetInfo(petID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":         info.ID,
		"name":       info.Name,
		"level":      info.Level,
		"owner":      info.Owner,
		"attributes": info.Attributes,
		"skills":     info.Skills,
	}

	return result, nil
}

// CreatePet 创建新宠物
func (a *app) CreatePet(name, owner string, level int) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.database.CreatePet(name, owner, level)
}

// UpdatePetBasicInfo 更新宠物基本信息
func (a *app) UpdatePetBasicInfo(petID int, owner string, level int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdatePetBasicInfo(petID, owner, level)
}

// DeletePet 删除宠物
func (a *app) DeletePet(petID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeletePet(petID)
}

// AddPetAttribute 添加宠物属性
func (a *app) AddPetAttribute(petID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddPetAttribute(petID, name, description, value)
}

// DeletePetAttribute 删除宠物属性
func (a *app) DeletePetAttribute(attributeID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeletePetAttribute(attributeID)
}

// UpdatePetAttribute 更新宠物属性
func (a *app) UpdatePetAttribute(attributeID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdatePetAttribute(attributeID, name, description, value)
}

// AddPetSkill 添加宠物技能
func (a *app) AddPetSkill(petID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddPetSkill(petID, name, description)
}

// DeletePetSkill 删除宠物技能
func (a *app) DeletePetSkill(skillID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeletePetSkill(skillID)
}

// UpdatePetSkill 更新宠物技能
func (a *app) UpdatePetSkill(skillID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdatePetSkill(skillID, name, description)
}

// ============ 道具相关接口 ============

// GetAllDaoju 获取所有道具列表
func (a *app) GetAllDaoju() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.database.GetAllDaoju()
}

// GetDaojuInfo 获取道具详细信息
func (a *app) GetDaojuInfo(daojuID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetDaojuInfo(daojuID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":        info.ID,
		"name":      info.Name,
		"level":     info.Level,
		"holder":    info.Holder,
		"functions": info.Functions,
	}

	return result, nil
}

// CreateDaoju 创建新道具
func (a *app) CreateDaoju(name, holder string, level int) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	id, err := a.database.CreateDaoju(name, level, holder)
	return int(id), err
}

// UpdateDaojuBasicInfo 更新道具基本信息
func (a *app) UpdateDaojuBasicInfo(daojuID int, holder string, level int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateDaojuBasicInfo(daojuID, level, holder)
}

// DeleteDaoju 删除道具
func (a *app) DeleteDaoju(daojuID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteDaoju(daojuID)
}

// GetDaojuFunctions 获取道具功能列表
func (a *app) GetDaojuFunctions(daojuID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	functions, err := a.database.GetDaojuFunctions(daojuID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(functions))
	for i, function := range functions {
		result[i] = map[string]interface{}{
			"id":          function.ID,
			"name":        function.Name,
			"description": function.Description,
		}
	}

	return result, nil
}

// AddDaojuFunction 添加道具功能
func (a *app) AddDaojuFunction(daojuID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddDaojuFunction(daojuID, name, description)
}

// DeleteDaojuFunction 删除道具功能
func (a *app) DeleteDaojuFunction(functionID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteDaojuFunction(functionID)
}

// UpdateDaojuFunction 更新道具功能
func (a *app) UpdateDaojuFunction(functionID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateDaojuFunction(functionID, name, description)
}

// ============ 任务相关接口 ============

// GetAllShiqing 获取所有任务列表
func (a *app) GetAllShiqing() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.database.GetAllShiqing()
}

// GetShiqingInfo 获取任务详细信息
func (a *app) GetShiqingInfo(shiqingID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetShiqingInfo(shiqingID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":       info.ID,
		"name":     info.Name,
		"location": info.Location,
		"time":     info.Time,
		"details":  info.Details,
	}

	return result, nil
}

// CreateShiqing 创建新任务
func (a *app) CreateShiqing(name, location, time string) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	id, err := a.database.CreateShiqing(name, location, time)
	return int(id), err
}

// UpdateShiqingBasicInfo 更新任务基本信息
func (a *app) UpdateShiqingBasicInfo(shiqingID int, location, time string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateShiqingBasicInfo(shiqingID, location, time)
}

// DeleteShiqing 删除任务
func (a *app) DeleteShiqing(shiqingID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteShiqing(shiqingID)
}

// GetShiqingDetails 获取任务详情列表
func (a *app) GetShiqingDetails(shiqingID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	details, err := a.database.GetShiqingDetails(shiqingID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		result[i] = map[string]interface{}{
			"id":          detail.ID,
			"description": detail.Description,
		}
	}

	return result, nil
}

// AddShiqingDetail 添加任务详情
func (a *app) AddShiqingDetail(shiqingID int, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddShiqingDetail(shiqingID, description)
}

// DeleteShiqingDetail 删除任务详情
func (a *app) DeleteShiqingDetail(detailID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteShiqingDetail(detailID)
}

// UpdateShiqingDetail 更新任务详情
func (a *app) UpdateShiqingDetail(detailID int, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateShiqingDetail(detailID, description)
}

// ============ 怪物相关接口 ============

// GetAllGuaiwu 获取所有怪物列表
func (a *app) GetAllGuaiwu() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	guaiwuList, err := a.database.GetAllGuaiwu()
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(guaiwuList))
	for i, g := range guaiwuList {
		result[i] = map[string]interface{}{
			"id":      g.ID,
			"name":    g.Name,
			"type":    g.Type,
			"level":   g.Level,
			"health":  g.Health,
			"attack":  g.Attack,
			"defense": g.Defense,
			"rewards": g.Rewards,
		}
	}

	return result, nil
}

// GetGuaiwuInfo 获取怪物详细信息
func (a *app) GetGuaiwuInfo(guaiwuID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetGuaiwuInfo(guaiwuID)
	if err != nil {
		return nil, err
	}

	// 获取属性和技能
	attributes, err := a.database.GetGuaiwuAttributes(guaiwuID)
	if err != nil {
		return nil, err
	}

	skills, err := a.database.GetGuaiwuSkills(guaiwuID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":         info.ID,
		"name":       info.Name,
		"type":       info.Type,
		"level":      info.Level,
		"health":     info.Health,
		"attack":     info.Attack,
		"defense":    info.Defense,
		"rewards":    info.Rewards,
		"attributes": attributes,
		"skills":     skills,
	}

	return result, nil
}

// CreateGuaiwu 创建新怪物
func (a *app) CreateGuaiwu(name, guaiwuType string, level, health, attack, defense int, rewards string) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.database.CreateGuaiwu(name, guaiwuType, level, health, attack, defense, rewards)
}

// UpdateGuaiwuBasicInfo 更新怪物基本信息
func (a *app) UpdateGuaiwuBasicInfo(guaiwuID int, level, health, attack, defense int, rewards string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateGuaiwuBasicInfo(guaiwuID, level, health, attack, defense, rewards)
}

// DeleteGuaiwu 删除怪物
func (a *app) DeleteGuaiwu(guaiwuID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteGuaiwu(guaiwuID)
}

// GetGuaiwuAttributes 获取怪物属性列表
func (a *app) GetGuaiwuAttributes(guaiwuID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	attributes, err := a.database.GetGuaiwuAttributes(guaiwuID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(attributes))
	for i, attr := range attributes {
		result[i] = map[string]interface{}{
			"id":          attr.ID,
			"guaiwu_id":   attr.GuaiwuID,
			"name":        attr.Name,
			"description": attr.Description,
			"value":       attr.Value,
		}
	}

	return result, nil
}

// AddGuaiwuAttribute 添加怪物属性
func (a *app) AddGuaiwuAttribute(guaiwuID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddGuaiwuAttribute(guaiwuID, name, description, value)
}

// DeleteGuaiwuAttribute 删除怪物属性
func (a *app) DeleteGuaiwuAttribute(attributeID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteGuaiwuAttribute(attributeID)
}

// UpdateGuaiwuAttribute 更新怪物属性
func (a *app) UpdateGuaiwuAttribute(attributeID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateGuaiwuAttribute(attributeID, name, description, value)
}

// GetGuaiwuSkills 获取怪物技能列表
func (a *app) GetGuaiwuSkills(guaiwuID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	skills, err := a.database.GetGuaiwuSkills(guaiwuID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(skills))
	for i, skill := range skills {
		result[i] = map[string]interface{}{
			"id":          skill.ID,
			"guaiwu_id":   skill.GuaiwuID,
			"name":        skill.Name,
			"description": skill.Description,
		}
	}

	return result, nil
}

// AddGuaiwuSkill 添加怪物技能
func (a *app) AddGuaiwuSkill(guaiwuID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddGuaiwuSkill(guaiwuID, name, description)
}

// DeleteGuaiwuSkill 删除怪物技能
func (a *app) DeleteGuaiwuSkill(skillID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteGuaiwuSkill(skillID)
}

// UpdateGuaiwuSkill 更新怪物技能
func (a *app) UpdateGuaiwuSkill(skillID int, name, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateGuaiwuSkill(skillID, name, description)
}

// ============ 势力相关接口 ============

// GetAllShili 获取所有势力列表
func (a *app) GetAllShili() ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	shiliList, err := a.database.GetAllShili()
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(shiliList))
	for i, s := range shiliList {
		result[i] = map[string]interface{}{
			"id":           s.ID,
			"name":         s.Name,
			"level":        s.Level,
			"founder":      s.Founder,
			"wealth":       s.Wealth,
			"member_count": s.MemberCount,
			"max_members":  s.MaxMembers,
		}
	}

	return result, nil
}

// GetShiliInfo 获取势力详细信息
func (a *app) GetShiliInfo(shiliID int) (map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	info, err := a.database.GetShiliInfo(shiliID)
	if err != nil {
		return nil, err
	}

	// 获取职务和属性
	positions, err := a.database.GetShiliPositions(shiliID)
	if err != nil {
		return nil, err
	}

	attributes, err := a.database.GetShiliAttributes(shiliID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := map[string]interface{}{
		"id":           info.ID,
		"name":         info.Name,
		"level":        info.Level,
		"founder":      info.Founder,
		"wealth":       info.Wealth,
		"member_count": info.MemberCount,
		"max_members":  info.MaxMembers,
		"positions":    positions,
		"attributes":   attributes,
	}

	return result, nil
}

// CreateShili 创建新势力
func (a *app) CreateShili(name, founder string, level, wealth, maxMembers int) (int, error) {
	if a.database == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	id, err := a.database.CreateShili(name, founder, level, wealth, maxMembers)
	return int(id), err
}

// UpdateShiliBasicInfo 更新势力基本信息
func (a *app) UpdateShiliBasicInfo(shiliID int, level int, founder string, wealth, memberCount, maxMembers int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateShiliBasicInfo(shiliID, level, founder, wealth, memberCount, maxMembers)
}

// DeleteShili 删除势力
func (a *app) DeleteShili(shiliID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteShili(shiliID)
}

// GetShiliPositions 获取势力职务列表
func (a *app) GetShiliPositions(shiliID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	positions, err := a.database.GetShiliPositions(shiliID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(positions))
	for i, pos := range positions {
		result[i] = map[string]interface{}{
			"id":            pos.ID,
			"shili_id":      pos.ShiliID,
			"position_name": pos.PositionName,
			"person_name":   pos.PersonName,
			"description":   pos.Description,
		}
	}

	return result, nil
}

// AddShiliPosition 添加势力职务
func (a *app) AddShiliPosition(shiliID int, positionName, personName, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddShiliPosition(shiliID, positionName, personName, description)
}

// DeleteShiliPosition 删除势力职务
func (a *app) DeleteShiliPosition(positionID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteShiliPosition(positionID)
}

// UpdateShiliPosition 更新势力职务
func (a *app) UpdateShiliPosition(positionID int, positionName, personName, description string) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateShiliPosition(positionID, positionName, personName, description)
}

// GetShiliAttributes 获取势力属性列表
func (a *app) GetShiliAttributes(shiliID int) ([]map[string]interface{}, error) {
	if a.database == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	attributes, err := a.database.GetShiliAttributes(shiliID)
	if err != nil {
		return nil, err
	}

	// 转换为 map 以便 JSON 序列化
	result := make([]map[string]interface{}, len(attributes))
	for i, attr := range attributes {
		result[i] = map[string]interface{}{
			"id":          attr.ID,
			"shili_id":    attr.ShiliID,
			"name":        attr.Name,
			"description": attr.Description,
			"value":       attr.Value,
		}
	}

	return result, nil
}

// AddShiliAttribute 添加势力属性
func (a *app) AddShiliAttribute(shiliID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.AddShiliAttribute(shiliID, name, description, value)
}

// DeleteShiliAttribute 删除势力属性
func (a *app) DeleteShiliAttribute(attributeID int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.DeleteShiliAttribute(attributeID)
}

// UpdateShiliAttribute 更新势力属性
func (a *app) UpdateShiliAttribute(attributeID int, name, description string, value int) error {
	if a.database == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return a.database.UpdateShiliAttribute(attributeID, name, description, value)
}
