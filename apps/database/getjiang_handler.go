package database

import (
	"fmt"
	"math/rand"
	"time"
)

// PrizeInfo 奖品信息结构体
type PrizeInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Rate        float64 `json:"rate"` // 爆率（百分比，如 1.5 表示 1.5%）
	Description string  `json:"description"`
	Variety     string  `json:"variety"` // 品种
}

// DrawResult 抽奖结果
type DrawResult struct {
	PrizeID      int    `json:"prize_id"`
	PrizeName    string `json:"prize_name"`
	PrizeVariety string `json:"prize_variety"`
	Description  string `json:"description"`
	DrawnAt      string `json:"drawn_at"`
}

// GetAllPrizes 获取所有奖品
func (db *Database) GetAllPrizes() ([]PrizeInfo, error) {
	query := `SELECT id, name, rate, description, variety FROM prizes ORDER BY rate DESC`
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []PrizeInfo
	for rows.Next() {
		var prize PrizeInfo
		err := rows.Scan(&prize.ID, &prize.Name, &prize.Rate, &prize.Description, &prize.Variety)
		if err != nil {
			return nil, err
		}
		prizes = append(prizes, prize)
	}

	return prizes, nil
}

// GetPrizeInfo 获取指定奖品的详细信息
func (db *Database) GetPrizeInfo(id int) (PrizeInfo, error) {
	query := `SELECT id, name, rate, description, variety FROM prizes WHERE id = ?`
	var prize PrizeInfo
	err := db.db.QueryRow(query, id).Scan(&prize.ID, &prize.Name, &prize.Rate, &prize.Description, &prize.Variety)
	if err != nil {
		return PrizeInfo{}, err
	}
	return prize, nil
}

// CreatePrize 创建新奖品
func (db *Database) CreatePrize(name string, rate float64, description, prizeVariety string) (int64, error) {
	query := `INSERT INTO prizes (name, rate, description, variety) VALUES (?, ?, ?, ?)`
	result, err := db.db.Exec(query, name, rate, description, prizeVariety)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// DeletePrize 删除奖品
func (db *Database) DeletePrize(id int) error {
	query := `DELETE FROM prizes WHERE id = ?`
	_, err := db.db.Exec(query, id)
	return err
}

// UpdatePrize 更新奖品信息
func (db *Database) UpdatePrize(id int, name string, rate float64, description, prizeVariety string) error {
	query := `UPDATE prizes SET name = ?, rate = ?, description = ?, variety = ? WHERE id = ?`
	_, err := db.db.Exec(query, name, rate, description, prizeVariety, id)
	return err
}

// DrawPrize 单抽
func (db *Database) DrawPrize() (DrawResult, error) {
	// 获取所有奖品
	prizes, err := db.GetAllPrizes()
	if err != nil {
		return DrawResult{}, err
	}

	if len(prizes) == 0 {
		return DrawResult{}, fmt.Errorf("奖池中没有奖品")
	}

	// 基于爆率进行抽奖
	prize, err := drawByRate(prizes)
	if err != nil {
		return DrawResult{}, err
	}

	// 记录抽奖结果
	result := DrawResult{
		PrizeID:      prize.ID,
		PrizeName:    prize.Name,
		PrizeVariety: prize.Variety,
		Description:  prize.Description,
		DrawnAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	// 将抽奖结果保存到数据库
	insertQuery := `INSERT INTO draw_history (prize_id, prize_name, drawn_at) VALUES (?, ?, ?)`
	_, err = db.db.Exec(insertQuery, result.PrizeID, result.PrizeName, result.DrawnAt)
	if err != nil {
		return DrawResult{}, fmt.Errorf("保存抽奖历史失败: %v", err)
	}

	return result, nil
}

// DrawTenPrizes 十连抽
func (db *Database) DrawTenPrizes() ([]DrawResult, error) {
	var results []DrawResult

	for i := 0; i < 10; i++ {
		result, err := db.DrawPrize()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// drawByRate 根据爆率进行抽奖
func drawByRate(prizes []PrizeInfo) (PrizeInfo, error) {
	if len(prizes) == 0 {
		return PrizeInfo{}, fmt.Errorf("奖品列表为空")
	}

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 计算总爆率
	totalRate := 0.0
	for _, prize := range prizes {
		totalRate += prize.Rate
	}

	// 如果总爆率不是 100，按比例调整
	if totalRate != 100.0 {
		// 生成 0-100 的随机数
		randomNum := rand.Float64() * 100

		// 按比例计算中奖
		accumulatedRate := 0.0
		for _, prize := range prizes {
			accumulatedRate += prize.Rate
			if randomNum <= accumulatedRate {
				return prize, nil
			}
		}

		// 如果由于浮点数精度问题没有匹配到，返回最后一个奖品
		return prizes[len(prizes)-1], nil
	}

	// 如果总爆率正好是 100，直接抽奖
	randomNum := rand.Float64() * 100
	accumulatedRate := 0.0

	for _, prize := range prizes {
		accumulatedRate += prize.Rate
		if randomNum <= accumulatedRate {
			return prize, nil
		}
	}

	// 如果由于浮点数精度问题没有匹配到，返回最后一个奖品
	return prizes[len(prizes)-1], nil
}

// GetDrawHistory 获取抽奖历史记录
func (db *Database) GetDrawHistory(limit int) ([]DrawResult, error) {
	query := `SELECT prize_id, prize_name, drawn_at FROM draw_history ORDER BY drawn_at DESC LIMIT ?`
	rows, err := db.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DrawResult
	for rows.Next() {
		var result DrawResult
		err := rows.Scan(&result.PrizeID, &result.PrizeName, &result.DrawnAt)
		if err != nil {
			return nil, err
		}
		history = append(history, result)
	}

	return history, nil
}

// ClearDrawHistory 清空抽奖历史
func (db *Database) ClearDrawHistory() error {
	query := `DELETE FROM draw_history`
	_, err := db.db.Exec(query)
	return err
}
