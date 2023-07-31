package payment

import (
	"github.com/go-chi/chi/v5"
	"github.com/nextpricedevelopers/go-next/pkg/service/payment"
)

func RegisterPaymentAPIHandlers(r chi.Router, service payment.PaymentServiceInterface) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/payment", createPayment(service))
		r.Put("/client/{id}", updatePayment(service))

	})
}
