package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type Hits struct {
	Timestamp               int64  `json:"_timestamp"`
	Body                    string `json:"Body"`
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`
	ContentType             string `json:"Content-Type"`
	Date                    string `json:"Date"`
	From                    string `json:"From"`
	MessageID               string `json:"Message-ID"`
	MimeVersion             string `json:"Mime-Version"`
	Subject                 string `json:"Subject"`
	To                      string `json:"To"`
	Xbcc                    string `json:"X-bcc"`
	Xcc                     string `json:"X-cc"`
	XFileName               string `json:"X-FileName"`
	XFolder                 string `json:"X-Folder"`
	XFrom                   string `json:"X-From"`
	XOrigin                 string `json:"X-Origin"`
	XTo                     string `json:"X-To"`
}

type Response struct {
	Took       int `json:"took"`
	TookDetail struct {
		Total            int `json:"total"`
		WaitQueue        int `json:"wait_queue"`
		ClusterTotal     int `json:"cluster_total"`
		ClusterWaitQueue int `json:"cluster_wait_queue"`
	} `json:"took_detail"`
	Hits        []Hits `json:"hits"`
	Total       int    `json:"total"`
	From        int    `json:"from"`
	Size        int    `json:"size"`
	ScanSize    int    `json:"scan_size"`
	ScanRecords int    `json:"scan_records"`
	SessionID   string `json:"session_id"`
}

func searchData(data []byte) (*Response, error) {
	//user := "jgarcia025@gmail.com"
	//key := "kq8JH24j3i1675AQl90p"
	//url := "https://api.openobserve.ai/api/jhon_jairo_organization_6007_ktWTtLeSrd5sd63/default/_json"
	user := "root@example.com"
	key := "Complexpass#123"
	url := "http://localhost:5080/api/openobserve/_search"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Permitir solicitudes desde cualquier origen
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Permitir solicitudes GET y POST
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	// Permitir el encabezado Content-Type
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Verificar si la solicitud es de tipo OPTIONS (preflight)
	if r.Method == http.MethodOptions {
		// Respondemos sin hacer nada si es una solicitud OPTIONS
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parsear la consulta de búsqueda de los parámetros de la solicitud
	query := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")
	if query == "" {
		http.Error(w, "Missing 'q' parameter", http.StatusBadRequest)
		return
	}
	if page == "" {
		http.Error(w, "Missing 'page' parameter", http.StatusBadRequest)
		return
	}

	// Construir la consulta JSON con la frase de búsqueda y la página variable
	jsonString := fmt.Sprintf(`{
        "query": {
            "sql": "SELECT * FROM default WHERE Body LIKE '%%%s%%'",
            "start_time": 1707503203844188,
            "end_time":   1708456667721113,
            "from": %s,
            "size": 10
        }
    }`, query, page)
	jsonBytes := []byte(jsonString)

	// Realizar la búsqueda
	response, err := searchData(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Codificar y escribir la respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Hits)
}

func main() {
	// Crear un enrutador chi
	r := chi.NewRouter()

	// Manejar las solicitudes de búsqueda con chi
	r.Get("/search", searchHandler)
	// Agregar soporte para preflight OPTIONS
	r.Options("/search", searchHandler)

	// Iniciar el servidor con el enrutador chi
	fmt.Println("Servidor escuchando en el puerto :8080")
	http.ListenAndServe(":8080", r)
}
