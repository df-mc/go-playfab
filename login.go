package playfab

import (
	"errors"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

// LoginConfig represents a base structure used for signing in to the PlayFab through a method.
// An implementation of IdentityProvider may embed LoginConfig with additional fields required
// to fill in.
type LoginConfig struct {
	Title title.Title `json:"TitleId,omitempty"`
	// CreateAccount specifies to automatically create a PlayFab account if
	// one is not currently linked to the ID.
	CreateAccount bool `json:"CreateAccount,omitempty"`
	// CustomTags is the optional custom tags associated with the request.
	CustomTags map[string]any `json:"CustomTags,omitempty"`
	// EncryptedRequest is a base64-encoded body that is encrypted with the
	// Title's public RSA key (Enterprise Only).
	EncryptedRequest  []byte             `json:"EncryptedRequest,omitempty"`
	RequestParameters *RequestParameters `json:"InfoRequestParameters,omitempty"`
	// PlayerSecret that is used to verify API request signatures (Enterprise Only).
	PlayerSecret string `json:"PlayerSecret,omitempty"`
}

// IdentityProvider implements a Login method, which signs in to the PlayFab using [LoginConfig.Login] with a LoginConfig
// with additional fields required to fill in, through the path '/Client/LoginWithXXX', where X is normally the method to
// sign in. IdentityProvider is implemented by several platforms which supports signing in to the PlayFab with their identity.
type IdentityProvider interface {
	Login(config LoginConfig) (*Identity, error)
}

// RequestParameters is a set of requested parameters included in ResponseParameters, which can be retrieved
// through [Identity.ResponseParameters]. Users may set RequestParameters as a part of LoginConfig to include
// additional parameters while signing in to PlayFab.
type RequestParameters struct {
	CharacterInventories bool               `json:"GetCharacterInventories,omitempty"`
	CharacterList        bool               `json:"GetCharacterList,omitempty"`
	PlayerProfile        bool               `json:"GetPlayerProfile,omitempty"`
	PlayerStatistics     bool               `json:"GetPlayerStatistics,omitempty"`
	TitleData            bool               `json:"GetTitleData,omitempty"`
	UserAccountInfo      bool               `json:"GetUserAccountInfo,omitempty"`
	UserData             bool               `json:"GetUserData,omitempty"`
	UserInventory        bool               `json:"GetUserInventory,omitempty"`
	UserReadOnlyData     bool               `json:"GetUserReadOnlyData,omitempty"`
	UserVirtualCurrency  bool               `json:"GetUserVirtualCurrency,omitempty"`
	PlayerStatisticNames []string           `json:"PlayerStatisticNames,omitempty"`
	ProfileConstraints   ProfileConstraints `json:"ProfileConstraints,omitempty"`
	TitleDataKeys        []string           `json:"TitleDataKeys,omitempty"`
	UserDataKeys         []string           `json:"UserDataKeys,omitempty"`
	UserReadOnlyDataKeys []string           `json:"UserReadOnlyDataKeys,omitempty"`
}

// ProfileConstraints specifies the properties to return from the player profile, it is included as
// [RequestParameters.ProfileConstraints] to request some of the properties specified on the fields
// as [ResponseParameters.PlayerProfile].
type ProfileConstraints struct {
	ShowAvatarURL                     bool `json:"ShowAvatarUrl,omitempty"`
	ShowBannedUntil                   bool `json:"ShowBannedUntil,omitempty"`
	ShowCampaignAttributions          bool `json:"ShowCampaignAttributions,omitempty"`
	ShowContactEmailAddresses         bool `json:"ShowContactEmailAddresses,omitempty"`
	ShowCreated                       bool `json:"ShowCreated,omitempty"`
	ShowDisplayName                   bool `json:"ShowDisplayName,omitempty"`
	ShowExperimentVariants            bool `json:"ShowExperimentVariants,omitempty"`
	ShowLastLogin                     bool `json:"ShowLastLogin,omitempty"`
	ShowLinkedAccounts                bool `json:"ShowLinkedAccounts,omitempty"`
	ShowLocations                     bool `json:"ShowLocations,omitempty"`
	ShowMemberships                   bool `json:"ShowMemberships,omitempty"`
	ShowOrigination                   bool `json:"ShowOrigination,omitempty"`
	ShowPushNotificationRegistrations bool `json:"ShowPushNotificationRegistrations,omitempty"`
	ShowStatistics                    bool `json:"ShowStatistics,omitempty"`
	ShowTags                          bool `json:"ShowTags,omitempty"`
	ShowTotalValueToDateInUSD         bool `json:"ShowTotalValueToDateInUsd,omitempty"`
	ShowValuesToDate                  bool `json:"ShowValuesToDate,omitempty"`
}

// Login signs in to PlayFab using the request body and the path. The path normally follows the format
// 'Client/LoginWithXXX' where X typically is the method for signing in. The request body is generally
// a structure that may embed LoginConfig with additional fields required for the specific path, such
// as a token of IdentityProvider.
func (l LoginConfig) Login(path string, body any) (*Identity, error) {
	if l.Title == "" {
		return nil, errors.New("playfab: LoginConfig: Title not set")
	}
	return internal.Post[*Identity](l.Title.URL().JoinPath(path), body)
}
