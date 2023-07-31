package payment

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/nextpricedevelopers/go-next/pkg/model"
	"github.com/nextpricedevelopers/go-next/pkg/service/payment"
	"github.com/nextpricedevelopers/go-next/pkg/service/validation"

	"github.com/go-chi/chi/v5"
)

func createPayment(service payment.PaymentServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payment := &model.PaymentOptions{}
		name := r.URL.Query().Get("name")

		if name == "" {
			http.Error(w, "Payment name cannot be empty", http.StatusBadRequest)
			return
		}

		payment.PaymentName = name
		payment.Enabled = true

		result, err := service.Create(r.Context(), payment)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error to insert the payment"+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}

func updatePayment(service payment.PaymentServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := service.GetByID(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Payment Not Found", http.StatusNotFound)
			return
		}

		payment := &model.PaymentOptions{}
		name := r.URL.Query().Get("name")
		enable := r.URL.Query().Get("enable")
		booleanValue, err := strconv.ParseBool(enable)

		if err != nil {

			http.Error(w, "Erro ao converter a string para boolean:", http.StatusBadRequest)
			return
		}

		if name == "" {
			http.Error(w, "Payment name cannot be empty", http.StatusBadRequest)
			return
		}

		payment.PaymentName = validation.CareString(name)
		payment.Enabled = booleanValue
		_, err = service.Update(r.Context(), chi.URLParam(r, "id"), *payment)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error to update payment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}
