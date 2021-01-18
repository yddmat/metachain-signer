package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"multichain-signer/signer"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type TxSigner interface {
	Sign(signerType string, txData json.RawMessage) ([]byte, error)
}

type Server struct {
	TxSigner TxSigner

	Host string
	Port string
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})

	router.Post("/api/v1/sign_transaction", s.signTransactionHandler)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Host, s.Port), router)
}

type SignRequest struct {
	Gate string          `json:"gate"`
	Tx   json.RawMessage `json:"tx"`
}

type SignResponse struct {
	Sign string `json:"sign"`
}

type SignErrResponse struct {
	Error string `json:"error"`
}

func (s *Server) signTransactionHandler(w http.ResponseWriter, r *http.Request) {
	request := SignRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	sig, err := s.TxSigner.Sign(request.Gate, request.Tx)
	if err != nil {
		if err == signer.ErrUnknownSignerType || err == signer.ErrInvalidTxData {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, SignErrResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("/sign: internal error:", err.Error())
		return
	}

	render.JSON(w, r, SignResponse{Sign: hex.EncodeToString(sig)})
}
