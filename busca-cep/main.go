package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if cep == "" {
		http.Error(w, "Parâmetro 'cep' ausente na consulta", http.StatusBadRequest)
		return
	}

	res, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		http.Error(w, "Erro ao buscar CEP", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var data CepResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Erro ao decodificar JSON: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

	fmt.Printf("CEP recebido: %s\n", cep)
	fmt.Printf("Informações do CEP: %+v\n", data)
}

func main() {
	http.HandleFunc("/cep", handler)
	fmt.Println("Servidor ouvindo na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
