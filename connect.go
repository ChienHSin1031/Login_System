// 1. id數要判別
// 3. 判斷帳密是否重複
// 4. 前端太醜

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	GetWebData()
	// insert()
}

func GetWebData() {
	http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		name := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Fprintf(w, "name = %s\n", name)
		fmt.Fprintf(w, "password = %s\n", password)

		//檢查密碼是否正確
		databasecheck(name, password)

		databaseWrite(name, password)
	}
}

//檢查密碼是否正確
func databasecheck(name string, password string) {
	db, err := sql.Open("mysql", "root:@/login")
	// var id int
	rows, err := db.Query("SELECT * FROM login.new_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var pw, count string
		if err := rows.Scan(&count, &pw); err != nil {
			log.Fatal(err)
		}
		if count == name && pw == password {
			fmt.Printf("登入成功\n")
			break
		}
		fmt.Printf("%s,%s\n", count, pw)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func databaseWrite(name string, password string) {
	//資料庫寫入
	db, err := sql.Open("mysql", "root:@/login")
	stmt, _ := db.Prepare(`INSERT INTO new_table (name, password) VALUES (?, ?)`)
	defer stmt.Close()

	ret, err := stmt.Exec(name, password)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:", RowsAffected)
	}

}

//DataBase Part
//新增
func insert() {
	db, err := sql.Open("mysql", "root:@/login")
	stmt, _ := db.Prepare(`INSERT INTO new_table (id, name, password) VALUES (?,?, ?)`)
	defer stmt.Close()
	i := 7
	ret, err := stmt.Exec(i, "123", " 451")
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return
	}

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:", RowsAffected)
	}

}

//錯誤檢查
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
