package catalog

import (
	"encoding/json"
	"github.com/df-mc/go-playfab/entity"
	"time"
)

type Item struct {
	// AlternateIDs is the alternate IDs associated with the Item. An alternate
	// ID can be set to 'FriendlyId' or any of the supported marketplace names.
	AlternateIDs []AlternateID `json:"AlternateIds,omitempty"`
	// ContentType is the client-defined type of the Item.
	ContentType string `json:"ContentType,omitempty"`
	// Contents is the set of content/files associated with the Item. Up to 100
	// files can be added to an Item. In Minecraft, it includes a set of URL for 'XForge'
	// where it contains a ZIP file containing a set of encrypted packs.
	Contents []Content `json:"Contents,omitempty"`
	// CreationDate is the date and time when the Item was created.
	CreationDate time.Time `json:"CreationDate,omitempty"`
	// CreatorEntity is the [entity.Key] of the creator of the Item.
	CreatorEntity entity.Key `json:"CreatorEntity,omitempty"`
	// DeepLinks is the set of platform specific deep links for the Item.
	DeepLinks []DeepLink `json:"DeepLinks,omitempty"`
	// DefaultStackID is the stack ID that will be used as default for the Item
	// in inventory when an explicit one is not provided. The DefaultStackID can be
	// a static stack ID or '{GUID}', which will generate a unique stack ID for the
	// Item. If empty, inventory's default stack ID will be used.
	DefaultStackID string `json:"DefaultStackId,omitempty"`
	// Description is a Dictionary of localized descriptions. Descriptions have
	// a 10000-character limit per country code.
	Description Dictionary[string] `json:"Description,omitempty"`
	// DisplayProperties is a game-specific properties for display purposes. It is
	// an arbitrary JSON blob. The fields of DisplayProperties has a 10000-byte
	// limit per Item. In Minecraft, it contains the name of creator, whether
	// the Item is purchasable, a URL of video trailer, prices, address and port of
	// server (If ContentType is '3PP' or '3PP_V2.0').
	DisplayProperties map[string]json.RawMessage `json:"DisplayProperties,omitempty"`
	// DisplayVersion is the user-provided version of the Item for display
	// purposes. It has a maximum character length of 50.
	DisplayVersion string `json:"DisplayVersion,omitempty"`
	// ETag is the current ETag value that can be used for optimistic
	// concurrency in the 'If-None-Match' header.
	ETag string `json:"ETag,omitempty"`
	// EndDate is the date of when the Item will cease to be available. If left
	// a zero [time.Time] then the product will be available indefinitely.
	EndDate time.Time `json:"EndDate,omitempty"`
	// ID is the unique ID of the Item. It can be specified to [Query.ID].
	ID string `json:"Id,omitempty"`
	// Images is the images associated with the Item. Images can be thumbnails
	// or screenshots. Up to 100 images can be added to an Item. Only .png, .jpg,
	// .gif, and .bmp file types can be uploaded.
	Images []Image `json:"Images,omitempty"`
	// Hidden indicates if the Item is hidden.
	Hidden bool `json:"IsHidden,omitempty"`
	// ItemReferences is the item references associated with the Item. Every Item
	// can have up to 50 item references.
	ItemReferences []ItemReference `json:"ItemReferences,omitempty"`
	// Keywords is a Dictionary of localized keywords. Keywords have a 50-character
	// limit per keyword and up to 32 keywords can be added per country code.
	Keywords Dictionary[*Keyword] `json:"Keywords,omitempty"`
	// LastModifiedDate is the date and time the Item was last updated.
	LastModifiedDate time.Time `json:"LastModifiedDate,omitempty"`
	// Moderation is the moderation state for the Item.
	Moderation ModerationState `json:"Moderation,omitempty"`
	// Platforms is the platforms supported by the Item.
	Platforms []string `json:"Platforms,omitempty"`
	// PriceOptions is the prices the Item can be purchased for.
	PriceOptions PriceOptions `json:"PriceOptions,omitempty"`
	// Rating s the rating summary for the Item.
	Rating Rating `json:"Rating,omitempty"`
	// StartDate is the date of when the Item will be available. If left as
	// a zero [time.Time] then the product will appear immediately.
	StartDate time.Time `json:"StartDate,omitempty"`
	// StoreDetails is an optional details for stores items.
	StoreDetails StoreDetails `json:"StoreDetails,omitempty"`
	// Tags is the list of tags that are associated with the Item. Up to 32 tags
	// can be added to an Item.
	Tags []string `json:"Tags,omitempty"`
	// Title is a Dictionary of localized titles. Titles have a 512-character limit
	// per country code.
	Title Dictionary[string] `json:"Title,omitempty"`
	// Type is the high-level type of the Item. It is one of constants defined below.
	Type string `json:"Type,omitempty"`
}

