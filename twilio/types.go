// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, you can
// obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2017 Tim Heckman

package twilio

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// HTTPClientInterface is an interface that should be satisfied by the HTTP
// client in net/http. This interface is used for the HTTP client in our Twilio
// client, which allows us to provide an alternative HTTP client implementation
// if we desire.
type HTTPClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

// Time is our own representation of time.Time type because Twilio does not use
// the time format most commonly seen in JSON/JavaScript (RFC-3339/ISO8601).
// Instead Twilio uses the RFC-2822 (RFC-1123) time format, which means the JSON
// package is unable to parse the JSON data in to a time.Time value
// automagically. As such, this type implements the json.Marshaler and
// json.Unmarshaler interfaces to allow Marhsaling and Unmarshaling the time
// format expected by Twilio.
//
// This type also provides a Time() method for returning a time.Time value for
// use in other places that expect an instance of time.Time.
type Time time.Time

// Time is amethod to convert a twilio.Time value in to a time.Time value. The
// inverse can be done with `twilio.Time(t)`, where t is a time.Time value.
func (t Time) Time() time.Time { return time.Time(t) }

// UnmarshalJSON implements the json.Unmarshaler interface. This expects the
// time string to be in an RFC-2822 (RFC-1123) compliant format.
func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "null" {
		*t = Time(time.Time{})
		return nil
	}

	parsedTime, err := time.Parse(time.RFC1123Z, s)

	if err != nil {
		return err
	}

	*t = Time(parsedTime)

	return nil
}

// MarshalJSON implements the json.Marshaler interface. This marshals the time
// in to an RFC-2822 (RFC-1123) compliant format.
func (t Time) MarshalJSON() ([]byte, error) {
	tm := t.Time()

	// RFC2822 says the value for the year field is in the format of 4*DIGIT
	// where DIGIT is any character [0-9]. It also says:
	// The year is any numeric year 1900 or later.
	if year := tm.Year(); year < 1900 || year > 9999 {
		return nil, errors.New("twilio.Time.MarshalJSON(): year outside of range [1900,9999]")
	}

	b := make([]byte, 0, len(time.RFC1123Z)+2)
	b = append(b, '"')
	b = tm.AppendFormat(b, time.RFC1123Z)
	b = append(b, '"')

	return b, nil
}
