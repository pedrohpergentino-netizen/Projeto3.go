package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

type Usuarios struct {
	Idade int    `json:"idade"`
	Nome  string `json:"nome"`
}

func main() {
	db, err := sql.Open("sqlite", "./usuarios.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS usuarios (idade INTEGER, nome TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT OR IGNORE INTO usuarios (idade, nome) VALUES (39, 'Hagno'), (20, 'Gabriel'), (19, 'Pedro')")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT idade, nome FROM usuarios")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var pessoas []Usuarios
		for rows.Next() {
			var u Usuarios
			err := rows.Scan(&u.Idade, &u.Nome)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pessoas = append(pessoas, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pessoas)
	})

	log.Println("Servidor iniciando, use http://localhost:6820/usuarios")
	log.Fatal(http.ListenAndServe(":6820", nil))
}
