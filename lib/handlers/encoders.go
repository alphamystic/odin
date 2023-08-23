package handlers

import (
	"fmt"
	"bytes"
	"strings"
	"net/url"
	"encoding/base64"
)
// EncodingMethod represents an encoding method that can be used to bypass a WAF.
type EncodingMethod int

type Encoder interface{
	Encode() string //returns a string of encoded payload
}
const (
	// URLEncoding represents URL encoding.
	URLEncoding EncodingMethod = iota
	// HTMLEncoding represents HTML encoding.
	HTMLEncoding
	// Base64Encoding represents base64 encoding.
	Base64Encoding
	// DoubleEncoding represents double encoding.
	DoubleEncoding
	// LeadingZeros represents adding leading zeros.
	LeadingZeros
	// UnicodeEscaping represents Unicode escaping.
	UnicodeEscaping
	// OctalEscaping represents octal escaping.
	OctalEscaping
	// MultipleEncoding represents multiple encoding.
	MultipleEncoding
	// HexEncoding represents hex encoding.
	HexEncoding
)

// encodePayload encodes the given payload using the specified encoding method.
func encodePayload(payload string, method EncodingMethod) string {
	switch method {
	case URLEncoding:
		return url.QueryEscape(payload)
	case HTMLEncoding:
		var b bytes.Buffer
		//html.Escape(&b, []byte(payload))
		return b.String()
	case Base64Encoding:
		return base64.StdEncoding.EncodeToString([]byte(payload))
	case DoubleEncoding:
		encodedPayload := url.QueryEscape(payload)
		return url.QueryEscape(encodedPayload)
	case LeadingZeros:
		return strings.Replace(payload, ".", "%00.", -1)
	case UnicodeEscaping:
		var b bytes.Buffer
		for _, c := range payload {
			fmt.Fprintf(&b, "\\u%04x", c)
		}
		return b.String()
	case OctalEscaping:
		var b bytes.Buffer
		for _, c := range payload {
			fmt.Fprintf(&b, "\\%03o", c)
		}
		return b.String()
	case MultipleEncoding:
		encodedPayload := url.QueryEscape(payload)
		encodedPayload = strings.Replace(encodedPayload, ".", "%00.", -1)
		encodedPayload = base64.StdEncoding.EncodeToString([]byte(encodedPayload))
		return encodedPayload
	case HexEncoding:
		var b bytes.Buffer
		for _, c := range payload {
			fmt.Fprintf(&b, "\\x%02x", c)
		}
		return b.String()
	default:
		return payload
	}
}
