package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSNBScraper_Scrape(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)

		// Sample RSS response with SNB interest rates
		xml := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Current interest rates</title>
    <item>
      <title>CH: 0.25 SNBLZ 2025-04-04</title>
      <description>SNB policy rate / valid from 04.04.2025</description>
      <pubDate>Fri, 04 Apr 2025 10:16:33 GMT</pubDate>
      <country>CH</country>
      <code>SNBLZ</code>
      <value>0.25</value>
      <unit>percent</unit>
      <date>2025-04-04</date>
    </item>
    <item>
      <title>CH: 0.75 LSFF 2025-04-04</title>
      <description>Special rate (liquidity-shortage financing facility)</description>
      <pubDate>Fri, 04 Apr 2025 10:16:33 GMT</pubDate>
      <country>CH</country>
      <code>LSFF</code>
      <value>0.75</value>
      <unit>percent</unit>
      <date>2025-04-04</date>
    </item>
    <item>
      <title>CH: 0.386 R10 2025-04-04</title>
      <description>Yield on Swiss Confederation bonds</description>
      <pubDate>Fri, 04 Apr 2025 10:16:33 GMT</pubDate>
      <country>CH</country>
      <code>R10</code>
      <value>0.386</value>
      <unit>percent</unit>
      <date>2025-04-04</date>
    </item>
  </channel>
</rss>`
		_, _ = w.Write([]byte(xml))
	}))
	defer mockServer.Close()

	// Create scraper with mock server URL
	scraper := &SNBScraper{
		rssURL:     mockServer.URL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// Run the scraper
	results, err := scraper.Scrape(context.Background())
	require.NoError(t, err, "Scrape should not return an error")
	require.Len(t, results, 1, "Should return exactly 1 result")

	result := results[0]
	assert.Equal(t, "snb_interest_rates", result.Source, "Result source should match scraper name")

	// Type assertion
	rates, ok := result.Data.([]SNBInterestRate)
	require.True(t, ok, "Result data should be of type []SNBInterestRate")
	require.Len(t, rates, 3, "Should return exactly 3 rate items")

	// Test specific values
	expectedRates := map[string]float64{
		"SNBLZ": 0.25,
		"LSFF":  0.75,
		"R10":   0.386,
	}

	expectedDate, _ := time.Parse("2006-01-02", "2025-04-04")

	for _, rate := range rates {
		expectedValue, exists := expectedRates[rate.Code]
		require.True(t, exists, "Unexpected rate code: %s", rate.Code)

		assert.Equal(t, expectedValue, rate.Value, "Rate %s should have correct value", rate.Code)
		assert.True(t, rate.Date.Equal(expectedDate), "Rate %s should have correct date", rate.Code)
		assert.Equal(t, "percent", rate.Unit, "Rate %s should have unit 'percent'", rate.Code)
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"0.25", 0.25, false},
		{"0.25 percent", 0.25, false},
		{"1.75%", 0.0, true}, // Should fail as % isn't handled
		{"invalid", 0.0, true},
		{"25", 25.0, false},
		{" 0.5 ", 0.5, false},
	}

	for _, test := range tests {
		value, err := parseValue(test.input)

		if test.hasError {
			assert.Error(t, err, "Input '%s' should cause an error", test.input)
		} else {
			assert.NoError(t, err, "Input '%s' should not cause an error", test.input)
			assert.Equal(t, test.expected, value, "Input '%s' should parse to expected value", test.input)
		}
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"2025-04-04", time.Date(2025, 4, 4, 0, 0, 0, 0, time.UTC), false},
		{"2023-12-31", time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), false},
		{"04-04-2025", time.Time{}, true}, // Wrong format
		{"invalid", time.Time{}, true},
	}

	for _, test := range tests {
		date, err := parseDate(test.input)

		if test.hasError {
			assert.Error(t, err, "Input '%s' should cause an error", test.input)
		} else {
			assert.NoError(t, err, "Input '%s' should not cause an error", test.input)
			assert.True(t, test.expected.Equal(date), "Input '%s' should parse to expected date", test.input)
		}
	}
}
