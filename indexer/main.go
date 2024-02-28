package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Email struct {
	MessageID               string `json:"Message-ID"`
	Date                    string `json:"Date"`
	From                    string `json:"From"`
	To                      string `json:"To"`
	Subject                 string `json:"Subject"`
	MimeVersion             string `json:"Mime-Version"`
	ContentType             string `json:"Content-Type"`
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`
	XFrom                   string `json:"X-From"`
	XTo                     string `json:"X-To"`
	Xcc                     string `json:"X-cc"`
	Xbcc                    string `json:"X-bcc"`
	XFolder                 string `json:"X-Folder"`
	XOrigin                 string `json:"X-Origin"`
	XFileName               string `json:"X-FileName"`
	Body                    string `json:"Body"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Falta especificar el directorio a indexar: go run main.go <directorio>")
		return
	}

	directory := os.Args[1]
	emails := make([]Email, 0)
	cant := 1
	cantmax := 40000

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if isValidFile(info.Name()) {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				email := parseEmail(string(content))
				emails = append(emails, email)
				if cant%cantmax == 0 {
					err = sendEmails(emails, cant)
					if err != nil {
						return err
					}
					// Vaciar emails
					emails = make([]Email, 0)
				}
				cant++
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", directory, err)
		return
	}
	// Enviar los correos restantes
	sendEmails(emails, cant-1)
} // Fin main

func isValidFile(filename string) bool {
	// Expresión regular para hacer coincidir los nombres de archivo que comienzan con números y terminan con "_"
	re := regexp.MustCompile(`^\d+_$`)
	return re.MatchString(filename)
} // Fin isValidFile

func parseEmail(content string) Email {
	email := Email{}

	// Dividir el contenido del correo electrónico en líneas
	lines := strings.Split(content, "\n")

	// Iterar sobre las líneas y parsear los campos relevantes
	for _, line := range lines {
		if strings.HasPrefix(line, "Message-ID:") {
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Date:") {
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		} else if strings.HasPrefix(line, "From:") {
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To:") {
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject:") {
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		} else if strings.HasPrefix(line, "Mime-Version:") {
			email.MimeVersion = strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:"))
		} else if strings.HasPrefix(line, "Content-Type:") {
			email.ContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
		} else if strings.HasPrefix(line, "Content-Transfer-Encoding:") {
			email.ContentTransferEncoding = strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:"))
		} else if strings.HasPrefix(line, "X-From:") {
			email.XFrom = strings.TrimSpace(strings.TrimPrefix(line, "X-From:"))
		} else if strings.HasPrefix(line, "X-To:") {
			email.XTo = strings.TrimSpace(strings.TrimPrefix(line, "X-To:"))
		} else if strings.HasPrefix(line, "X-cc:") {
			email.Xcc = strings.TrimSpace(strings.TrimPrefix(line, "X-cc:"))
		} else if strings.HasPrefix(line, "X-bcc:") {
			email.Xbcc = strings.TrimSpace(strings.TrimPrefix(line, "X-bcc:"))
		} else if strings.HasPrefix(line, "X-Folder:") {
			email.XFolder = strings.TrimSpace(strings.TrimPrefix(line, "X-Folder:"))
		} else if strings.HasPrefix(line, "X-Origin:") {
			email.XOrigin = strings.TrimSpace(strings.TrimPrefix(line, "X-Origin:"))
		} else if strings.HasPrefix(line, "X-FileName:") {
			email.XFileName = strings.TrimSpace(strings.TrimPrefix(line, "X-FileName:"))
		} else {
			email.Body += line + "\n"
		}
	}

	return email
} // Fin parseEmail

func sendEmails(emails []Email, cant int) error {
	if len(emails) > 0 {
		// Convertir las estructuras de Email en JSON
		emailsJSON, err := json.Marshal(emails)
		if err != nil {
			return err
		}
		// Indexar los datos en ZincSearch
		if err := indexData(emailsJSON); err != nil {
			return err
		}
	}
	return nil
} // Fin sendEmails

func indexData(data []byte) error {
	//user := "jgarcia025@gmail.com"
	//key := "kq8JH24j3i1675AQl90p"
	//url := "https://api.openobserve.ai/api/jhon_jairo_organization_6007_ktWTtLeSrd5sd63/default/_json"
	user := "root@example.com"
	key := "Complexpass#123"
	url := "http://localhost:5080/api/openobserve/default/_json"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.SetBasicAuth(user, key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	//log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
} // Fin indexData
