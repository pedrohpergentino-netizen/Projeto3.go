package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

type Usuario struct {
	Idade int    `json:"idade"`
	Nome  string `json:"nome"`
}

func main() {
	db, err := sql.Open("sqlite", "./usuarios.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS usuarios (
			idade INTEGER,
			nome TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			listarUsuarios(db, w, r)

		case http.MethodPost:
			criarUsuario(db, w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Servidor rodando em http://localhost:6820/usuarios")
	log.Fatal(http.ListenAndServe(":6820", nil))
}

func listarUsuarios(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")

	var rows *sql.Rows
	var err error

	if nome != "" {
		rows, err = db.Query("SELECT idade, nome FROM usuarios WHERE nome = ?", nome)
	} else {
		rows, err = db.Query("SELECT idade, nome FROM usuarios")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pessoas []Usuario

	for rows.Next() {
		var u Usuario
		err := rows.Scan(&u.Idade, &u.Nome)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pessoas = append(pessoas, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

func criarUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var u Usuario

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO usuarios (idade, nome) VALUES (?, ?)", u.Idade, u.Nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
