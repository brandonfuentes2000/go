package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Usuario struct {
	UUID     string `json:"uuid"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Ciudad   string `json:"ciudad"`
	Pais     string `json:"pais"`
}

type RandomUserResponse struct {
	Resultados []struct {
		Login struct {
			UUID string `json:"uuid"`
		} `json:"login"`
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email    string `json:"email"`
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
	} `json:"results"`
}

func usuariosGet() ([]Usuario, error) {
	usuarios := []Usuario{}
	usuariosUnicos := make(map[string]bool)

	apiURL := "https://randomuser.me/api/?results=5000"
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos de la API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %w", err)
	}

	var apiResponse RandomUserResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %w", err)
	}

	for _, resultado := range apiResponse.Resultados {
		uuid := resultado.Login.UUID

		if _, existe := usuariosUnicos[uuid]; !existe {
			user := Usuario{
				UUID:     uuid,
				Nombre:   resultado.Name.First,
				Apellido: resultado.Name.Last,
				Email:    resultado.Email,
				Ciudad:   resultado.Location.City,
				Pais:     resultado.Location.Country,
			}

			usuarios = append(usuarios, user)
			usuariosUnicos[uuid] = true
		}
	}

	return usuarios, nil
}

func getUsuariosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usuarios, err := usuariosGet()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/usuarios", getUsuariosHandler).Methods("GET")

	port := ":3000"
	log.Printf("Servidor corriendo en http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
