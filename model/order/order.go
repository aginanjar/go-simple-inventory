package order

import (
	"database/sql"
	"fmt"
	"time"
)

type Order struct {
	ID          int       `json:"id"`
	Sku         string    `json:"sku"`
	Name        string    `json:"name"`
	Qty         int       `json:"qty"`
	CreatedDate time.Time `json:"created_date"`
}

func (i *Order) GetItem(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM order WHERE sku=%s", i.Sku)
	return db.QueryRow(statement).Scan(&i.Sku, &i.Name)
}

func (i *Order) UpdateItem(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE order SET name='%s', sku=%s WHERE id=%d", i.Name, i.Sku, i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *Order) DeleteItem(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM order WHERE id=%d", i.ID)
	_, err := db.Exec(statement)
	return err
}

func (i *Order) CreateItem(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO order(sku,name,qty,created_date) VALUES('%s', '%s', %d, '%s')", i.Sku, i.Name, i.Qty, i.CreatedDate)
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

	statement := fmt.Sprintf("SELECT * FROM order LIMIT %d OFFSET %d", count, start)
	rows, _ := db.Query(statement)

	defer rows.Close()

	orders := []Order{}

	fmt.Println(rows.Next())
	for rows.Next() {
		var i order
		err := rows.Scan(&i.ID, &i.Sku, &i.Name, &i.Qty, &i.CreatedDate)

		if err != nil {
			return nil, err
		}

		orders = append(orders, i)
	}

	return orders, nil
}
