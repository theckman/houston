// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, you can
// obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2017 Tim Heckman

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
	// A 34 character string that uniquely identifies this address.
	SID string `json:"sid"`

	// The unique id of the Account responsible for this address.
	AccountSID string `json:"account_sid"`

	// A human-readable description of the address. Maximum 64 characters.
	FriendlyName string `json:"friendly_name"`

	// Your name or business name, or that of your customer.
	CustomerName string `json:"customer_name"`

	// The number and street address where you or your customer is located.
	Street string `json:"street"`

	// The city in which you or your customer is located.
	City string `json:"city"`

	// The state or region in which you or your customer is located.
	Region string `json:"region"`

	// The postal code in which you or your customer is located.
	PostalCode string `json:"postal_code"`

	// The ISO country code of your or your customer's address.
	IsoCountry string `json:"iso_country"`

	// The URI for this resource, relative to https://api.twilio.com.
	URI string `json:"uri"`

	// This is a value that indicates if emergency calling has been enabled on
	// this number. Possible values are true or false.
	EmergencyEnabled bool `json:"emergency_enabled"`

	// In some countries, addresses are validated to comply with local
	// regulation. In those countries, if the address you provide does not pass
	// validation, it will not be accepted as an Address. This value will be
	// true if the Address has been validated, or false for countries that don't
	// require validation or if the Address is non-compliant.
	Validated bool `json:"validated"`
}
