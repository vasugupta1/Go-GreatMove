package models

import (
	"time"
)

type RightMoveSearchResponse struct {
	ResultCount string     `json:"resultCount"`
	Properties  []Property `json:"properties"`
}

type Property struct {
	ID                          *int            `json:"id,omitempty" bson:"_id,omitempty"`
	Bedrooms                    *int            `json:"bedrooms,omitempty"`
	Bathrooms                   *int            `json:"bathrooms,omitempty"`
	NumberOfImages              *int            `json:"numberOfImages,omitempty"`
	Summary                     *string         `json:"summary,omitempty"`
	DisplayAddress              *string         `json:"displayAddress,omitempty"`
	Location                    *Location       `json:"location,omitempty"`
	PropertyImages              *PropertyImages `json:"propertyImages,omitempty"`
	Price                       *Price          `json:"price,omitempty"`
	TransactionType             *string         `json:"transactionType,omitempty"`
	PropertyTypeFullDescription *string         `json:"propertyTypeFullDescription,omitempty"`
	Heading                     *string         `json:"heading,omitempty"`
	PropertyUrl                 *string         `json:"propertyUrl,omitempty"`
}

type Location struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type Image struct {
	SrcUrl  *string `json:"srcUrl,omitempty"`
	Url     string  `json:"url"`
	Caption *string `json:"caption,omitempty"`
	Summary *string // will be used for llm text generation
}

type PropertyImages struct {
	Images          []Image `json:"images"`
	MainImageSrc    *string `json:"mainImageSrc,omitempty"`
	MainMapImageSrc *string `json:"mainMapImageSrc,omitempty"`
}

type ListingUpdate struct {
	ListingUpdateReason *string    `json:"listingUpdateReason,omitempty"`
	ListingUpdateDate   *time.Time `json:"listingUpdateDate,omitempty"`
}

type Price struct {
	Amount        *float64          `json:"amount,omitempty"`
	Frequency     *string           `json:"frequency,omitempty"`
	CurrencyCode  *string           `json:"currencyCode,omitempty"`
	DisplayPrices *[]map[string]any `json:"displayPrices,omitempty"`
}
