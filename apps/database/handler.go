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

	// 创建帮派表
	if err := d.createShiliTable(); err != nil {
		return err
	}

	// 创建怪物表
	if err := d.createGuaiwuTable(); err != nil {
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
		position TEXT,
		member_count INTEGER DEFAULT 0,
		attributes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	return err
}

// createGuaiwuTable 创建怪物表
func (d *Database) createGuaiwuTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS guaiwu (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
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
	tables := []string{"renwu", "renwu_attributes", "renwu_skills", "wuqi", "wuqi_attributes", "wuqi_skills", "shiqing", "shili", "guaiwu", "daoju", "chongwu"}

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
