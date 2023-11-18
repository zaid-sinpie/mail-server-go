package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"net/smtp"
	"os"
)

var envFile, _ = godotenv.Read(".env")

type Viewer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var viewer Viewer
	err := json.NewDecoder(r.Body).Decode(&viewer)
	if err != nil {
		fmt.Println("Request wasn't properly decoded")
		return
	}
	fmt.Println(viewer)
	from := envFile["FROM_GMAIL"]
	if from == "" {
		from = os.Getenv("FROM_GMAIL")
	}
	password := envFile["GMAIL_PASSWORD"]
	if password == "" {
		password = os.Getenv("GMAIL_PASSWORD")
	}
	to := envFile["TO_GMAIL"]
	if to == "" {
		to = os.Getenv("TO_GMAIL")
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	subject := viewer.Subject
	body := viewer.Content
	message := fmt.Sprintf("Subject: %s(%s)\r\n\r\n%s", subject, viewer.Email, body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("Success")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	fmt.Println("Hello World, changed")
	mux := http.NewServeMux()
	mux.HandleFunc("/sendmail", sendMail)
	http.ListenAndServe(":"+port, mux)
}
