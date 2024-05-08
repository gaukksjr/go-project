package repository

import (
	"database/sql"
	"fmt"
	"go/project/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) Register(username, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	println(username, password, role, string(hashedPassword))

	_, err = r.Db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", username, string(hashedPassword), role)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Login(username, password string) (int, error) {
	var userID int
	var storedPassword string
	log.Println(username, password)
	err := r.Db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &storedPassword)
	if err != nil {
		fmt.Println("Here error")
		return 0, err
	}
	fmt.Println(userID, storedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Repository) GetMenu() ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	rows, err := r.Db.Query(`SELECT id, name, price, description, available FROM menu_items`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menuItem models.MenuItem
		err := rows.Scan(&menuItem.ID, &menuItem.Name, &menuItem.Price, &menuItem.Description, &menuItem.Available)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (r *Repository) AddMenuItems(items []models.MenuItem) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO menu_items (name, price, description, available) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(item.Name, item.Price, item.Description, item.Available)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (r *Repository) SaveOrder(userID int, items []models.OrderItem) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO orders (user_id, item_id, quantity, total_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(userID, item.ItemID, item.Quantity, item.TotalPrice)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (r *Repository) UpdateOrderStatus(orderTime time.Time) error {
	query := `
        UPDATE orders
        SET status = 'ready'
        WHERE order_time <= ?
        AND status = 'processing'
    `

	_, err := r.Db.Exec(query, orderTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRoleByUserID(userID int) (string, error) {
	var role string

	query := "SELECT role FROM users WHERE id = ?"

	err := r.Db.QueryRow(query, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %d not found", userID)
		}
		return "", err
	}

	return role, nil
}

func (r *Repository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	query := "SELECT id, user_id, item_id, quantity, total_price, status, order_time FROM orders"

	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ItemID, &order.Quantity, &order.TotalPrice, &order.Status, &order.OrderTime); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetUsernameByID(userID int) (string, error) {
	var username string

	query := "SELECT username FROM users WHERE id = ?"

	err := r.Db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}
