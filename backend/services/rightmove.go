package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/greatmove/backend/models"
)

type RightMoveService interface {
	GetLocationIdentifiers(query string) ([]string, error)
	SearchProperties(locationId string) ([]models.Property, error)
}

type RightMove struct {
	HTTPClient HttpClient
}

const (
	ResultPerPage = 24
	MaxResults    = 1000
)

func ConstructRightMove(client HttpClient) *RightMove {
	return &RightMove{
		HTTPClient: client,
	}
}

func (fl *RightMove) GetLocationIdentifiers(query string) ([]string, error) {
	upper := strings.ToUpper(query)
	var builder strings.Builder
	for i, char := range upper {
		builder.WriteRune(char)
		if (i+1)%2 == 0 && i < len(upper)-1 {
			builder.WriteRune('/')
		}
	}
	tokenziedQuery := builder.String()
	token := strings.TrimSuffix(tokenziedQuery, "/")
	url := fmt.Sprintf("https://www.rightmove.co.uk/typeAhead/uknostreet/%s/", token)
	resp, err := fl.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	locations := extractLocations(bodyBytes)
	return locations, nil
}

func (fl *RightMove) SearchProperties(locationId string) ([]models.Property, error) {
	var offset = 0
	firstPageUrl := constructSearchUrl(offset, locationId, ResultPerPage)
	searchResponse, err := getAndParseJsonResponse(firstPageUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching properties: %w", err)
	}
	var totalResult []models.Property
	totalResult = append(totalResult, searchResponse.Properties...)
	for offset := ResultPerPage; offset < MaxResults; offset += ResultPerPage {
		searchUrl := constructSearchUrl(offset, locationId, ResultPerPage)
		searchResponse, err := getAndParseJsonResponse(searchUrl)
		if err != nil {
			return nil, fmt.Errorf("error fetching properties: %w", err)
		}
		if len(searchResponse.Properties) == 0 {
			break
		}
		totalResult = append(totalResult, searchResponse.Properties...)
	}

	return totalResult, nil
}

func getAndParseJsonResponse(url string) (*models.RightMoveSearchResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response models.RightMoveSearchResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}
	return &response, nil
}

func constructSearchUrl(offset int, locationId string, resultPerPage int) string {
	baseUrl := "https://www.rightmove.co.uk/api/_search"
	params := url.Values{}
	params.Set("areaSizeUnit", "sqft")
	params.Set("channel", "BUY")
	params.Set("currencyCode", "GBP")
	params.Set("includeSSTC", "false")
	params.Set("index", strconv.Itoa(offset))
	params.Set("isFetching", "false")
	params.Set("locationIdentifier", locationId)
	params.Set("numberOfPropertiesPerPage", strconv.Itoa(resultPerPage))
	params.Set("radius", "0.0")
	params.Set("sortType", "6")
	params.Set("viewType", "LIST")
	return fmt.Sprintf("%s?%s", baseUrl, params.Encode())
}

func extractLocations(jsonResponse []byte) []string {
	var response models.RightMoveLocationResponse
	err := json.Unmarshal(jsonResponse, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON response:", err)
		return nil
	}
	var locations []string
	for _, location := range response.TypeAheadLocations {
		locations = append(locations, location.LocationIdentifier)
	}
	return locations
}
