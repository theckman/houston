package twilio

import (
	"bytes"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Run("Time.Time()", testTimeTime)
	t.Run("Time.MarshalJSON()", testTimeMarshalJSON)
	t.Run("Time.UnmarshalJSON()", testTimeUnmarshalJSON)
}

func testTimeTime(t *testing.T) {
	var t1, t2 time.Time
	var tt1 Time

	t1 = time.Now()
	tt1 = Time(t1)

	if t2 = tt1.Time(); t1.UnixNano() != t2.UnixNano() {
		t.Errorf("twilio.Time.Time() = %q; want %q", t2.String(), t1.String())
	}
}

func testTimeMarshalJSON(t *testing.T) {
	var jsonOut []byte
	var err error

	tests := []struct {
		in   Time
		out  []byte
		err  bool
		desc string
	}{
		{in: Time(time.Unix(0, 0).UTC()), out: []byte(`"Thu, 01 Jan 1970 00:00:00 +0000"`), desc: `Time should be the UNIX epoch`},
		{in: Time(time.Unix(-2208988801, 0).UTC()), err: true, desc: `Times before 1900 are invalid per RFC-2822`},
	}

	for _, test := range tests {
		jsonOut, err = test.in.MarshalJSON()

		if !test.err && err != nil { // error found when not expected
			t.Errorf("%#v.MarshalJSON() = _, %q; want nil", test.in, err)
			continue
		}

		if test.err && err == nil {
			t.Errorf("%#v.MarshalJSON() = _, nil; want err != nil", test.in)
			continue
		}

		// XX(heckman): we make sure this error, if it exists, is expected above
		// so we can safely move on to the next iteration here...
		if err != nil {
			continue
		}

		if !bytes.Equal(jsonOut, test.out) {
			t.Errorf("\nDescription: %s\n%#v.MarshalJSON() = %#v, nil; want %#v", test.desc, test.in, string(jsonOut), string(test.out))
		}
	}
}

func testTimeUnmarshalJSON(t *testing.T) {
	var tm Time
	var err error

	tests := []struct {
		in   []byte
		out  Time
		err  bool
		desc string
	}{
		{in: []byte(`"Thu, 01 Jan 1970 00:00:00 +0000"`), out: Time(time.Unix(0, 0).UTC()), desc: `Time should be the UNIX epoch`},
		{in: []byte(`null`), out: Time(time.Unix(-62135596800, 0).UTC()), desc: `Time should be the time.Time zero value`},
	}

	for _, test := range tests {
		tm = Time{}
		err = tm.UnmarshalJSON(test.in)

		if !test.err && err != nil { // error found when not expected
			t.Errorf("%#v.MarshalJSON() = _, %q; want nil", test.in, err)
			continue
		}

		if test.err && err == nil {
			t.Errorf("%#v.MarshalJSON() = _, nil; want err != nil", test.in)
			continue
		}

		// XX(heckman): we make sure this error, if it exists, is expected above
		// so we can safely move on to the next iteration here...
		if err != nil {
			continue
		}

		if tm.Time().UnixNano() != test.out.Time().UnixNano() {
			t.Errorf("\nDescription: %s\nTime.UnmarshalJSON() = %#v; want %#v", test.desc, tm.Time().String(), test.out.Time().String())
		}
	}
}
