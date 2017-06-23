package twilio

// The Exception struct is the representation of an error resource returned from
// the Twilio API when an error happens with a request.
type Exception struct {
	Status   int    `json:"Status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

// An Account instance resource represents a single Twilio account.
type Account struct {
	// A 34 character string that uniquely identifies this account.
	SID string `json:"sid"`

	// The date this account was created.
	DateCreated Time `json:"date_created"`

	// The date that this account was last updated.
	DateUpdated Time `json:"date_updated"`

	// A human readable description of this account, up to 64 characters long.
	// By default the FriendlyName is your email address.
	FriendlyName string `json:"friendly_name"`

	// The type of this account. Either Trial or Full if you've upgraded.
	Type string `json:"type"`

	// The status of this account. Usually active, but can be suspended or
	// closed.
	Status string `json:"status"`

	// The authorization token for this account. This token should be kept a
	// secret, so no sharing.
	AuthToken string `json:"auth_token,omitempty"`

	// The URI for this resource, relative to https://api.twilio.com.
	URI string `json:"uri"`

	// The list of subresources under this account.
	SubresourceURIs map[string]string `json:"subresource_uris"`

	// The Sid of the parent account for this account. The OwnerAccountSid of a
	// parent account is its own sid.
	OwnerAccountSID string `json:"owner_account_sid"`
}

// An Address instance resource represents a single Twilio address.
//
// From the docs: An Address instance resource represents your or your
// customer’s physical location within a country.
type Address struct {
	SID              string `json:"sid"`
	AccountSID       string `json:"account_sid"`
	FriendlyName     string `json:"friendly_name"`
	CustomerName     string `json:"customer_name"`
	Street           string `json:"street"`
	City             string `json:"city"`
	Region           string `json:"region"`
	PostalCode       string `json:"postal_code"`
	IsoCountry       string `json:"iso_country"`
	URI              string `json:"uri"`
	EmergencyEnabled bool   `json:"emergency_enabled"`
	Validated        bool   `json:"validated"`
}
