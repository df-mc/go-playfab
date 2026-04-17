package catalog

import (
	"encoding/json"
	"time"

	"github.com/df-mc/go-playfab/entity"
	"github.com/google/uuid"
)

// Item represents a catalog item in the PlayFab Economy v2 catalog.
// It contains metadata, localized content, pricing, and other properties
// associated with a purchasable or displayable entity in the catalog.
//
// See: https://learn.microsoft.com/en-us/rest/api/playfab/economy/catalog/get-item?view=playfab-rest#catalogitem
type Item struct {
	// AlternateIDs is the list of alternate IDs associated with this item.
	AlternateIDs []AlternateID `json:"AlternateIds"`
	// ContentType is the title-defined type of the item.
	ContentType string
	// Contents is the set of content/files associated with this item.
	// Up to 100 files can be added to an item.
	Contents []Content
	// CreatorEntity is the entity key identifying the creator of this catalog item.
	CreatorEntity entity.Key
	// DeepLinks is the set of platform-specific deep links for this item.
	DeepLinks []DeepLink
	// DefaultStackID is the stack ID used as default for this item in Inventory
	// when an explicit one is not provided. It can be a static stack ID or '{guid}',
	// which generates a unique stack ID for the item. If empty, Inventory's
	// default stack ID will be used.
	DefaultStackID string

	// Title is a dictionary of localized titles for this item.
	// Key is a language code and the value is the localized string.
	// Each title has a 512 character limit per locale.
	Title Dictionary[string]
	// Description is a dictionary of localized descriptions for this item.
	// Key is a language code and the value is the localized string.
	// Each description has a 10000 character limit per locale.
	Description Dictionary[string]
	// DisplayProperties contains game-specific properties for display purposes.
	// This is an arbitrary JSON blob with a 10000 byte limit per item.
	DisplayProperties json.RawMessage
	// DisplayVersion is the user-provided version of the item for display purposes.
	// Maximum character length is 50.
	DisplayVersion string
	// ETag is the current ETag value that can be used for optimistic concurrency
	// in the If-None-Match header.
	ETag string

	// CreationDate is the date and time when this item was created.
	CreationDate time.Time
	// StartDate is the date when the item will become available.
	// If not provided, the item will appear immediately in the catalog.
	StartDate time.Time
	// EndDate is the date when the item will cease to be available.
	// If not provided, the item will be available indefinitely.
	EndDate time.Time
	// LastModifiedDate is the date and time when this item was last updated.
	LastModifiedDate time.Time

	// ID is the unique ID of the item.
	ID uuid.UUID `json:"Id"`

	// Images is the set of images associated with this item.
	// Images can be thumbnails or screenshots.
	// Up to 100 images can be added to an item.
	// Only .png, .jpg, .gif, and .bmp file types can be uploaded.
	Images []Image
	// Hidden indicates whether the item is currently hidden from the catalog.
	Hidden bool `json:"IsHidden"`
	// ItemReferences is the list of item references associated with this item,
	// such as the items contained in a Bundle, Store, or Subscription.
	// Every item can have up to 50 item references.
	ItemReferences []ItemReference

	// Keywords is a dictionary of localized keywords associated with this item.
	// Key is a language code and the value is the localized list of keywords.
	// Keywords have a 50 character limit each, and up to 32 keywords can be added per locale.
	Keywords Dictionary[KeywordSet]

	// Moderation is the moderation state for this item.
	// It is typically used for community-provided (UGC) items.
	Moderation ModerationState
	// Platforms lists the platforms supported by this item.
	Platforms []string

	// PriceOptions contains the prices this item can be purchased for.
	// An item can have up to 15 prices.
	PriceOptions PriceOptions

	// Rating is the rating summary for this item.
	Rating Rating
	// RealMoneyPrices contains the real-money prices for this item, scoped per
	// marketplace platform. Each entry is a map from ISO 4217 currency code to
	// the price in the smallest currency unit (e.g. cents). Currently, only USD is supported.
	RealMoneyPrices RealMoneyPrices `json:"RealMoneyPriceDetails"`

	// Tags is the list of tags associated with this item.
	// Up to 32 tags can be added to an item.
	Tags []string
	// Type is the high-level type of the item.
	// It can be one of the constants prefixed with ItemType* defined below.
	Type string
}

