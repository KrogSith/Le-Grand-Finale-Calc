package application

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type (
	User struct {
		ID       string
		Name     string
		Password string
	}

	Expression struct {
		ID         string  `json:"id"`
		Expression string  `json:"expression"`
		UserID     string  `json:"userId"`
		Status     string  `json:"status"`
		Result     float64 `json:"result"`
	}
)

func NewUser(name string, password string) User {
	return User{
		ID:       uuid.NewString(),
		Name:     name,
		Password: password,
	}
}

func (u User) Print() string {
	return "ID: " + u.ID + " Name: " + u.Name + " Password: " + u.Password
}

func NewExpression(expression string, userId string, status string, result float64) Expression {
	return Expression{
		ID:         uuid.NewString(),
		Expression: expression,
		UserID:     userId,
		Status:     status,
		Result:     result,
	}
}

func (e Expression) Print() string {
	result := strconv.FormatFloat(e.Result, 'f', 10, 64)
	return "ID: " + e.ID + " Expression" + e.Expression + " UserID:" + e.UserID + " Status:" + e.Status + " Result:" + result
}

func insertUser(ctx context.Context, db *sql.DB, user User) error {
	var q = `
	INSERT INTO users (id, name, password) values ($1, $2, $3)
	`
	_, err := db.ExecContext(ctx, q, user.ID, user.Name, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func insertExpression(ctx context.Context, db *sql.DB, expression Expression) error {
	var q = `
	INSERT INTO expressions (id, expression, user_id, status, result) values ($1, $2, $3, $4, $5)
	`
	_, err := db.ExecContext(ctx, q, expression.ID, expression.Expression, expression.UserID, expression.Status, expression.Result)
	if err != nil {
		return err
	}

	return nil
}

func selectUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	var users []User
	var q = "SELECT id, name, password FROM users"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func selectExpressions(ctx context.Context, db *sql.DB) ([]Expression, error) {
	var expressions []Expression
	var q = "SELECT id, expression, user_id, status, result FROM expressions"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Status, &e.Result)
		if err != nil {
			return nil, err
		}
		if e.UserID == current_user {
			expressions = append(expressions, e)
		}
	}

	return expressions, nil
}

func selectUserByID(ctx context.Context, db *sql.DB, id string) (User, error) {
	u := User{}
	var q = "SELECT id, name, password FROM users WHERE id = $1"
	err := db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT NOT NULL, 
		name TEXT NOT NULL,
		password TEXT NOT NULL
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id TEXT NOT NULL, 
		expression TEXT,
		user_id TEXT NOT NULL,
		status TEXT NOT NULL,
		result INT,
	
		FOREIGN KEY (user_id)  REFERENCES expressions (id)
	);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}

func InitiateDatabase() (context.Context, *sql.DB, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	if err = createTables(ctx, db); err != nil {
		return nil, nil, err
	}

	return ctx, db, nil

	// user := &User{
	// 	Name:     "Petr",
	// 	Password: "1234",
	// }

	// userID, err := insertUser(ctx, db, user)
	// if err != nil {
	// 	panic(err)
	// }

	// expression := &Expression{
	// 	Expression: "2+2",
	// 	UserID:     userID,
	// }

	// expressionID, err := insertExpression(ctx, db, expression)
	// if err != nil {
	// 	panic(err)
	// }

	// expression.ID = expressionID

	// users, err := selectUsers(ctx, db)
	// if err != nil {
	// 	panic(err)
	// }

	// for i := range users {
	// 	log.Println(users[i].Print())
	// }

	// expressions, err := selectExpressions(ctx, db)
	// if err != nil {
	// 	panic(err)
	// }

	// for i := range expressions {
	// 	log.Println(expressions[i].Print())
	// }

	// u, err := selectUserByID(ctx, db, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(u.Print())
}

func UserIsAuthorised() bool {
	return current_user != ""
}
