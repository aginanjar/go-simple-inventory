package stockindetail

import (
	"database/sql"
	"fmt"
	"time"
)

type StockInDetail struct {
	ID            int       `json:"id"`
	NoKwitansi    string    `json:"no_kwitansi"`
	AcceptedStock string    `json:"accepted_stock"`
	Date          time.Time `json:"date"`
}

func (i *StockInDetail) GetItem(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM stock_in_detail WHERE no_kwitansi=%s", i.NoKwitansi)
	return db.QueryRow(statement).Scan(&i.NoKwitansi, &i.AcceptedStock, &i.Date)
}

func (i *StockInDetail) UpdateItem(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE stock_in_detail SET accepted_stock=%d WHERE no_kwitansi=%d", i.AcceptedStock, i.NoKwitansi)
	_, err := db.Exec(statement)
	return err
}

func (i *StockInDetail) DeleteItem(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM stock_in_detail WHERE no_kwitansi=%d", i.NoKwitansi)
	_, err := db.Exec(statement)
	return err
}

func (i *StockInDetail) CreateItem(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO stock_in_detail(no_kwitansi, accepted_stock, date) VALUES('%s',%d,'%s')", i.NoKwitansi, i.AcceptedStock, i.Date)
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

func GetItems(db *sql.DB, start, count int) ([]StockInDetail, error) {

	statement := fmt.Sprintf("SELECT * FROM stock_in_detail LIMIT %d OFFSET %d", count, start)
	rows, _ := db.Query(statement)

	defer rows.Close()

	items := []StockInDetail{}

	fmt.Println(rows.Next())
	for rows.Next() {
		var i StockInDetail
		err := rows.Scan(&i.ID, &i.NoKwitansi, &i.AcceptedStock, &i.Date)

		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}