const (
	// ItemTypeCatalogItem is the item type for standard catalog items available for purchase.
	ItemTypeCatalogItem = "catalogItem"
	// ItemTypeCurrency is the item type representing an in-game currency.
	// e.g. in Minecraft, this is used to represent Minecoins.
	ItemTypeCurrency = "currency"
	// ItemTypeStore is the item type for a store.
	// A store holds a list of items and prices, and can override the base catalog prices for those items.
	ItemTypeStore = "store"
	// ItemTypeUserGeneratedContent is the item type for user-generated content (UGC).
	ItemTypeUserGeneratedContent = "ugc"
	// ItemTypeSubscription is the item type for subscription items.
	// e.g. in Minecraft, this is used to represent Realms plans.
	ItemTypeSubscription = "subscription"
)

// AlternateID describes an alternate ID associated with a catalog item.
type AlternateID struct {
	// Type is the type of the alternate ID. It can be 'FriendlyId' or any other marketplace names.
	Type string
	// Value is the value of the alternate ID.
	Value string
}

// Content represents a file or binary content associated with a catalog item.
type Content struct {
	// ID is the unique ID of this content entry.
	ID string `json:"Id"`
	// MaxClientVersion is the maximum client version this content is compatible with.
	// Versions follow semantic versioning with up to 3 dot-separated segments (e.g. "1.2.3"),
	// where each segment can be at most 65535.
	MaxClientVersion string
	// MinClientVersion is the minimum client version this content is compatible with.
	// Versions follow semantic versioning with up to 3 dot-separated segments (e.g. "1.2.3"),
	// where each segment can be at most 65535.
	MinClientVersion string
	// Tags is the list of tags associated with this content.
	Tags []string
	// Type is the title-defined type of this content.
	Type string
	// URL is the Azure CDN URL for retrieval of the binary content.
	URL string `json:"Url"`
}

// DeepLink represents a platform-specific deep link associated with a catalog item.
type DeepLink struct {
	// Platform is the target platform for this deep link.
	Platform string
	// URL is the deep link URL for the target platform.
	URL string `json:"Url"`
}

// Image represents an image associated with a catalog item.
// Images can be defined as either a thumbnail or a screenshot.
// There can only be one thumbnail image per item.
// Only .png, .jpg, .gif, and .bmp file types can be uploaded.
type Image struct {
	// ID is the unique ID of this image.
	ID string `json:"Id"`
	// Tag is the title-defined tag associated with this image.
	Tag string
	// Type indicates whether this image is a thumbnail or a screenshot.
	// It can be one of the constants prefixed with ImageType* defined below.
	Type string
	// URL is the URL for retrieval of this image.
	URL string `json:"Url"`
}

const (
	// ImageTypeThumbnail is the image type for a thumbnail image.
	// Only one thumbnail image is allowed per catalog item.
	ImageTypeThumbnail = "thumbnail"
	// ImageTypeScreenshot is the image type for a screenshot image.
	ImageTypeScreenshot = "screenshot"
)

// ItemReference represents a reference to another catalog item.
// It is used within bundles, stores, and subscriptions to list contained items.
type ItemReference struct {
	// Amount is the quantity of the referenced catalog item.
	Amount int
	// ID is the unique ID of the referenced catalog item.
	ID string `json:"Id"`
	// PriceOptions contains the prices at which the referenced item can be purchased.
	PriceOptions PriceOptions
}

// PriceOptions is the list of prices a catalog item can be purchased for.
// An item can have up to 15 prices.
type PriceOptions []Price

// MarshalJSON implements [json.Marshaler] for PriceOptions,
// encoding it as a JSON object with a "Prices" field as required by the PlayFab API.
func (p *PriceOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Prices []Price
	}{Prices: *p})
}

// UnmarshalJSON implements [json.Unmarshaler] for PriceOptions,
// decoding it from a JSON object with a "Prices" field as returned by the PlayFab API.
func (p *PriceOptions) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &struct {
		Prices *[]Price
	}{Prices: (*[]Price)(p)})
}

