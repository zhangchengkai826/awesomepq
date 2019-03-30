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
	Cmd string
	History []interface{}
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
var user string
var pw string
var host string
var port string
var cmd string
var hist []interface{}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmplForm := template.Must(template.ParseFiles("index.html"))

	var tableInfo []TableInfo
	var cols []ColInfo
	var rs []Row
	var outs []interface{}

	if dbname == "" && r.Method != http.MethodPost {
		// 1. Not log in
		err := tmplForm.Execute(w, TmplData{true, tableInfo, cols, rs, outs, cmd, hist})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	// 2. Already log in.
	if dbname == "" {
		dbname = r.FormValue("dbname")
		user = r.FormValue("user")
		pw = r.FormValue("pw")
		host = r.FormValue("host")
		port = r.FormValue("port")
	}
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", dbname, user, pw, host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodPost && r.FormValue("cmd") != "" {
		cmd = r.FormValue("cmd")
		if r.FormValue("exec") != "" {
			// @Deal with user command (exec)
			_, err = db.Exec(cmd)
			if err != nil {
				outs = append(outs, err)
				goto ExecErrHander
			}
			hist = append(hist, cmd)
			cmd = ""
		}
	}
ExecErrHander:
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

	customQuery := false
	if r.Method == http.MethodPost && r.FormValue("cmd") != "" {
		cmd = r.FormValue("cmd")
		if r.FormValue("query") != "" {
			// @Deal with user command (query)
			rows, err = db.Query(cmd)
			if err != nil {
				outs = append(outs, err)
				goto QueryErrHandler
			}
			defer rows.Close()
			colNames, err := rows.Columns()
			if err != nil {
				outs = append(outs, err)
				goto QueryErrHandler
			}
			colTypes, err := rows.ColumnTypes()
			if err != nil {
				outs = append(outs, err)
				goto QueryErrHandler
			}
			for i, _ := range colNames {
				cols = append(cols, ColInfo{colNames[i], colTypes[i].DatabaseTypeName()})
			}
			for rows.Next() {
				cellsRaw := make([]interface{}, len(cols))
				cells := make([]Cell, len(cols))
				for i := range cellsRaw {
					cellsRaw[i] = &cells[i].Data
				}
				err := rows.Scan(cellsRaw...)
				if err != nil {
					outs = append(outs, err)
					goto QueryErrHandler
				}
				rs = append(rs, Row{cells})
			}
			err = rows.Err()
			if err != nil {
				outs = append(outs, err)
				goto QueryErrHandler
			}
			customQuery = true
			hist = append(hist, cmd)
			cmd = ""
		}
	}
QueryErrHandler:
	if tbname, ok := mux.Vars(r)["tbname"]; ok && !customQuery {
		// @Show a complete table
		// retrieve cols
		rows, err = db.Query(fmt.Sprintf("SELECT * FROM %s", tbname))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		colNames, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}
		colTypes, err := rows.ColumnTypes()
		if err != nil {
			log.Fatal(err)
		}
		for i, _ := range colNames {
			cols = append(cols, ColInfo{colNames[i], colTypes[i].DatabaseTypeName()})
		}
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

	err = tmplForm.Execute(w, TmplData{false, tableInfo, cols, rs, outs, cmd, hist})
	if err != nil {
		log.Fatal(err)
	}
}