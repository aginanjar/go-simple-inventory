package item

import (
	"database/sql"
	"fmt"
	"time"
)

type Item struct {
	ID           int       `json:"id"`
	Sku          string    `json:"sku"`
	Name         string    `json:"name"`
	SalePrice    float32   `json:"sale_price"`
	CurrentStock int       `json:"current_stock"`
	CreatedDate  time.Time `json:"created_date"`
	UpdatedDate  time.Time `json:"updated_date"`
}

func (i *Item) GetItem(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM item WHERE sku=%s", i.Sku)
	return db.QueryRow(statement).Scan(&i.Sku, &i.Name)
}

func (i *Item) UpdateItem(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE item SET name='%s', sku=%s WHERE id=%d", i.Name, i.Sku, i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *Item) DeleteItem(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM item WHERE id=%d", i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *Item) CreateItem(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO item(name, sku) VALUES('%s', %s)", i.Name, i.Sku)
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

func GetItems(db *sql.DB, start, count int) ([]Item, error) {

	statement := fmt.Sprintf("SELECT * FROM item LIMIT %d OFFSET %d", count, start)
	rows, _ := db.Query(statement)

	defer rows.Close()

	items := []Item{}

	fmt.Println(rows.Next())
	for rows.Next() {
		var i Item
		err := rows.Scan(&i.ID, &i.Sku, &i.Name, &i.SalePrice, &i.CurrentStock, &i.CreatedDate, &i.UpdatedDate)

		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}
