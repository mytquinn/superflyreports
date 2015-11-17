package formHandler

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
    "log"
	"config"
	"net/http"
	"strings"
	"encoding/json"
)

var db *sql.DB



func GetType (w http.ResponseWriter, r *http.Request) {
	typeTable := strings.Split(r.RequestURI, "/")[2]
	if !strings.Contains(typeTable, "type") {
		w.WriteHeader(402)
		w.Write([]byte("Only useable on type tables"))
		return
	}

	query := "select * from " + typeTable
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
	}
	rows, err := stmt.Query()
	if err == nil {
		defer rows.Close()
		columns, _ := rows.Columns()
		tableData := make([]map[string]interface{}, 0)
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for rows.Next() {
			for i := 0; i < len(columns); i++ {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			entry := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				entry[col] = v
			}
			tableData = append(tableData, entry)
		}
		jsonData, err := json.Marshal(tableData)
		if err != nil {
			log.Println(err)
		}
        w.WriteHeader(302)
		w.Write([]byte(string(jsonData)))
	} else {
		log.Println(err)
	}
}

func JournalSubmit (w http.ResponseWriter, r *http.Request) {


}


func init(){
	db, _ = sql.Open("mysql", myConfig.Db_user + ":" + myConfig.Db_password + "@" + myConfig.Db_address + "/" + myConfig.Db_schema)
	err := db.Ping()
	if err == nil {
		log.Println("DB responded")
	} else {
		log.Println("DB not responding: ", err)
	}
}