package cli

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/galihrivanto/kotak/config"
	"github.com/spf13/cobra"
)

var SendEmailCmd = &cobra.Command{
	Use:   "send-email",
	Short: "Send an email to a temporary email address",
	Run: func(cmd *cobra.Command, args []string) {
		email := promptInput("Email To: ")
		subject := promptInput("Subject: ")
		body := promptInput("Body: ")

		fmt.Printf(`Sending email to "%s" with subject "%s" and body "%s"\n`, email, subject, body)

		config := config.FromContext(cmd.Context())

		// construct the email
		msg := []byte("To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(
			fmt.Sprintf("%s:%s", config.SmtpServer.Host, config.SmtpServer.Port),
			nil,
			"no-reply@test.com",
			[]string{email},
			msg,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Email sent successfully")
	},
}

func promptInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	reader := bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	// Trim the trailing newline character
	return strings.TrimSpace(input)
}
