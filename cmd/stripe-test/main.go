package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
)

func main() {
	// Set Stripe API key from environment
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY environment variable is required")
	}
	
	stripe.Key = stripeKey
	
	fmt.Println("🧪 Testing Stripe Integration")
	fmt.Println("=============================")
	
	// Test 1: Create a test customer
	fmt.Println("\n1. Testing customer creation...")
	testCustomer, err := createTestCustomer()
	if err != nil {
		log.Printf("❌ Customer creation failed: %v", err)
	} else {
		fmt.Printf("✅ Customer created: %s (%s)\n", testCustomer.Name, testCustomer.ID)
	}
	
	// Test 2: List existing products
	fmt.Println("\n2. Testing product listing...")
	products, err := listProducts()
	if err != nil {
		log.Printf("❌ Product listing failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d products\n", len(products))
		for _, prod := range products {
			fmt.Printf("   - %s (%s)\n", prod.Name, prod.ID)
		}
	}
	
	// Test 3: Create test products if none exist
	if len(products) == 0 {
		fmt.Println("\n3. Creating test products...")
		err := createTestProducts()
		if err != nil {
			log.Printf("❌ Product creation failed: %v", err)
		} else {
			fmt.Println("✅ Test products created successfully")
		}
	}
	
	// Test 4: List prices
	fmt.Println("\n4. Testing price listing...")
	prices, err := listPrices()
	if err != nil {
		log.Printf("❌ Price listing failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d prices\n", len(prices))
		for _, p := range prices {
			fmt.Printf("   - %s: $%.2f/%s\n", 
				p.Nickname, 
				float64(p.UnitAmount)/100, 
				p.Recurring.Interval)
		}
	}
	
	fmt.Println("\n🎉 Stripe integration test complete!")
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   1. Update your .env file with Stripe keys")
	fmt.Println("   2. Run: ./scripts/setup-stripe.sh")
	fmt.Println("   3. Launch beta program: ./scripts/launch-beta.sh")
	fmt.Println("   4. Start making money! 💰")
}

func createTestCustomer() (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String("Test Beta User"),
		Email: stripe.String("test@example.com"),
		Metadata: map[string]string{
			"source": "stripe-test",
			"tier":   "professional",
		},
	}
	
	return customer.New(params)
}

func listProducts() ([]*stripe.Product, error) {
	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "10")
	
	i := product.List(params)
	var products []*stripe.Product
	
	for i.Next() {
		products = append(products, i.Product())
	}
	
	return products, i.Err()
}

func listPrices() ([]*stripe.Price, error) {
	params := &stripe.PriceListParams{}
	params.Filters.AddFilter("limit", "", "20")
	
	i := price.List(params)
	var prices []*stripe.Price
	
	for i.Next() {
		prices = append(prices, i.Price())
	}
	
	return prices, i.Err()
}

func createTestProducts() error {
	// Create Starter product
	starterProduct, err := product.New(&stripe.ProductParams{
		Name:        stripe.String("AI Crypto Trading - Starter (Test)"),
		Description: stripe.String("Test product for Starter tier"),
		Type:        stripe.String("service"),
	})
	if err != nil {
		return fmt.Errorf("failed to create starter product: %v", err)
	}
	
	// Create Starter monthly price
	_, err = price.New(&stripe.PriceParams{
		Product:    stripe.String(starterProduct.ID),
		UnitAmount: stripe.Int64(4900), // $49.00
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		Nickname: stripe.String("Starter Monthly (Test)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create starter monthly price: %v", err)
	}
	
	// Create Professional product
	professionalProduct, err := product.New(&stripe.ProductParams{
		Name:        stripe.String("AI Crypto Trading - Professional (Test)"),
		Description: stripe.String("Test product for Professional tier"),
		Type:        stripe.String("service"),
	})
	if err != nil {
		return fmt.Errorf("failed to create professional product: %v", err)
	}
	
	// Create Professional monthly price
	_, err = price.New(&stripe.PriceParams{
		Product:    stripe.String(professionalProduct.ID),
		UnitAmount: stripe.Int64(19900), // $199.00
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		Nickname: stripe.String("Professional Monthly (Test)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create professional monthly price: %v", err)
	}
	
	// Create Enterprise product
	enterpriseProduct, err := product.New(&stripe.ProductParams{
		Name:        stripe.String("AI Crypto Trading - Enterprise (Test)"),
		Description: stripe.String("Test product for Enterprise tier"),
		Type:        stripe.String("service"),
	})
	if err != nil {
		return fmt.Errorf("failed to create enterprise product: %v", err)
	}
	
	// Create Enterprise monthly price
	_, err = price.New(&stripe.PriceParams{
		Product:    stripe.String(enterpriseProduct.ID),
		UnitAmount: stripe.Int64(99900), // $999.00
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		Nickname: stripe.String("Enterprise Monthly (Test)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create enterprise monthly price: %v", err)
	}
	
	fmt.Println("✅ Created test products:")
	fmt.Printf("   - Starter: %s\n", starterProduct.ID)
	fmt.Printf("   - Professional: %s\n", professionalProduct.ID)
	fmt.Printf("   - Enterprise: %s\n", enterpriseProduct.ID)
	
	return nil
}
