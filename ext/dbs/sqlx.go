package dbs

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	var err error
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		panic(err)
	}
	DB.SetMaxOpenConns(200)
	DB.SetMaxIdleConns(10)
}
func Insert(TablesName string, date map[string]interface{}) {
	// insert into user (id,name) values (1,"s")
	ntime := time.Now().Local().Unix()
	lens := len(date) + 1
	CoNames := make([]string, lens)
	CoVs := make([]string, lens)
	Args := make([]interface{}, lens)
	i := 0
	for k, v := range date {
		CoNames[i] = k
		CoVs[i] = "?"
		Args[i] = v
		i++
	}
	CoNames[i] = "create_time"
	CoVs[i] = "?"
	Args[i] = ntime

	CoName := strings.Join(CoNames, ",")
	CoV := strings.Join(CoVs, ",")
	query := fmt.Sprint("insert into ", TablesName, " (", CoName, ") values (", CoV, ")")
	DB.MustExec(query, Args...)
}
func Update(TablesName string, date map[string]interface{}, whereData map[string]interface{}) {
	// UPDATE table_name
	// SET column1=value1,column2=value2,...
	// WHERE some_column=some_value;
	Set, Sa := set(date)
	Where, Wa := where(whereData)
	query := fmt.Sprint("update  ", TablesName, " SET ", Set, " WHERE ", Where)
	Args := append(Sa, Wa...)
	fmt.Println()
	fmt.Println(query)
	fmt.Println()
	DB.MustExec(query, Args...)
}
func Delete(TablesName string, whereData map[string]interface{}) {
	a := make(map[string]interface{})
	a["delete_time"] = time.Now().Local().Unix()
	Update(TablesName, a, whereData)
}
func set(data map[string]interface{}) (Set string, Args []interface{}) {
	lens := len(data) + 1
	Sets := make([]string, lens)
	Args = make([]interface{}, lens)
	i := 0
	for k, v := range data {
		Sets[i] = fmt.Sprint(k, "=?")
		Args[i] = v
		i++
	}
	Sets[i] = "update_time=?"
	Args[i] = time.Now().Local().Unix()
	Set = strings.Join(Sets, ",")
	return
}
func where(data map[string]interface{}) (Where string, Args []interface{}) {
	lens := len(data)
	Wheres := make([]string, lens)
	Args = make([]interface{}, lens)
	i := 0
	for k, v := range data {
		Wheres[i] = fmt.Sprint(k, "=?")
		Args[i] = v
		i++
	}
	Where = strings.Join(Wheres, " AND ")
	return
}
