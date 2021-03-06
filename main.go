package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"database/sql"
	//"fmt"
	"log"
	//"net/http"
	//"strconv"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
var err error 
type User struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
}
type NumOne struct{
	Num1 float64 `json:"num1"`
}
//init database..
func initDB(){
	db, err = sql.Open("mysql", "erchizhang:123456@tcp(127.0.0.1)/Trail1")
	if err != nil {
		fmt.Println("righthere")
	}
}
func login(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
		fmt.Println(decoder)
		var(
			numsData NumOne
		)
		decoder.Decode(&numsData)
		fmt.Println(numsData)
		numsData.Num1++
		
		data, _ := json.Marshal(numsData.Num1)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers","*")
		w.Write(data)
}
func saveUser(w http.ResponseWriter, r *http.Request){
	decoder:= json.NewDecoder(r.Body)
	//fmt.Println(decoder)
	if err!=nil{
		fmt.Println(err)
	}
	var(
		user User
	)
	decoder.Decode(&user)
	fmt.Println(user)
	message := saveIn(user.Username, user.Pwd)
	data,_:=json.Marshal(message)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers","*")
	w.Write(data)
	fmt.Println(user)
}
func queryName(name string)int{
	var timeCount int
	sqlStr:= "select count from user where username=?;"
	err = db.QueryRow(sqlStr, name).Scan(&timeCount)
	switch{
	case err==sql.ErrNoRows:
		return -1
	case err != nil:
		return -10
	default:
		return timeCount
	}
}
func saveIn(name string, pwd string) string {
	str:=""
	timesR:=queryName(name)
	if timesR == -1 {//no same name found
	if (name=="")||(pwd==""){
		str="username and password cannot be empty"
	}else if (strings.ContainsAny(name,"%<>/\\")||strings.ContainsAny(pwd,"%<>/\\")){
		str="username and password caannot contain %<>/\\"
	}else if (len(name)>20)||(len(pwd)>20){
		str= "username and password cannot exceed 20 characters long"
	}else{
		insert, err := db.Exec("insert into User(username, pwd, count) values(?, ?, 1)",name, pwd)
		//insert, err:= db.Exec("insert into finalTable(username, times) values("+"hehe"+",1)")
		if err != nil {
			fmt.Println(
				//insert
				insert, err)
		}
		str= "注册成功！"
	}
	
	}else{
		str="该用户名已被占用"
	}
	return str
}

func main() {
	//user := User{"abc", 123}
	initDB()
	http.HandleFunc("/", login)
	http.HandleFunc("/signup", saveUser)
	err:=http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	db.Close()
}
