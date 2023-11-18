package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
    "github.com/joho/godotenv"
)
var envFile, _ = godotenv.Read(".env")
type Viewer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var viewer Viewer
	err := json.NewDecoder(r.Body).Decode(&viewer)
	if err != nil {
		fmt.Println("Request wasn't properly decoded")
        return
	}
    fmt.Println(viewer)
	from := "dulangeraviraj@gmail.com"
    password := envFile["GMAIL_PASSWORD"]
    if password == "" {
        password = os.Getenv("GMAIL_PASSWORD")
    }
    fmt.Println("Password = ", password)
	to := []string{"zaidmasuldar@gmail.com", "dulangeraviraj@gmail.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	subject := viewer.Subject
	body := viewer.Content
	message := fmt.Sprintf("Subject: %s(%s)\r\n\r\n%s", subject, viewer.Email, body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}
	fmt.Println("Success")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Hello World")
	http.HandleFunc("/sendmail", sendMail)
	http.ListenAndServe(":"+port, nil)

}