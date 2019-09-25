package sqlite

import (
	"database/sql"
	log "github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const (
	CREATE_TABLE_SQL = `
		CREATE TABLE IF NOT EXISTS key_node 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			podkey VARCHAR(64) NULL,
			nodename VARCHAR(64) NULL,
			count INT(10) NULL
		)
	;`
)

var KeyNodeCilent KeyNodeTable

func init() {
	db, err := sql.Open("sqlite3", os.Getenv("DATA_PATH"))
	err = db.Ping()
	checkErr(err)
	CreateTable(db, CREATE_TABLE_SQL)
	KeyNodeCilent = KeyNodeTable{
		SQLiteDB: db,
	}
}

func CreateTable(db *sql.DB, sql string) {
	_, error := db.Exec(sql)
	checkErr(error)
}

func (kn *KeyNodeTable) KeyNodeInsert(key, nodeName string, count int) (id int, err error) {
	tx, err := kn.SQLiteDB.Begin()
	stmt, err := tx.Prepare("INSERT INTO key_node(podkey, nodename, count) VALUES (?,?,?)")
	if err != nil {
		return 1, err
	}
	defer stmt.Close()

	stmt.Exec(key, nodeName, count)
	checkErr(err)
	tx.Commit()
	return id, nil
}

func (kn *KeyNodeTable) KeyNodeSearch(podKey string, nodeName string) ([]KeyNodeTable, error) {
	rows, err := kn.SQLiteDB.Query("SELECT *  FROM key_node WHERE podkey = ? AND nodename = ?", podKey, nodeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []KeyNodeTable{}

	for rows.Next() {
		var id, count int
		var podkey, nodename string

		err = rows.Scan(&id, &podkey, &nodename, &count)
		checkErr(err)

		res = append(res, KeyNodeTable{
			Id:     id,
			PodKey: podkey,
			Count:  count,
		})
	}
	return res, nil

}

func (kn *KeyNodeTable) KeyNodeUpdate(id int, count int) (int64, error) {
	//update
	stmt, err := kn.SQLiteDB.Prepare("update key_node set count=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(count, id)
	rowId, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return rowId, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
