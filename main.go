package main

import (
	"bytes"
	"context"
	_ "database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var a = App{}

func main() {
	log.Println("starting API server")

	s := gocron.NewScheduler(time.UTC)

	a.Initialize()
	OAuthGmailService()
	//create a new router
	log.Println("creating routes")
	http.Handle("/", a.Router)

	//10AM every Friday
	//job, _ := s.Cron("0 10 * * 5").Do(func() {
	job, _ := s.Every(1).Hour().Do(func() {
		fmt.Println("Sending cron scheduled emails...")
		var companies, nil = a.GetCompanies()
		for _, company := range companies {
			_, err := SendEmailOAUTH2(company.UserEmail)
			if err != nil {
				fmt.Println("Error sending email for " + company.UserEmail + "!")
			}
		}
	})
	//log next run time for emails
	go func() {
		for {
			fmt.Println("Next run", job.NextRun())
			time.Sleep(time.Hour * 24)
		}
	}()
	//don't block other code execution
	s.StartAsync()
	//start and listen to requests
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("RedirectURL"),
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("AccessToken"),
		RefreshToken: os.Getenv("RefreshToken"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized")
	}
}

func SendEmailOAUTH2(to string) (bool, error) {
	var message gmail.Message

	type Data struct {
		URL     string
		Company string
	}

	emailBody, err := parseTemplate(Data{os.Getenv("SurveyEndpoint"), "tensure"})
	if err != nil {
		return false, errors.New("unable to parse email template")
	}

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Let us know how we're doing!" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send(os.Getenv("EmailUserName"), &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

func parseTemplate(data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("templates/email.html"))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
