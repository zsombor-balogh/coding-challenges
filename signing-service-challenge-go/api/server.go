package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

const (
	apiVersion = "v0"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	repo          persistence.SignatureDeviceRepository
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, repo persistence.SignatureDeviceRepository) *Server {
	return &Server{
		listenAddress: listenAddress,
		repo:          repo,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := mux.NewRouter()

	router.
		HandleFunc(fmt.Sprintf("/api/%s/health", apiVersion), s.Health).
		Methods(http.MethodGet)
	router.
		HandleFunc(fmt.Sprintf("/api/%s/devices", apiVersion), s.CreateSignatureDeviceHandler).
		Methods(http.MethodPost)
	router.
		HandleFunc(fmt.Sprintf("/api/%s/transactions/sign", apiVersion), s.SignTransactionHandler).
		Methods(http.MethodPost)
	router.
		HandleFunc(fmt.Sprintf("/api/%s/devices/{device_id}", apiVersion), s.GetSignatureDeviceHandler).
		Methods(http.MethodGet)
	router.
		HandleFunc(fmt.Sprintf("/api/%s/devices", apiVersion), s.ListSignatureDevicesHandler).
		Methods(http.MethodGet)

	server := &http.Server{
		Addr:    s.listenAddress,
		Handler: router,
	}

	log.Println("Server listening on", s.listenAddress)
	return server.ListenAndServe()
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
