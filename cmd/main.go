package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	priceModels "github.com/jupiter-go/gen/go/price/models"
	priceServices "github.com/jupiter-go/gen/go/price/services"
	swapModels "github.com/jupiter-go/gen/go/swap/models"
	swapServices "github.com/jupiter-go/gen/go/swap/services"
	tokenServices "github.com/jupiter-go/gen/go/token/services"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// Jupiter API base URLs
	PriceAPIBase = "https://api.jup.ag/price/v2"
	SwapAPIBase  = "https://quote-api.jup.ag/v6"
	TokenAPIBase = "https://tokens.jup.ag"
)

// JupiterClient wraps HTTP client with Jupiter API endpoints
type JupiterClient struct {
	httpClient *http.Client
}

// NewJupiterClient creates a new Jupiter API client
func NewJupiterClient() *JupiterClient {
	return &JupiterClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// normalizeEnumValues fixes enum values to match protobuf expectations
func normalizeEnumValues(jsonData []byte) []byte {
	jsonStr := string(jsonData)

	// Fix confidence level enum values
	jsonStr = strings.ReplaceAll(jsonStr, `"confidenceLevel":"high"`, `"confidenceLevel":"CONFIDENCE_LEVEL_HIGH"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"confidenceLevel":"medium"`, `"confidenceLevel":"CONFIDENCE_LEVEL_MEDIUM"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"confidenceLevel":"low"`, `"confidenceLevel":"CONFIDENCE_LEVEL_LOW"`)

	// Fix swap mode enum values
	jsonStr = strings.ReplaceAll(jsonStr, `"swapMode":"ExactIn"`, `"swapMode":"SWAP_MODE_EXACTIN"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"swapMode":"ExactOut"`, `"swapMode":"SWAP_MODE_EXACTOUT"`)

	// Handle null values that cause parsing issues
	jsonStr = strings.ReplaceAll(jsonStr, `:null`, `:""`)

	return []byte(jsonStr)
}

// GetPrice fetches token prices using the Price API
func (c *JupiterClient) GetPrice(ctx context.Context, req *priceServices.RootGetRequest) (*priceServices.RootGetResponse, error) {
	// Build query parameters
	params := url.Values{}
	if req.Ids != "" {
		params.Set("ids", req.Ids)
	}
	if req.VsToken != "" {
		params.Set("vsToken", req.VsToken)
	}
	if req.ShowExtraInfo != "" {
		params.Set("showExtraInfo", req.ShowExtraInfo)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s?%s", PriceAPIBase, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse response into protobuf
	response := &priceServices.RootGetResponse{}
	if resp.StatusCode == 200 {
		// Normalize enum values before parsing
		normalizedBody := normalizeEnumValues(body)

		var priceResp priceModels.PriceResponse
		if err := protojson.Unmarshal(normalizedBody, &priceResp); err != nil {
			// If protojson fails, try regular JSON parsing for debugging
			var rawData map[string]interface{}
			if jsonErr := json.Unmarshal(body, &rawData); jsonErr == nil {
				fmt.Printf("Raw JSON response: %+v\n", rawData)
			}
			return nil, fmt.Errorf("unmarshaling price response: %w", err)
		}
		response.Response = &priceServices.RootGetResponse_PriceResponse_200{
			PriceResponse_200: &priceResp,
		}
	} else {
		// Handle error responses
		response.Response = &priceServices.RootGetResponse_Empty_400{}
	}

	return response, nil
}

// GetQuote fetches swap quotes using the Swap API
func (c *JupiterClient) GetQuote(ctx context.Context, req *swapServices.QuoteGetRequest) (*swapServices.QuoteGetResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("inputMint", req.InputMint)
	params.Set("outputMint", req.OutputMint)
	params.Set("amount", strconv.Itoa(int(req.Amount)))

	if req.SlippageBps > 0 {
		params.Set("slippageBps", strconv.Itoa(int(req.SlippageBps)))
	}
	if req.SwapMode != "" {
		params.Set("swapMode", req.SwapMode)
	}
	if len(req.Dexes) > 0 {
		for _, dex := range req.Dexes {
			params.Add("dexes", dex)
		}
	}
	if req.OnlyDirectRoutes {
		params.Set("onlyDirectRoutes", "true")
	}
	if req.PlatformFeeBps > 0 {
		params.Set("platformFeeBps", strconv.Itoa(int(req.PlatformFeeBps)))
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/quote?%s", SwapAPIBase, params.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse response into protobuf
	response := &swapServices.QuoteGetResponse{}
	if resp.StatusCode == 200 {
		// Normalize enum values before parsing
		normalizedBody := normalizeEnumValues(body)

		var quoteResp swapModels.QuoteResponse
		if err := protojson.Unmarshal(normalizedBody, &quoteResp); err != nil {
			// If protojson fails, try regular JSON parsing for debugging
			var rawData map[string]interface{}
			if jsonErr := json.Unmarshal(body, &rawData); jsonErr == nil {
				fmt.Printf("Raw JSON response: %+v\n", rawData)
			}
			return nil, fmt.Errorf("unmarshaling quote response: %w", err)
		}
		response.Response = &swapServices.QuoteGetResponse_QuoteResponse_200{
			QuoteResponse_200: &quoteResp,
		}
	}

	return response, nil
}

// GetTokenInfo fetches token information
func (c *JupiterClient) GetTokenInfo(ctx context.Context, address string) (*tokenServices.TokenAddressGetResponse, error) {
	url := fmt.Sprintf("%s/token/%s", TokenAPIBase, address)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// For demonstration, just return a basic response
	// In a real implementation, you'd parse the JSON into the appropriate protobuf type
	response := &tokenServices.TokenAddressGetResponse{}

	// Parse the JSON to show token info nicely
	var tokenInfo map[string]interface{}
	if err := json.Unmarshal(body, &tokenInfo); err == nil {
		fmt.Printf("✅ Token info received:\n")
		if name, ok := tokenInfo["name"].(string); ok {
			fmt.Printf("   Name: %s\n", name)
		}
		if symbol, ok := tokenInfo["symbol"].(string); ok {
			fmt.Printf("   Symbol: %s\n", symbol)
		}
		if decimals, ok := tokenInfo["decimals"].(float64); ok {
			fmt.Printf("   Decimals: %.0f\n", decimals)
		}
		if volume, ok := tokenInfo["daily_volume"].(float64); ok {
			fmt.Printf("   Daily Volume: $%.2f\n", volume)
		}
	}

	return response, nil
}

func main() {
	ctx := context.Background()
	client := NewJupiterClient()

	fmt.Println("🚀 Jupiter Go Proto Client Demo")
	fmt.Println("================================")

	// Example 1: Get SOL price
	fmt.Println("\n📊 Getting SOL price...")
	priceReq := &priceServices.RootGetRequest{
		Ids: "So11111111111111111111111111111111111111112", // SOL mint address
	}

	priceResp, err := client.GetPrice(ctx, priceReq)
	if err != nil {
		log.Printf("Error getting price: %v", err)
	} else if priceData := priceResp.GetPriceResponse_200(); priceData != nil {
		fmt.Printf("✅ Price data received (time taken: %.3fs)\n", priceData.TimeTaken)
		for tokenId, info := range priceData.Data {
			fmt.Printf("   Token: %s\n", tokenId)
			fmt.Printf("   Price: $%s\n", info.Price)
			fmt.Printf("   Type: %s\n", info.Type)
		}
	}

	// Example 2: Get price with extra info
	fmt.Println("\n📈 Getting SOL price with extra info...")
	priceReqExtra := &priceServices.RootGetRequest{
		Ids:           "So11111111111111111111111111111111111111112",
		ShowExtraInfo: "true",
	}

	priceRespExtra, err := client.GetPrice(ctx, priceReqExtra)
	if err != nil {
		log.Printf("Error getting price with extra info: %v", err)
	} else if priceData := priceRespExtra.GetPriceResponse_200(); priceData != nil {
		fmt.Printf("✅ Price data with extra info received\n")
		for tokenId, info := range priceData.Data {
			fmt.Printf("   Token: %s\n", tokenId)
			fmt.Printf("   Price: $%s\n", info.Price)
			if extraInfo := info.ExtraInfo; extraInfo != nil {
				fmt.Printf("   Confidence: %s\n", extraInfo.ConfidenceLevel.String())
				if quoted := extraInfo.QuotedPrice; quoted != nil && quoted.BuyPrice != "" {
					fmt.Printf("   Buy Price: $%s\n", quoted.BuyPrice)
					fmt.Printf("   Sell Price: $%s\n", quoted.SellPrice)
				}
			}
		}
	}

	// Example 3: Get swap quote (SOL to USDC)
	fmt.Println("\n🔄 Getting swap quote (SOL → USDC)...")
	quoteReq := &swapServices.QuoteGetRequest{
		InputMint:   "So11111111111111111111111111111111111111112",  // SOL
		OutputMint:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		Amount:      1000000000,                                     // 1 SOL (in lamports)
		SlippageBps: 50,                                             // 0.5% slippage
		SwapMode:    "ExactIn",                                      // Use the API's expected value
	}

	quoteResp, err := client.GetQuote(ctx, quoteReq)
	if err != nil {
		log.Printf("Error getting quote: %v", err)
	} else if quoteData := quoteResp.GetQuoteResponse_200(); quoteData != nil {
		fmt.Printf("✅ Quote received\n")
		fmt.Printf("   Input Amount: %s\n", quoteData.InAmount)
		fmt.Printf("   Output Amount: %s\n", quoteData.OutAmount)
		fmt.Printf("   Price Impact: %s%%\n", quoteData.PriceImpactPct)
		fmt.Printf("   Route Plan Steps: %d\n", len(quoteData.RoutePlan))

		if len(quoteData.RoutePlan) > 0 && quoteData.RoutePlan[0].SwapInfo != nil {
			fmt.Printf("   First DEX: %s\n", quoteData.RoutePlan[0].SwapInfo.Label)
		}
	}

	// Example 4: Get token information
	fmt.Println("\n🪙 Getting token information for USDC...")
	_, err = client.GetTokenInfo(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	if err != nil {
		log.Printf("Error getting token info: %v", err)
	}

	// Example 5: Demonstrate error handling
	fmt.Println("\n❌ Testing error handling with invalid token...")
	invalidPriceReq := &priceServices.RootGetRequest{
		Ids: "invalid-token-address",
	}

	invalidResp, err := client.GetPrice(ctx, invalidPriceReq)
	if err != nil {
		fmt.Printf("✅ Properly handled error: %v\n", err)
	} else if invalidResp.GetEmpty_400() != nil {
		fmt.Printf("✅ Properly handled 400 error response\n")
	}

	fmt.Println("\n🎉 Demo completed!")
	fmt.Println("\nThis demonstrates how to:")
	fmt.Println("• Use generated protobuf types for type-safe API calls")
	fmt.Println("• Make HTTP requests to Jupiter's REST APIs")
	fmt.Println("• Handle different response types and error cases")
	fmt.Println("• Parse JSON responses into protobuf structures")
	fmt.Println("• Normalize enum values for proper protobuf parsing")
}
