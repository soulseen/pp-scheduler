package sqlite

import "database/sql"

type KeyNodeTable struct {
	SQLiteDB *sql.DB `json:",omitempty"`
	Id       int     `json:",omitempty"`
	PodKey   string  `json:",omitempty"`
	NodeName string  `json:",omitempty"`
	Count    int     `json:",omitempty"`
}
