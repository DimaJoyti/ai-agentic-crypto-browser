package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ai-agentic-browser/internal/billing"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v76"
)

// StripeWebhookHandlers handles Stripe webhook events
type StripeWebhookHandlers struct {
	stripeProcessor     *billing.StripePaymentProcessor
	subscriptionManager *billing.SubscriptionManager
}

// NewStripeWebhookHandlers creates new Stripe webhook handlers
func NewStripeWebhookHandlers(
	stripeProcessor *billing.StripePaymentProcessor,
	subscriptionManager *billing.SubscriptionManager,
) *StripeWebhookHandlers {
	return &StripeWebhookHandlers{
		stripeProcessor:     stripeProcessor,
		subscriptionManager: subscriptionManager,
	}
}

// RegisterRoutes registers Stripe webhook routes
func (swh *StripeWebhookHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/webhooks/stripe", swh.HandleStripeWebhook).Methods("POST")
}

// HandleStripeWebhook processes incoming Stripe webhook events
func (swh *StripeWebhookHandlers) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Get the Stripe signature header
	sigHeader := r.Header.Get("Stripe-Signature")

	// Verify and construct the event
	event, err := swh.stripeProcessor.HandleWebhook(payload, sigHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the event
	err = swh.processStripeEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// processStripeEvent processes different types of Stripe events
func (swh *StripeWebhookHandlers) processStripeEvent(event *stripe.Event) error {
	// Process the webhook event using the Stripe processor
	return swh.stripeProcessor.ProcessWebhookEvent(event)
}
