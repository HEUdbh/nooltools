package database

// ShoppingInfo 商城商品信息结构体
type ShoppingInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
}

// GetAllShopping 获取所有商城商品
func (db *Database) GetAllShopping() ([]ShoppingInfo, error) {
	query := `SELECT id, name, value, description, condition FROM shopping ORDER BY id`
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ShoppingInfo
	for rows.Next() {
		var item ShoppingInfo
		err := rows.Scan(&item.ID, &item.Name, &item.Value, &item.Description, &item.Condition)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetShoppingInfo 获取指定商城商品的详细信息
func (db *Database) GetShoppingInfo(id int) (ShoppingInfo, error) {
	query := `SELECT id, name, value, description, condition FROM shopping WHERE id = ?`
	var item ShoppingInfo
	err := db.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Value, &item.Description, &item.Condition)
	if err != nil {
		return ShoppingInfo{}, err
	}
	return item, nil
}

// CreateShopping 创建新的商城商品
func (db *Database) CreateShopping(name string, value int, description, condition string) (int64, error) {
	query := `INSERT INTO shopping (name, value, description, condition) VALUES (?, ?, ?, ?)`
	result, err := db.db.Exec(query, name, value, description, condition)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// DeleteShopping 删除商城商品
func (db *Database) DeleteShopping(id int) error {
	query := `DELETE FROM shopping WHERE id = ?`
	_, err := db.db.Exec(query, id)
	return err
}

// UpdateShopping 更新商城商品信息
func (db *Database) UpdateShopping(id int, name string, value int, description, condition string) error {
	query := `UPDATE shopping SET name = ?, value = ?, description = ?, condition = ? WHERE id = ?`
	_, err := db.db.Exec(query, name, value, description, condition, id)
	return err
}
