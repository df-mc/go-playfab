package playfab

import (
	"encoding/json"
	"time"

	"github.com/df-mc/go-playfab/entity"
)

// Identity a session identity that can subsequently be used for API which requires an authentication.
// It is generally returned as a result of [IdentityProvider.Login] or [LoginConfig.Login].
type Identity struct {
	// EntityToken is an [entity.Token] of [entity.TypeTitlePlayerAccount].
	// API requests will mostly require an [entity.Token] of [entity.TypeMasterPlayerAccount]
	// so you may exchange it with [entity.Token.Exchange] with PlayFabID.
	EntityToken *entity.Token `json:"EntityToken,omitempty"`
	// ResponseParameters is a set of parameters requested to be included in [LoginConfig.RequestParameters].
	// It includes an additional data specific to the player/entity signed in.
	ResponseParameters json.RawMessage `json:"InfoResultPayload,omitempty"`
	// LastLoginTime is the time of previous login. If there was no previous login, it is zero [time.Time].
	LastLoginTime time.Time `json:"LastLoginTime,omitempty"`
	// NewlyCreated is true if the account was newly created on login.
	NewlyCreated bool `json:"NewlyCreated,omitempty"`
	// PlayFabID is the unique ID of player. It can be used for exchanging an [entity.Token]
	// of [entity.TypeMasterPlayerAccount] with [entity.Token.Exchange].
	PlayFabID string `json:"PlayFabId,omitempty"`
	// SessionTicket is a unique token authorizing the user and game at server level, for the
	// current session. In Minecraft, it is used for authorizing with franchise API using a PlayFab token.
	SessionTicket string `json:"SessionTicket,omitempty"`
	// SettingsForUser is the settings specific to the user.
	SettingsForUser json.RawMessage `json:"SettingsForUser,omitempty"`
	// The experimentation TreatmentAssignment for this user at the time of login.
	TreatmentAssignment json.RawMessage `json:"TreatmentAssignment,omitempty"`
}
