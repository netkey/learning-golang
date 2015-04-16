/*
  test postgres database with go


  sql:

  CREATE TABLE userinfo
(
    uid serial NOT NULL,
    username character varying(100) NOT NULL,
    departname character varying(500) NOT NULL,
    Created date,
    CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE userdeatail
(
    uid integer,
    intro character varying(100),
    profile character varying(100)
)
WITH(OIDS=FALSE);


*/
//
package main

import (
	"database/sql"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"os"
)

// database
var engine *xorm.Engine

func main() {

	var err error
	engine, err = xorm.NewEngine("postgres", "user=postgres password=postgres dbname=tsingcloud sslmode=disable")
	engine.ShowSQL = true
	engine.ShowWarn = true

	//db, err := sql.Open("postgres", "user=postgres password=postgres dbname=tsingcloud sslmode=disable")
	//checkError(err)

	// log for database
	f, err := os.Create("sql.log")
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()
	engine.Logger = xorm.NewSimpleLogger(f)

	//插入数据
	stmt, err := enginePrepare("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) RETURNING uid")
	checkError(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkError(err)

	//pg不支持这个函数，因为他没有类似MySQL的自增ID
	//id, err := res.LastInsertId()
	//checkError(err)

	//fmt.Println(id)

	//更新数据
	stmt, err = enginePrepare("update userinfo set username=$1 where uid=$2")
	checkError(err)

	res, err = stmt.Exec("xxxastaxieupdate", 1)
	checkError(err)

	affect, err := res.RowsAffected()
	checkError(err)

	fmt.Println(affect)

	//查询数据
	rows, err := engineQuery("SELECT * FROM userinfo")
	checkError(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkError(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//删除数据
	stmt, err = enginePrepare("delete from userinfo where uid=$1")
	checkError(err)

	res, err = stmt.Exec(1)
	checkError(err)

	affect, err = res.RowsAffected()
	checkError(err)

	fmt.Println(affect)

	engineClose()

}

func checkError(err error) {
	if err != nil {
		panic(err)
		//fmt.Println("error is s%", err)
	}
	os.Exit(1)
}
