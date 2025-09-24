package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Makhaev/projectname/pkg/auth"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" || req.Password != expectedPassword {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(expectedPassword)
	if err != nil {
		http.Error(w, `{"error": "Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