// Price represents a single purchasable price for a catalog item.
type Price struct {
	// Amounts is the list of currency amounts required for this price.
	// Each price can have up to 15 item amounts.
	Amounts []PriceAmount
	// UnitAmount is the per-unit quantity this price allows the player to purchase.
	UnitAmount int
	// UnitDurationInSeconds is the per-unit duration this price allows the player to purchase.
	// The maximum duration is 100 years.
	UnitDurationInSeconds int
}

// PriceAmount represents a single currency component of a price.
type PriceAmount struct {
	// Value is the amount of currency required.
	Value int `json:"Amount"`
	// ItemID is the ID of the catalog item used as currency for this amount.
	ItemID uuid.UUID `json:"ItemId"`
}

// KeywordSet is a list of localized keywords associated with a catalog item.
// Keywords have a 50 character limit each, and up to 32 can be added per locale.
type KeywordSet []string

// MarshalJSON implements [json.Marshaler] for KeywordSet,
// encoding it as a JSON object with a "Values" field as required by the PlayFab API.
func (s *KeywordSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Values []string
	}{Values: *s})
}

// UnmarshalJSON implements [json.Unmarshaler] for KeywordSet,
// decoding it from a JSON object with a "Values" field as returned by the PlayFab API.
func (s *KeywordSet) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &struct {
		Values *[]string
	}{Values: (*[]string)(s)}); err != nil {
		return err
	}
	return nil
}

// ModerationState represents the moderation status of a catalog item.
// Moderation is typically applied to community-provided (UGC) items.
type ModerationState struct {
	// LastModifiedDate is the date and time this moderation state was last updated.
	LastModifiedDate time.Time
	// Reason is the stated reason for the item being moderated, if applicable.
	Reason string
	// Status is the current moderation status of the item.
	// It can be one of the constants prefixed with ModerationStatus* defined below.
	Status string
}

const (
	// ModerationStatusUnknown indicates an unknown moderation status.
	ModerationStatusUnknown = "Unknown"
	// ModerationStatusAwaitingModeration indicates the item is pending moderation review.
	ModerationStatusAwaitingModeration = "AwaitingModeration"
	// ModerationStatusApproved indicates the item has been approved by moderation.
	ModerationStatusApproved = "Approved"
	// ModerationStatusRejected indicates the item has been rejected by moderation.
	ModerationStatusRejected = "Rejected"
)

// Rating represents the aggregated rating summary for a catalog item.
type Rating struct {
	// Average is the average star rating for this item.
	Average float32
	// Count1Star is the total number of 1-star ratings for this item.
	Count1Star int
	// Count2Star is the total number of 2-star ratings for this item.
	Count2Star int
	// Count3Star is the total number of 3-star ratings for this item.
	Count3Star int
	// Count4Star is the total number of 4-star ratings for this item.
	Count4Star int
	// Count5Star is the total number of 5-star ratings for this item.
	Count5Star int
	// TotalCount is the total number of ratings submitted for this item.
	TotalCount int
}

// RealMoneyPrices contains the real-money prices for a catalog item across
// multiple marketplace platforms. Each field is a map from ISO 4217 currency code
// to the price expressed in the smallest currency unit (e.g. 139 for $1.39 USD).
// Currently, only United States Dollar (USD) is supported.
type RealMoneyPrices struct {
	// AppleAppStorePrices is the price map for the Apple App Store, keyed by currency code.
	AppleAppStorePrices map[string]int
	// GooglePlayPrices is the price map for Google Play, keyed by currency code.
	GooglePlayPrices map[string]int
	// MicrosoftStorePrices is the price map for the Microsoft Store, keyed by currency code.
	MicrosoftStorePrices map[string]int
	// NintendoEShopPrices is the price map for the Nintendo eShop, keyed by currency code.
	NintendoEShopPrices map[string]int
	// PlayStationStorePrices is the price map for the PlayStation Store, keyed by currency code.
	PlayStationStorePrices map[string]int
	// SteamPrices is the price map for Steam, keyed by currency code.
	SteamPrices map[string]int
}
