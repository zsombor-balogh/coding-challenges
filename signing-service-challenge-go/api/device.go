package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type SignTransactionRequest struct {
	DeviceId string `json:"device_id"`
	Data     string `json:"data"`
}

func (s *Server) CreateSignatureDeviceHandler(response http.ResponseWriter, request *http.Request) {
	var createReq CreateSignatureDeviceRequest
	if err := json.NewDecoder(request.Body).Decode(&createReq); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	deviceID := uuid.New().String()

	device, err := domain.NewSignatureDevice(deviceID, createReq.Algorithm, createReq.Label)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	err = s.repo.SaveSignatureDevice(device)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, http.StatusCreated, device)
}

func (s *Server) SignTransactionHandler(response http.ResponseWriter, request *http.Request) {
	var signReq SignTransactionRequest
	if err := json.NewDecoder(request.Body).Decode(&signReq); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	device, err := s.repo.GetSignatureDevice(signReq.DeviceId)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{err.Error()})
		return
	}

	signature, signedData, err := device.SignTransaction(signReq.Data)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, http.StatusOK, map[string]string{
		"signature":   signature,
		"signed_data": signedData,
	})
}

func (s *Server) ListSignatureDevicesHandler(response http.ResponseWriter, request *http.Request) {
	devices, err := s.repo.ListSignatureDevices()
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, http.StatusOK, devices)
}

func (s *Server) GetSignatureDeviceHandler(response http.ResponseWriter, request *http.Request) {
	deviceId := mux.Vars(request)["device_id"]

	if deviceId == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{"Missing device_id parameter"})
		return
	}

	device, err := s.repo.GetSignatureDevice(deviceId)
	if err != nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, http.StatusOK, device)
}
