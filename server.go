package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type TmplData struct {
	NeedLogIn bool
	Tables []TableInfo
	Cols []ColInfo
	Rows []Row
	Outputs []interface{}
}

type TableInfo struct {
	Name string
}

type ColInfo struct {
	Name string
	Type string
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Data interface{}
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	router.HandleFunc("/", mainHandler)
	router.HandleFunc("/tables/{tbname}", mainHandler)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

var dbname string

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmplForm := template.Must(template.ParseFiles("index.html"))

	var tableInfo []TableInfo
	var cols []ColInfo
	var rs []Row
	var outs []interface{}

	if dbname == "" && r.Method != http.MethodPost {
		// 1. Not log in
		err := tmplForm.Execute(w, TmplData{true, tableInfo, cols, rs, outs})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	// 2. Already log in.
	if dbname == "" {
		dbname = r.FormValue("dbname")
	}
	connStr := fmt.Sprintf("dbname=%s user=andys sslmode=disable", dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// @First deal with user command
	if r.Method == http.MethodPost && r.FormValue("cmd") != "" {
		cmd := r.FormValue("cmd")
		_, err = db.Exec(cmd)
		if err != nil {
			outs = append(outs, err)
		}
	}

	var name string
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		tableInfo = append(tableInfo, TableInfo{name})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// 2.1 Already choose a table.
	if tbname, ok := mux.Vars(r)["tbname"]; ok {
		// retrieve cols
		var (
			colName string
			colType string
		)
		rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", tbname)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&colName, &colType)
			if err != nil {
				log.Fatal(err)
			}
			cols = append(cols, ColInfo{colName, colType})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal()
		}

		rows, err = db.Query(fmt.Sprintf("SELECT * FROM %s", tbname))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			cellsRaw := make([]interface{}, len(cols))
			cells := make([]Cell, len(cols))
			for i := range cellsRaw {
				cellsRaw[i] = &cells[i].Data
			}
			err := rows.Scan(cellsRaw...)
			if err != nil {
				log.Fatal(err)
			}
			rs = append(rs, Row{cells})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tmplForm.Execute(w, TmplData{false, tableInfo, cols, rs, outs})
	if err != nil {
		log.Fatal(err)
	}
}