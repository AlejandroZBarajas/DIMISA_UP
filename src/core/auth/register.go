package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	DB *sql.DB
}

type registerRequest struct {
	Nombres   string `json:"nombres"`
	Apellido1 string `json:"apellido1"`
	Apellido2 string `json:"apellido2"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type registerResponse struct {
	IdUsuario int32  `json:"id_usuario"`
	IdRol     int32  `json:"id_rol"`
	Username  string `json:"username"`
	Nombre    string `json:"nombre_usuario"`
}

func (rh *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	req.Nombres = strings.TrimSpace(req.Nombres)
	req.Apellido1 = strings.TrimSpace(req.Apellido1)
	req.Apellido2 = strings.TrimSpace(req.Apellido2)
	req.Username = strings.TrimSpace(req.Username)

	if req.Nombres == "" || req.Apellido1 == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Campos obligatorios incompletos", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encriptando contraseña", http.StatusInternalServerError)
		return
	}

	tx, err := rh.DB.Begin()
	if err != nil {
		http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	const idRolAdmin int32 = 1

	queryInsertUser := `
		INSERT INTO usuarios 
		(nombres, apellido1, apellido2, username, password, id_rol)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		queryInsertUser,
		req.Nombres,
		req.Apellido1,
		req.Apellido2,
		req.Username,
		string(hashedPassword),
		idRolAdmin,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creando usuario: %v", err), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error obteniendo ID del usuario", http.StatusInternalServerError)
		return
	}

	idUsuario := int32(lastID)

	queryInsertAdmin := `
		INSERT INTO admin_users 
		(id_user)
		VALUES (?)
	`

	if _, err := tx.Exec(queryInsertAdmin, idUsuario); err != nil {
		http.Error(w, fmt.Sprintf("Error asignando rol admin: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Error confirmando transacción", http.StatusInternalServerError)
		return
	}

	nombreCompleto := strings.TrimSpace(
		fmt.Sprintf("%s %s %s", req.Nombres, req.Apellido1, req.Apellido2),
	)

	resp := registerResponse{
		IdUsuario: idUsuario,
		IdRol:     idRolAdmin,
		Username:  req.Username,
		Nombre:    nombreCompleto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