type StoreReference struct {
	AlternateID AlternateID `json:"AlternateId,omitempty"`
	ID          string      `json:"Id,omitempty"`
}

type AlternateID struct {
	Type  string `json:"Type,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Content struct {
	// ID is the unique ID of the Content.
	ID string `json:"Id,omitempty"`
	// MaxClientVersion is the maximum client version that the Content is
	// compatible with. Client Versions can be up to 3 segments separated
	// by periods (.) and each segment can have a maximum value of 65535.
	MaxClientVersion string `json:"MaxClientVersion,omitempty"`
	// MinClientVersion is the minimum client version that the Content is
	// compatible with. Client Versions can be up to 3 segments separated
	// by periods (.) and each segment can have a maximum value of 65535.
	MinClientVersion string `json:"MinClientVersion,omitempty"`
	// Tags is the list of tags that are associated with the Content. Tags
	// must be defined in the Catalog Config before being used in Content.
	Tags []string `json:"Tags,omitempty"`
	// Type is the client-defined type of the Content. Types must be defined
	// in the Catalog Config before being used.
	Type string `json:"Type,omitempty"`
	// URL is the Azure CDN URL for retrieval of the Item binary content.
	// In Minecraft (and some other games), It is a URL for XForge asset.
	URL string `json:"Url,omitempty"`
}

type DeepLink struct {
	// Platform is the target platform for the DeepLink.
	Platform string `json:"Platform,omitempty"`
	// URL is the deep link for the Platform.
	URL string `json:"Url,omitempty"`
}

type Image struct {
	// ID is the unique ID of the Image.
	ID string `json:"Id,omitempty"`
	// Tag is the client-defined tag associated with the Image. Tags must be
	// in the Catalog Config before being used in Image.
	Tag string `json:"Tag,omitempty"`
	// Type is the type of the Image. It is one of constants defined below.
	// There can only be one Image of ImageTypeThumbnail per Item.
	Type string `json:"Type,omitempty"`
	// URL is the URL for retrieval of the Image.
	URL string `json:"Url,omitempty"`
}

const (
	ImageTypeThumbnail  = "thumbnail"
	ImageTypeScreenshot = "screenshot"
)

type ItemReference struct {
	// Amount is the amount of the catalog Item.
	Amount int `json:"Amount,omitempty"`
	// ID is the unique ID of the catalog Item.
	ID string `json:"Id,omitempty"`
	// PriceOptions is the prices that the Item referenced in the
	// ID can be purchased for.
	PriceOptions PriceOptions `json:"PriceOptions,omitempty"`
}

type PriceOptions []Price

func (opts PriceOptions) MarshalJSON() ([]byte, error) {
	type raw struct {
		Prices []Price `json:"Prices,omitempty"`
	}
	return json.Marshal(raw{Prices: opts})
}

func (opts *PriceOptions) UnmarshalJSON(b []byte) error {
	var raw struct {
		Prices []Price `json:"Prices,omitempty"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	*opts = raw.Prices
	return nil
}

type Price struct {
	// Amounts is the amounts of the Price. Each Price can have up to
	// 15 amounts.
	Amounts []PriceAmount `json:"Amounts,omitempty"`
	// UnitAmount is the per-unit amount the Price can be used to purchase.
	UnitAmount int `json:"UnitAmount,omitempty"`
	// UnitDurationInSeconds is the per-unit duration the Price can be used
	// to purchase. The maximum duration is 100 years.
	UnitDurationInSeconds int `json:"UnitDurationInSeconds,omitempty"`
}

