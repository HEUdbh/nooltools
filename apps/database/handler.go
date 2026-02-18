// 创建本地sqlite数据库，创建数据库表（人物、武器、任务、帮派、怪物、道具、宠物）
// 处理数据库相关操作
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Database 数据库处理器
type Database struct {
	db *sql.DB
}

// NewDatabase 创建新的数据库实例
func NewDatabase() (*Database, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户目录失败: %v", err)
	}

	// 创建应用数据目录
	appDataDir := filepath.Join(homeDir, ".nooltools")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建应用目录失败: %v", err)
	}

	// 数据库文件路径
	dbPath := filepath.Join(appDataDir, "nooltools.db")

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	database := &Database{db: db}

	// 初始化数据库表
	if err := database.initTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("初始化数据库表失败: %v", err)
	}

	return database, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB 获取数据库连接
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// initTables 初始化所有数据库表
func (d *Database) initTables() error {
	// 创建人物表
	if err := d.createRenwuTable(); err != nil {
		return err
	}

	// 更新人物表结构（处理旧版本数据库）
	if err := d.updateRenwuTableSchema(); err != nil {
		return err
	}

	// 创建人物属性表
	if err := d.createRenwuAttributesTable(); err != nil {
		return err
	}

	// 创建人物技能表
	if err := d.createRenwuSkillsTable(); err != nil {
		return err
	}

	// 创建武器表
	if err := d.createWuqiTable(); err != nil {
		return err
	}

	// 创建武器属性表
	if err := d.createWuqiAttributesTable(); err != nil {
		return err
	}

	// 创建武器技能表
	if err := d.createWuqiSkillsTable(); err != nil {
		return err
	}

	// 创建任务表
	if err := d.createShiqingTable(); err != nil {
		return err
	}

	// 创建任务详情表
	if err := d.createShiqingDetailsTable(); err != nil {
		return err
	}

	// 创建帮派表
	if err := d.createShiliTable(); err != nil {
		return err
	}

	// 更新势力表结构（处理旧版本数据库）
	if err := d.updateShiliTableSchema(); err != nil {
		return err
	}

	// 创建势力职务表
	if err := d.createShiliPositionsTable(); err != nil {
		return err
	}

	// 创建势力属性表
	if err := d.createShiliAttributesTable(); err != nil {
		return err
	}

	// 创建怪物表
	if err := d.createGuaiwuTable(); err != nil {
		return err
	}

	// 更新怪物表结构（处理旧版本数据库）
	if err := d.updateGuaiwuTableSchema(); err != nil {
		return err
	}

	// 创建怪物属性表
	if err := d.createGuaiwuAttributesTable(); err != nil {
		return err
	}

	// 创建怪物技能表
	if err := d.createGuaiwuSkillsTable(); err != nil {
		return err
	}

	// 创建道具表
	if err := d.createDaojuTable(); err != nil {
		return err
	}

	// 创建宠物表
	if err := d.createChongwuTable(); err != nil {
		return err
	}

	// 创建宠物属性表
	if err := d.createChongwuAttributesTable(); err != nil {
		return err
	}

	// 创建宠物技能表
	if err := d.createChongwuSkillsTable(); err != nil {
		return err
	}

	// 创建道具功能表
	if err := d.createDaojuFunctionsTable(); err != nil {
		return err
	}

	// 创建背包表
	if err := d.createBeibaoTable(); err != nil {
		return err
	}

	// 创建背包物品表
	if err := d.createBeibaoItemsTable(); err != nil {
		return err
	}

	// 创建商城表
	if err := d.createShoppingTable(); err != nil {
		return err
	}

	// 创建奖池表
	if err := d.createPrizesTable(); err != nil {
		return err
	}

	// 创建抽奖历史表
	if err := d.createDrawHistoryTable(); err != nil {
		return err
	}

	return nil
}

