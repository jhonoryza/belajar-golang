package belajar_database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"testing"
)

func TestKoneksi(t *testing.T) {
	// membuka database pooling
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/belajar_go")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestInsert(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "INSERT INTO customer (name) VALUES (?)"
	param := "Budi"
	_, err := db.ExecContext(ctx, query, param)
	if err != nil {
		panic(err)
	}

	fmt.Println("Insert berhasil")
}

func TestQuery(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println(id, name)
	}
}

func TestSelect(t *testing.T) {
	db := GetConnection()

	defer db.Close()
	ctx := context.Background()
	query := "select id, name, email, balance, rating, birth_date, married, created_at from customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var email sql.NullString
		var balance int
		var rating float64
		var birthDate sql.NullTime
		var married bool
		var createdAt sql.NullTime
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("id %d \n", id)
		fmt.Printf("name %s \n", name)
		if email.Valid {
			val, _ := email.Value()
			fmt.Printf("email %v \n", val)
		}
		fmt.Printf("balance %d \n", balance)
		fmt.Printf("rating %f \n", rating)
		if birthDate.Valid {
			val, _ := birthDate.Value()
			fmt.Printf("birthdate %v \n", val)
		}
		fmt.Printf("married %t \n", married)
		if createdAt.Valid {
			val, _ := createdAt.Value()
			fmt.Printf("created_at %v \n", val)
		}
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "SELECT username, password FROM user WHERE username = ? LIMIT 1"
	param := "admin'; SELECT * FROM user --"
	rows, err := db.QueryContext(ctx, query, param)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		var email string
		err := rows.Scan(&username, &email)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("username: ", username)
		fmt.Println("email: ", email)
	} else {
		fmt.Println("user not found")
	}
}

func TestComment(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	email := "jgQpQ@example.com"
	comment := "test comment"

	result, err := db.ExecContext(ctx, query, email, comment)
	if err != nil {
		panic(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("insertId: ", insertId)
}

func TestManyInsert(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "test comment" + strconv.Itoa(i)
		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}
		insertId, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("insertId: ", insertId)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// do transaction

	statement, err := tx.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	isError := false

	for i := 0; i < 10; i++ {
		email := "eko@gmail.com"
		//email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "test comment" + strconv.Itoa(i)
		result, err2 := statement.ExecContext(ctx, email, comment)
		if err2 != nil {
			log.Println(err2)
			isError = true
			continue
		}
		insertId, err2 := result.LastInsertId()
		if err2 != nil {
			log.Println(err2)
			isError = true
			continue
		}
		fmt.Println("insertId: ", insertId)
	}

	if isError == true {
		// do rollback
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("rollback")
	} else {
		// do commit
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("commit")
	}

	result, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM comments")
	if err != nil {
		panic(err)
	}
	defer result.Close()

	if result.Next() {
		var count int
		err = result.Scan(&count)
		if err != nil {
			panic(err)
		}
		fmt.Println("count: ", count)
	}
}
