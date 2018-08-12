package stockin

import (
	"database/sql"
	"fmt"
	"time"
)

type StockIn struct {
	ID            int       `json:"id"`
	NoKwitansi    string    `json:"no_kwitansi"`
	SKU           string    `json:"sku"`
	StockIn       float32   `json:"stock_in"`
	AcceptedStock int       `json:"accepted_stock"`
	PurchasePrice float32   `json:"purchase_price"`
	Total         float32   `json:"total"`
	Status        int       `json:"status"`
	CreatedDate   time.Time `json:"created_date"`
	UpdatedDate   time.Time `json:"updated_date"`
}

func (i *StockIn) GetStockIn(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM stock_in WHERE sku=%s", i.SKU)
	return db.QueryRow(statement).Scan(&i.SKU)
}

func (i *StockIn) UpdateStockIn(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE stock_in SET stock_in=%d, updated_date='%s' WHERE id=%d", i.StockIn, i.UpdatedDate, i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *StockIn) DeleteItem(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM stock_in WHERE id=%d", i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *StockIn) CreateItem(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO stock_in(no_kwitansi,created_date,updated_date,sku,stock_in,,accepted_stock,purchase_price,total,status) VALUES('%s', '%s', %d, %f, %f, %d, %s, %s)",
		i.NoKwitansi, i.SKU, i.AcceptedStock, i.PurchasePrice, i.Total, i.Status, i.CreatedDate, i.UpdatedDate)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetItems(db *sql.DB, start, count int) ([]StockIn, error) {

	statement := fmt.Sprintf("SELECT * FROM stock_in LIMIT %d OFFSET %d", count, start)
	rows, _ := db.Query(statement)

	defer rows.Close()

	stockins := []StockIn{}

	fmt.Println(rows.Next())
	for rows.Next() {
		var i StockIn
		err := rows.Scan(&i.ID, &i.NoKwitansi, &i.SKU, &i.AcceptedStock, &i.PurchasePrice, &i.Total, &i.Status, &i.CreatedDate, &i.UpdatedDate)

		if err != nil {
			return nil, err
		}

		stockins = append(stockins, i)
	}

	return stockins, nil
}