type PriceAmount struct {
	Amount int    `json:"Amount,omitempty"`
	ItemID string `json:"ItemId,omitempty"`
}

type Keyword []string

func (k *Keyword) MarshalJSON() ([]byte, error) {
	type raw struct {
		Values []string `json:"Values,omitempty"`
	}
	return json.Marshal(raw{Values: *k})
}

func (k *Keyword) UnmarshalJSON(b []byte) error {
	var raw struct {
		Values []string `json:"Values,omitempty"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	*k = raw.Values
	return nil
}

type ModerationState struct {
	// LastModifiedDate is the date and time the ModerationState was last updated.
	LastModifiedDate time.Time `json:"LastModifiedDate,omitempty"`
	// Reason is the current stated reason for the associated Item being moderated.
	Reason string `json:"Reason,omitempty"`
	// Status is the current moderation status for the associated Item.
	Status string `json:"Status,omitempty"`
}

const (
	ModerationStatusApproved           string = "Approved"
	ModerationStatusAwaitingModeration string = "AwaitingModeration"
	ModerationStatusRejected           string = "Rejected"
	ModerationStatusUnknown            string = "Unknown"
)

type Rating struct {
	// Average is the average rating for the Item.
	Average float32 `json:"Average,omitempty"`
	// Count1Star is the total count of 1-star ratings for the Item.
	Count1Star int `json:"Count1Star,omitempty"`
	// Count2Star is the total count of 2-star ratings for the Item.
	Count2Star int `json:"Count2Star,omitempty"`
	// Count3Star is the total count of 3-star ratings for the Item.
	Count3Star int `json:"Count3Star,omitempty"`
	// Count4Star is the total count of 4-star ratings for the Item.
	Count4Star int `json:"Count4Star,omitempty"`
	// Count5Star is the total count of 5-star ratings for the Item.
	Count5Star int `json:"Count5Star,omitempty"`
	// TotalCount is the total count of ratings for the Item.
	TotalCount int `json:"TotalCount,omitempty"`
}

type StoreDetails struct {
	// FilterOptions is the options for the filter in filter-based stores.
	// There options are mutually exclusive with item references.
	FilterOptions FilterOptions `json:"FilterOptions,omitempty"`
	// PriceOptionsOverride is the global prices utilized in the store. These
	// options are mutually exclusive with price options in ItemReference.
	PriceOptionsOverride PriceOptionsOverride `json:"PriceOptionsOverride,omitempty"`
}

type FilterOptions struct {
	Filter          string `json:"Filter,omitempty"`
	IncludeAllItems bool   `json:"IncludeAllItems,omitempty"`
}

type PriceOptionsOverride []PriceOverride

func (opts PriceOptionsOverride) MarshalJSON() ([]byte, error) {
	type raw struct {
		Prices []PriceOverride `json:"Prices,omitempty"`
	}
	return json.Marshal(raw{Prices: opts})
}

func (opts *PriceOptionsOverride) UnmarshalJSON(b []byte) error {
	var raw struct {
		Prices []PriceOverride `json:"Prices,omitempty"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	*opts = raw.Prices
	return nil
}

type PriceOverride struct {
	// Amounts is the currency amounts utilized in the override for a singular price.
	Amounts []PriceAmountOverride `json:"Amounts,omitempty"`
}

type PriceAmountOverride struct {
	// FixedValue is the exact value that should be utilized in the PriceAmountOverride.
	FixedValue int `json:"FixedValue,omitempty"`
	// ItemID is the ID of the Item the PriceAmountOverride should utilize.
	ItemID string `json:"ItemId,omitempty"`
	// Multiplier is the multiplier that will be applied to the base catalog value
	// to determine what value should be utilized in the PriceAmountOverride.
	Multiplier int `json:"Multiplier,omitempty"`
}

const (
	ItemTypeBundle      = "bundle"
	ItemTypeCatalogItem = "catalogItem"
	ItemTypeCurrency    = "currency"
	ItemTypeStore       = "store"
	ItemTypeUGC         = "ugc"
)
