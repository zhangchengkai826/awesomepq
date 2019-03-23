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
	tmplForm := template.Must(template.ParseFiles("login.html"))

	var tableInfo []TableInfo
	var cols []ColInfo
	var rs []Row

	if tbname, ok := mux.Vars(r)["tbname"]; ok {
		// Already log in.
		// Already choose a table.
		var (
			colName string
			colType string
		)
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

		// retrieve column info
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
			log.Fatal(err)
		}

		err = tmplForm.Execute(w, TmplData{false, tableInfo, cols, rs})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	} else if r.Method != http.MethodPost {
		// Not log in
		err := tmplForm.Execute(w, TmplData{true, tableInfo, cols, rs})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	// Already log in.
	// Not choose a table.
	dbname = r.FormValue("dbname")
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

	err = tmplForm.Execute(w, TmplData{false, tableInfo, cols, rs})
	if err != nil {
		log.Fatal(err.Error())
	}
}
