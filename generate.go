/**
 * 数据库表 生成struct
 * @lifeibest
 * 使用 go run generate.go table
 * 2016.12.28
 */
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

type TableColumns struct {
	Columns string
}

//change this const for own database config
const (
	DB_TYPE = "mysql"
	DB_HOST = "127.0.0.1"
	DB_PORT = "3306"
	DB_USER = "root"
	DB_PASS = "root"
	DB_NAME = "qqxinli"
)

func main() {
	arg_num := len(os.Args)
	if arg_num < 2 {
		panic("No mysql table choose")
	}
	//参数 go run generate.go table_name
	var table_name = string(os.Args[1])
	//var table_columns map[string]string

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		panic(err)
	}

	//sql := "show columns  "
	sql := fmt.Sprintf("show columns from `%s`", table_name)

	query, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(query)

	results := DbRows(query)

	// type Column struct {
	// 	Name    string
	// 	Comment string
	// }

	//查询字段注释
	column := make(map[string]string)
	var sql_1 = "select COLUMN_NAME,column_comment from INFORMATION_SCHEMA.Columns where table_name='" + table_name + "' and table_schema='" + DB_NAME + "'"
	query_1, err := db.Query(sql_1)
	if err != nil {
		fmt.Println(err)
	}
	results_1 := DbRows(query_1)
	for _, v := range results_1 {
		var k = Upstr(v["COLUMN_NAME"])
		column[k] = v["column_comment"]
	}
	// fmt.Println(column)

	fmt.Printf("type %s struct {\n", Upstr(table_name))
	for _, v := range results { //查询出来的数组
		//fmt.Println(k)
		Field := v["Field"]
		Type := v["Type"]
		Key := v["Key"]

		var s_field, s_type, s_key string
		s_field = Upstr(Field)

		//转换类型
		if strings.Contains(Type, "var") {
			s_type = "string"
		} else if strings.Contains(Type, "int") {
			s_type = "int"
		} else if strings.Contains(Type, "datetime") {
			s_type = "time.Time"
		} else {
			s_type = "string"
		}

		//主键等
		if strings.Contains(Key, "PRI") {
			s_key = "`PK`"
		} else {
			s_key = ""
		}
		if s_type == "time.Time" {
			s_key = "`orm:\"auto_now_add;type(datetime)\"`"
		}
		fmt.Println(s_field, s_type, s_key, "//", column[s_field])
	}
	fmt.Println("}")

}

func DbRows(query *sql.Rows) (results []map[string]string) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	//results = make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results = append(results, row) //装入结果集中
		i++
	}
	// for i, k := range results {
	// 	fmt.Println(i)
	// 	fmt.Println(k)
	// }
	return
}

//字母转大写，_转驼峰
// camel string, xx_yy to XxYy
func Upstr(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
