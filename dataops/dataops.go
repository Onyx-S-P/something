package dataops

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/xedflix/auto-approval-system/model"
	_ "modernc.org/sqlite"
)

/*Schema for policy Table*/
var initSchema = `
CREATE TABLE IF NOT EXISTS policy(
id INTEGER PRIMARY KEY AUTOINCREMENT, 
apiversion varchar(20), 
kind varchar(20), 
metadata  varchar(200),
spec varchar(2000)
);
`
var initSchema2 = `
CREATE TABLE IF NOT EXISTS specs(
id INTEGER PRIMARY KEY AUTOINCREMENT, 
tags varchar(200), 
message varchar(200), 
selector  varchar(2000),
file varchar(3000)
);
`

/*DBconn struct stores sql.DB struct*/
type DBconn struct {
	conn *sql.DB
}

/*create executes init schema in db*/
func (db *DBconn) CreateTables() (error, error) {
	_, err := db.conn.Exec(initSchema)
	_, err1 := db.conn.Exec(initSchema2)
	return err, err1
}

/*create new sql cilent for given dsn. if database is not exists then
creates new db file*/
func NewSqliteCilent(dsn string) (DBconn, error) {
	_, err := os.Stat(dsn)
	if os.IsNotExist(err) {
		log.Println("Creating sqlite database ", dsn)
		_, err := os.OpenFile(dsn, os.O_CREATE|os.O_WRONLY, 0660)
		if err != nil {
			return DBconn{}, err
		}
	}
	db, err := sql.Open("sqlite", dsn)
	return DBconn{conn: db}, err
}

/*insert given policy in table*/
func (db *DBconn) InsertPolicy(p model.Policy) (sql.Result, error) {
	s := p.GetSpec()
	stm, err := db.conn.Prepare("INSERT INTO policy(apiversion,kind,metadata,spec) values(?,?,?,?)")
	if err != nil {
		return nil, err
	}

	return stm.Exec(
		p.APIVersion,
		p.Kind,
		p.GetMetadata(),
		string(s),
	)
}

func (db *DBconn) Readpolicy(p model.Policy) {
	//var spec model.Spec
	//s := model.NewSpecFrom(p.GetSpec())
	rows, err := db.conn.Query("SELECT spec FROM policy")
	if err != nil {
		log.Println(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		err := rows.Scan(&p.Spec)
		if err != nil {
			log.Println(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(p.Spec)
	}
	/*stm, err1 := db.conn.Prepare("INSERT INTO specs(tags,message,selector,file) values(?,?,?,?)")
	if err1 != nil {
		return nil, err1
	}

	return stm.Exec(
		s.Tags,
		s.Message,
		s.Selector,
		s.File,
	)*/
}

/*func (db *DBconn) Insertspecs(p model.Policy) (sql.Result, error) {
	var spec model.Spec
	stm, err := db.conn.Prepare("INSERT INTO specs(tags,message,selector,file) values(?,?,?,?)")
	if err != nil {
		return nil, err
	}

	return stm.Exec(
		spec.Tags,
		spec.Message,
		spec.Selector,
		spec.File,
	)
}*/
