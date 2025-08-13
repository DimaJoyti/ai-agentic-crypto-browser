package billing

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripePaymentProcessor handles Stripe payment processing
type StripePaymentProcessor struct {
	webhookSecret string
}

// NewStripePaymentProcessor creates a new Stripe payment processor
func NewStripePaymentProcessor() *StripePaymentProcessor {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	return &StripePaymentProcessor{
		webhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

// CreateCustomer creates a new Stripe customer
func (spp *StripePaymentProcessor) CreateCustomer(ctx context.Context, userID, email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}

	return customer.New(params)
}

// CreateSubscription creates a new Stripe subscription
func (spp *StripePaymentProcessor) CreateSubscription(ctx context.Context, customerID, priceID string, trialDays int64) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
		Expand: []*string{
			stripe.String("latest_invoice.payment_intent"),
		},
	}

	if trialDays > 0 {
		params.TrialPeriodDays = stripe.Int64(trialDays)
	}

	return subscription.New(params)
}

// CancelSubscription cancels a Stripe subscription
func (spp *StripePaymentProcessor) CancelSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{}
	return subscription.Cancel(subscriptionID, params)
}

// AttachPaymentMethod attaches a payment method to a customer
func (spp *StripePaymentProcessor) AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	return paymentmethod.Attach(paymentMethodID, params)
}

// CreateProducts creates Stripe products for subscription tiers
func (spp *StripePaymentProcessor) CreateProducts(ctx context.Context) error {
	tiers := []struct {
		ID           string
		Name         string
		Description  string
		MonthlyPrice int64
		AnnualPrice  int64
	}{
		{
			ID:           "starter",
			Name:         "Starter Plan",
			Description:  "Perfect for beginners getting started with AI trading",
			MonthlyPrice: 4900,  // $49.00 in cents
			AnnualPrice:  49000, // $490.00 in cents (2 months free)
		},
		{
			ID:           "professional",
			Name:         "Professional Plan",
			Description:  "Advanced features for serious traders",
			MonthlyPrice: 19900,  // $199.00 in cents
			AnnualPrice:  199000, // $1990.00 in cents (2 months free)
		},
		{
			ID:           "enterprise",
			Name:         "Enterprise Plan",
			Description:  "Full platform access with enterprise features",
			MonthlyPrice: 99900,  // $999.00 in cents
			AnnualPrice:  999000, // $9990.00 in cents (2 months free)
		},
	}

	for _, tier := range tiers {
		// Create product
		productParams := &stripe.ProductParams{
			ID:          stripe.String(tier.ID),
			Name:        stripe.String(tier.Name),
			Description: stripe.String(tier.Description),
			Type:        stripe.String("service"),
		}

		prod, err := product.New(productParams)
		if err != nil {
			log.Printf("Error creating product %s: %v", tier.ID, err)
			continue
		}

		// Create monthly price
		monthlyPriceParams := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(tier.MonthlyPrice),
			Currency:   stripe.String("usd"),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("month"),
			},
			Nickname: stripe.String(fmt.Sprintf("%s Monthly", tier.Name)),
		}

		_, err = price.New(monthlyPriceParams)
		if err != nil {
			log.Printf("Error creating monthly price for %s: %v", tier.ID, err)
		}

		// Create annual price
		annualPriceParams := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(tier.AnnualPrice),
			Currency:   stripe.String("usd"),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("year"),
			},
			Nickname: stripe.String(fmt.Sprintf("%s Annual", tier.Name)),
		}

		_, err = price.New(annualPriceParams)
		if err != nil {
			log.Printf("Error creating annual price for %s: %v", tier.ID, err)
		}

		log.Printf("Successfully created product and prices for %s", tier.Name)
	}

	return nil
}

// HandleWebhook processes Stripe webhooks
func (spp *StripePaymentProcessor) HandleWebhook(payload []byte, sigHeader string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, sigHeader, spp.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("webhook signature verification failed: %v", err)
	}

	return &event, nil
}

// ProcessWebhookEvent processes different types of webhook events
func (spp *StripePaymentProcessor) ProcessWebhookEvent(event *stripe.Event) error {
	switch event.Type {
	case "customer.subscription.created":
		// Handle subscription creation
		log.Printf("Subscription created: %s", event.Data.Object["id"])

	case "customer.subscription.updated":
		// Handle subscription updates
		log.Printf("Subscription updated: %s", event.Data.Object["id"])

	case "customer.subscription.deleted":
		// Handle subscription cancellation
		log.Printf("Subscription cancelled: %s", event.Data.Object["id"])

	case "invoice.payment_succeeded":
		// Handle successful payment
		log.Printf("Payment succeeded for invoice: %s", event.Data.Object["id"])

	case "invoice.payment_failed":
		// Handle failed payment
		log.Printf("Payment failed for invoice: %s", event.Data.Object["id"])

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

// GetSubscriptionPrices returns Stripe price IDs for subscription tiers
func (spp *StripePaymentProcessor) GetSubscriptionPrices() map[string]map[string]string {
	return map[string]map[string]string{
		"starter": {
			"monthly": "price_starter_monthly", // Replace with actual Stripe price IDs
			"annual":  "price_starter_annual",
		},
		"professional": {
			"monthly": "price_professional_monthly",
			"annual":  "price_professional_annual",
		},
		"enterprise": {
			"monthly": "price_enterprise_monthly",
			"annual":  "price_enterprise_annual",
		},
	}
}

// CreateBetaDiscountCoupon creates a discount coupon for beta users
func (spp *StripePaymentProcessor) CreateBetaDiscountCoupon() error {
	// This would create a 50% discount coupon for beta users
	// Implementation depends on your specific discount strategy
	log.Println("Beta discount coupon creation - implement based on your strategy")
	return nil
}