// createRenwuTable 创建人物表
func (d *Database) createRenwuTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS renwu (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT 'xx',
		shili TEXT DEFAULT 'xxx',
		level INTEGER DEFAULT 0,
		property INTEGER DEFAULT 100,
		health_current INTEGER DEFAULT 100,
		mana_current INTEGER DEFAULT 100,
		attack INTEGER DEFAULT 100,
		defense INTEGER DEFAULT 100,
		skills TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// updateRenwuTableSchema 更新人物表结构（处理旧版本数据库）
func (d *Database) updateRenwuTableSchema() error {
	// 检查并添加 shili 字段
	if err := d.addColumnIfNotExists("renwu", "shili", "TEXT DEFAULT 'xxx'"); err != nil {
		return err
	}

	// 检查并添加 property 字段
	if err := d.addColumnIfNotExists("renwu", "property", "INTEGER DEFAULT 100"); err != nil {
		return err
	}

	// 更新 level 默认值为 0（如果需要）
	// 注意：SQLite 不支持直接修改列的默认值，这里只是示例

	// 更新 attack 默认值为 100（如果需要）
	// 注意：SQLite 不支持直接修改列的默认值，这里只是示例

	// 更新 defense 默认值为 100（如果需要）
	// 注意：SQLite 不支持直接修改列的默认值，这里只是示例

	return nil
}

// addColumnIfNotExists 如果列不存在则添加列
func (d *Database) addColumnIfNotExists(tableName, columnName, columnDef string) error {
	// 检查列是否存在
	query := `PRAGMA table_info(` + tableName + `)`
	rows, err := d.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columnExists := false
	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var dfltValue interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}

		if name == columnName {
			columnExists = true
			break
		}
	}

	// 如果列不存在，则添加列
	if !columnExists {
		alterQuery := `ALTER TABLE ` + tableName + ` ADD COLUMN ` + columnName + ` ` + columnDef
		_, err := d.db.Exec(alterQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

// createWuqiTable 创建武器表
func (d *Database) createWuqiTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS wuqi (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		holder TEXT DEFAULT '',
		level INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createWuqiAttributesTable 创建武器属性表
func (d *Database) createWuqiAttributesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS wuqi_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wuqi_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (wuqi_id) REFERENCES wuqi(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createWuqiSkillsTable 创建武器技能表
func (d *Database) createWuqiSkillsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS wuqi_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wuqi_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (wuqi_id) REFERENCES wuqi(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createShiqingTable 创建任务表
func (d *Database) createShiqingTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shiqing (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		details TEXT,
		time TEXT,
		rewards TEXT,
		location TEXT,
		participants TEXT,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createShiliTable 创建帮派表
func (d *Database) createShiliTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shili (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		level INTEGER DEFAULT 1,
		founder TEXT NOT NULL DEFAULT '',
		wealth INTEGER DEFAULT 0,
		member_count INTEGER DEFAULT 0,
		max_members INTEGER DEFAULT 10,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// updateShiliTableSchema 更新势力表结构（处理旧版本数据库）
func (d *Database) updateShiliTableSchema() error {
	// 检查并添加 level 字段
	if err := d.addColumnIfNotExists("shili", "level", "INTEGER DEFAULT 1"); err != nil {
		return err
	}

	// 检查并添加 founder 字段
	if err := d.addColumnIfNotExists("shili", "founder", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}

	// 检查并添加 wealth 字段
	if err := d.addColumnIfNotExists("shili", "wealth", "INTEGER DEFAULT 0"); err != nil {
		return err
	}

	// 检查并添加 member_count 字段
	if err := d.addColumnIfNotExists("shili", "member_count", "INTEGER DEFAULT 0"); err != nil {
		return err
	}

	// 检查并添加 max_members 字段
	if err := d.addColumnIfNotExists("shili", "max_members", "INTEGER DEFAULT 10"); err != nil {
		return err
	}

	return nil
}

// createShiliPositionsTable 创建势力职务表
func (d *Database) createShiliPositionsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shili_positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shili_id INTEGER NOT NULL,
		position_name TEXT NOT NULL,
		person_name TEXT NOT NULL DEFAULT '',
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shili_id) REFERENCES shili(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createShiliAttributesTable 创建势力属性表
func (d *Database) createShiliAttributesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shili_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shili_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shili_id) REFERENCES shili(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createGuaiwuTable 创建怪物表
func (d *Database) createGuaiwuTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS guaiwu (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT '未命名',
		type TEXT NOT NULL,
		level INTEGER DEFAULT 1,
		health INTEGER DEFAULT 100,
		attack INTEGER DEFAULT 10,
		defense INTEGER DEFAULT 10,
		rewards TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// updateGuaiwuTableSchema 更新怪物表结构（处理旧版本数据库）
func (d *Database) updateGuaiwuTableSchema() error {
	// 检查并添加 name 字段
	if err := d.addColumnIfNotExists("guaiwu", "name", "TEXT NOT NULL DEFAULT '未命名'"); err != nil {
		return err
	}

	return nil
}

// createDaojuTable 创建道具表
func (d *Database) createDaojuTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS daoju (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		level INTEGER DEFAULT 1,
		function TEXT,
		durability INTEGER DEFAULT 100,
		holder TEXT DEFAULT '',
		price INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createChongwuTable 创建宠物表
func (d *Database) createChongwuTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS chongwu (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		level INTEGER DEFAULT 1,
		attributes TEXT,
		owner TEXT DEFAULT '',
		price INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// CheckDatabase 检查数据库和表是否存在
func (d *Database) CheckDatabase() (bool, error) {
	// 检查数据库连接
	if d.db == nil {
		return false, fmt.Errorf("数据库连接为空")
	}

	// 检查所有表是否存在
	tables := []string{"renwu", "renwu_attributes", "renwu_skills", "wuqi", "wuqi_attributes", "wuqi_skills", "shiqing", "shiqing_details", "shili", "shili_positions", "shili_attributes", "guaiwu", "guaiwu_attributes", "guaiwu_skills", "daoju", "chongwu", "chongwu_attributes", "chongwu_skills", "beibao", "beibao_items"}

	for _, table := range tables {
		if err := d.checkTableExists(table); err != nil {
			return false, err
		}
	}

	return true, nil
}

// checkTableExists 检查表是否存在
func (d *Database) checkTableExists(tableName string) error {
	query := `
	SELECT name FROM sqlite_master 
	WHERE type='table' AND name=?`

	var name string
	err := d.db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("表 %s 不存在", tableName)
		}
		return fmt.Errorf("检查表 %s 失败: %v", tableName, err)
	}

	return nil
}

// GetDatabasePath 获取数据库文件路径
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".nooltools", "nooltools.db")
	return dbPath, nil
}

// createRenwuAttributesTable 创建人物属性表
func (d *Database) createRenwuAttributesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS renwu_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		renwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (renwu_id) REFERENCES renwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createRenwuSkillsTable 创建人物技能表
func (d *Database) createRenwuSkillsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS renwu_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		renwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (renwu_id) REFERENCES renwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createChongwuAttributesTable 创建宠物属性表
func (d *Database) createChongwuAttributesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS chongwu_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chongwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (chongwu_id) REFERENCES chongwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createChongwuSkillsTable 创建宠物技能表
func (d *Database) createChongwuSkillsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS chongwu_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chongwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (chongwu_id) REFERENCES chongwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createDaojuFunctionsTable 创建道具功能表
func (d *Database) createDaojuFunctionsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS daoju_functions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		daoju_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (daoju_id) REFERENCES daoju(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createShiqingDetailsTable 创建任务详情表
func (d *Database) createShiqingDetailsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shiqing_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shiqing_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shiqing_id) REFERENCES shiqing(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createGuaiwuAttributesTable 创建怪物属性表
func (d *Database) createGuaiwuAttributesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS guaiwu_attributes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guaiwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (guaiwu_id) REFERENCES guaiwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createGuaiwuSkillsTable 创建怪物技能表
func (d *Database) createGuaiwuSkillsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS guaiwu_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guaiwu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (guaiwu_id) REFERENCES guaiwu(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createBeibaoTable 创建背包表
func (d *Database) createBeibaoTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS beibao (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT '默认背包',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createBeibaoItemsTable 创建背包物品表
func (d *Database) createBeibaoItemsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS beibao_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		beibao_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		quantity INTEGER DEFAULT 1,
		description TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (beibao_id) REFERENCES beibao(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}

// createShoppingTable 创建商城表
func (d *Database) createShoppingTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS shopping (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		value INTEGER NOT NULL,
		description TEXT DEFAULT '',
		condition TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createPrizesTable 创建奖池表
func (d *Database) createPrizesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS prizes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		rate REAL NOT NULL,
		description TEXT DEFAULT '',
		variety TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createDrawHistoryTable 创建抽奖历史表
func (d *Database) createDrawHistoryTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS draw_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prize_id INTEGER NOT NULL,
		prize_name TEXT NOT NULL,
		drawn_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (prize_id) REFERENCES prizes(id) ON DELETE CASCADE
	)`

	_, err := d.db.Exec(query)
	return err
}
