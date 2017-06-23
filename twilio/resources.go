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
	SID             string            `json:"sid"`
	DateCreated     Time              `json:"date_created"`
	DateUpdated     Time              `json:"date_updated"`
	FriendlyName    string            `json:"friendly_name"`
	Type            string            `json:"type"`
	Status          string            `json:"status"`
	AuthToken       string            `json:"auth_token,omitempty"`
	URI             string            `json:"uri"`
	SubresourceURIs map[string]string `json:"subresource_uris"`
	OwnerAccountSID string            `json:"owner_account_sid"`
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
