package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type FormData struct {
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	PlanType string  `json:"planType"`
	Category string  `json:"category"`
	Value    float64 `json:"value"`
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf(
		`📝 *Nova Simulação Recebida*:
👤 Nome: %s
📞 Telefone: %s
📧 E-mail: %s
💳 Tipo: %s
📂 Categoria: %s
💰 Valor: R$ %.2f`,
		data.Name, data.Phone, data.Email, data.PlanType, data.Category, data.Value,
	)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	fmt.Println("TELEGRAM_CHAT_ID:", chatID)
	fmt.Println("TELEGRAM_BOT_TOKEN:", botToken)

	telegramPayload := TelegramMessage{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	payloadBytes, _ := json.Marshal(telegramPayload)

	apiURL := os.Getenv("API_URL")
	fmt.Println("API_URL:", apiURL)

	resp, err := http.Post(
		fmt.Sprintf("%s%s/sendMessage", apiURL, botToken),
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)

	if err != nil || resp.StatusCode >= 400 {
		http.Error(w, "Erro ao enviar mensagem para o Telegram", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis do ambiente")
	}

	http.HandleFunc("/send", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback para desenvolvimento local
	}
	log.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	log.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
