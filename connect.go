// Bug                             				               ??
// 1.帳號密碼 未限制字元
// 1.函式無法進入login						  finish2020/3/10
// 2.函式無法進入register					finish2020/3/10
//3.前端整理
//4. 連接
//5.帳號註冊											finish2020/3/11
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
}

func GetWebData() {
	http.HandleFunc("/", sayhelloName)       // setting router rule
	http.HandleFunc("/login", login)         //登入頁面
	http.HandleFunc("/register", register)   //註冊頁面
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//註冊
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("register.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		Registername := r.FormValue("Registername")
		Registerpassword := r.FormValue("Registerpassword")
		Checkpassword := r.FormValue("Checkpassword")

		// CheckCountexist(Registername)    檢查帳號是否已註冊 exist 0 未註冊過 : 1 註冊過
		if CheckCountexist(Registername) == 0 && Registerpassword == Checkpassword { //如果帳號沒註冊過及密碼輸入2次都相同
			databaseWrite(Registername, Registerpassword) //寫入資料庫
			fmt.Fprintln(w, "註冊成功")
		}

	}
}

// 登入
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
		if databasecheck(name, password) == 0 {
			fmt.Fprintln(w, "請輸入正確帳號密碼")
		} else {
			fmt.Fprintln(w, "登入成功")
		}

	}
}

//檢查帳號是否已存在
func CheckCountexist(Registername string) int {
	exist := 0 //判斷帳號是否已存在  0 不存在 1存在
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
		if count == Registername {
			exist = 1
			break
		}
	}
	if exist == 1 {
		fmt.Println("帳號以被使用！")
	} else {
		fmt.Println("註冊成功")
	}
	return exist
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

//檢查密碼是否正確
func databasecheck(name string, password string) int {

	db, err := sql.Open("mysql", "root:@/login")
	// var id int
	rows, err := db.Query("SELECT * FROM login.new_table")
	if err != nil {
		log.Fatal(err)
	}
	loginstate := 0 //初始化 0 未登入   ; 1 登入
	defer rows.Close()
	for rows.Next() {
		var pw, count string
		if err := rows.Scan(&count, &pw); err != nil {
			log.Fatal(err)
		}

		if count == name && pw == password { //判斷密碼是否正確
			fmt.Printf("登入成功\n")
			loginstate = 1
			break
		}
	}

	if loginstate == 0 {
		fmt.Printf("帳號密碼錯誤")
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return loginstate
}

func databaseWrite(Registername string, Registerpassword string) {
	//資料庫寫入
	db, err := sql.Open("mysql", "root:@/login")
	stmt, _ := db.Prepare(`INSERT INTO new_table (name, password) VALUES (?, ?)`)
	defer stmt.Close()

	ret, err := stmt.Exec(Registername, Registerpassword)
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
