package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SNBInterestRate represents a Swiss National Bank interest rate data point
type SNBInterestRate struct {
	Code        string    `json:"code"`
	Value       float64   `json:"value"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
}

// SNBScraper implements the Scraper interface for Swiss National Bank interest rates
type SNBScraper struct {
	rssURL     string
	httpClient *http.Client
}

// NewSNBScraper creates a new SNB scraper instance
func NewSNBScraper() *SNBScraper {
	return &SNBScraper{
		rssURL:     "https://www.snb.ch/public/en/rss/interestRates",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns the unique identifier for this scraper
func (s *SNBScraper) Name() string {
	return "snb_interest_rates"
}

// Schedule returns the recommended scraping interval
func (s *SNBScraper) Schedule() time.Duration {
	// SNB typically updates rates daily or on business days
	return 6 * time.Hour
}

// Validate checks if the scraper configuration is valid
func (s *SNBScraper) Validate(ctx context.Context) error {
	if s.rssURL == "" {
		return fmt.Errorf("RSS URL is required")
	}
	return nil
}

// Init performs any necessary initialization
func (s *SNBScraper) Init(ctx context.Context) error {
	// No specific initialization needed
	return nil
}

// RSS feed structures
type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Items []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Country     string `xml:"country"`
	Code        string `xml:"code"`
	Value       string `xml:"value"`
	Unit        string `xml:"unit"`
	Date        string `xml:"date"`
}

// Scrape performs the data collection process for SNB interest rates
func (s *SNBScraper) Scrape(ctx context.Context) ([]Result, error) {
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SNB RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse XML
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Process items
	var rates []SNBInterestRate
	for _, item := range feed.Channel.Items {
		// Skip items without a code or value
		if item.Code == "" || item.Value == "" {
			continue
		}

		// Parse value
		value, err := parseValue(item.Value)
		if err != nil {
			// Log and skip invalid values
			continue
		}

		// Parse date
		date, err := parseDate(item.Date)
		if err != nil {
			// Use publication date as fallback
			date, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				// Use current time if all parsing fails
				date = time.Now()
			}
		}

		rate := SNBInterestRate{
			Code:        item.Code,
			Value:       value,
			Date:        date,
			Description: item.Description,
			Unit:        item.Unit,
		}

		rates = append(rates, rate)
	}

	// Create result
	result := Result{
		Source:    s.Name(),
		Timestamp: time.Now(),
		Data:      rates,
		Metadata: map[string]string{
			"url": s.rssURL,
		},
	}

	return []Result{result}, nil
}

// parseValue parses a string value to float64
func parseValue(val string) (float64, error) {
	// Remove any non-numeric characters except for decimal point
	val = strings.TrimSpace(val)
	if strings.Contains(val, " ") {
		val = strings.Split(val, " ")[0]
	}

	return strconv.ParseFloat(val, 64)
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
