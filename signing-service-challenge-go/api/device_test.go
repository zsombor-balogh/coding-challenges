package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func TestCreateSignatureDeviceHandler(t *testing.T) {
	mockRepo := persistence.NewMockRepository()
	server := NewServer(":8080", mockRepo)

	requestBody := CreateSignatureDeviceRequest{
		Algorithm: "RSA",
		Label:     "Test Device",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	server.CreateSignatureDeviceHandler(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response struct {
		Data *domain.SignatureDevice `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}

	fmt.Println(response.Data)

	if response.Data.Id == "" {
		t.Errorf("Response is missing 'id' field")
	}

	if response.Data.Algorithm != requestBody.Algorithm {
		t.Errorf("Expected algorithm %s, got %v", requestBody.Algorithm, response.Data.Algorithm)
	}
	if response.Data.Label != requestBody.Label {
		t.Errorf("Expected label %s, got %v", requestBody.Label, response.Data.Label)
	}
	if response.Data.LastSignature != "" {
		t.Errorf("Expected LastSignature to be empty, got %v", response.Data.LastSignature)
	}
	if response.Data.SignatureCounter != 0 {
		t.Errorf("Expected SignatureCounter to be 0, got %v", response.Data.SignatureCounter)
	}
}

func TestGetSignatureDeviceHandler(t *testing.T) {
	mockRepo := persistence.NewMockRepository()
	server := NewServer(":8080", mockRepo)

	deviceId := "device1"
	device := &domain.SignatureDevice{Id: deviceId, Algorithm: "rsa", Label: "Device 1"}
	mockRepo.Devices[deviceId] = device

	req, err := http.NewRequest("GET", "/devices/"+deviceId, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.Handle("/devices/{device_id}", http.HandlerFunc(server.GetSignatureDeviceHandler)).Methods("GET")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Data *domain.SignatureDevice `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}

	if response.Data == nil {
		t.Error("Response is missing 'data' field")
	} else {
		assertDeviceEquals(t, device, response.Data)
	}
}

func TestListSignatureDevicesHandler(t *testing.T) {
	mockRepo := persistence.NewMockRepository()
	server := NewServer(":8080", mockRepo)

	device1 := &domain.SignatureDevice{Id: "device1", Algorithm: "RSA", Label: "Device 1"}
	device2 := &domain.SignatureDevice{Id: "device2", Algorithm: "ECC", Label: "Device 2"}
	mockRepo.Devices["device1"] = device1
	mockRepo.Devices["device2"] = device2

	req, err := http.NewRequest("GET", "/devices", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	server.ListSignatureDevicesHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Data []*domain.SignatureDevice `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	if len(response.Data) != 2 {
		t.Errorf("Expected %d devices in response, got %d", 2, len(response.Data))
	}

	assertDeviceEquals(t, device1, response.Data[0])
	assertDeviceEquals(t, device2, response.Data[1])
}

func TestSignTransactionHandler(t *testing.T) {
	mockRepo := persistence.NewMockRepository()
	server := NewServer(":8080", mockRepo)

	deviceId := "mock-device-id"
	mockDevice, err := domain.NewSignatureDevice(deviceId, "RSA", "Test")
	mockRepo.Devices[deviceId] = mockDevice

	requestBody := SignTransactionRequest{
		DeviceId: deviceId,
		Data:     "test-data",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	server.SignTransactionHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Error unmarshaling response body")
	}

	if _, ok := data["signature"]; !ok {
		t.Errorf("Response is missing 'signature' field")
	}
	if _, ok := data["signed_data"]; !ok {
		t.Errorf("Response is missing 'signed_data' field")
	}
}

func assertDeviceEquals(t *testing.T, expected *domain.SignatureDevice, actual *domain.SignatureDevice) {
	if expected.Id != actual.Id {
		t.Errorf("Expected Id %s, got %s", expected.Id, actual.Id)
	}
	if expected.Algorithm != actual.Algorithm {
		t.Errorf("Expected algorithm %s, got %s", expected.Algorithm, actual.Algorithm)
	}
	if expected.Label != actual.Label {
		t.Errorf("Expected label %s, got %s", expected.Label, actual.Label)
	}
}
