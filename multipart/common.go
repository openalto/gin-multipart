package multipart

import (
	"net/textproto"
)

// Part defines a `body part`, which refers to an entity inside of a multipart
// entity.
// See https://tools.ietf.org/html/rfc2045#section-2.5
type Part struct {
	Header textproto.MIMEHeader
	Body   []byte
}
