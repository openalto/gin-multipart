// Copyright 2019 Jensen Zhang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package multipart

import (
	"bytes"
	"mime"
	mimeMultipart "mime/multipart"
	"net/http"
	"strings"
)

const relatedContentType = "multipart/related"

// Related defines a message of MIME type `multipart/related`.
// See https://tools.ietf.org/html/rfc2387
type Related struct {
	// Required Parameters
	Type     string
	Boundary string

	// Optional Parameters
	Start     string
	StartInfo string

	// Multipart
	RootPart      Part
	ChildrenParts []Part

	Data interface{}
}

// Render implements the callback method of Related type for gin data render
func (r Related) Render(w http.ResponseWriter) error {
	var b bytes.Buffer
	mw := mimeMultipart.NewWriter(&b)
	if r.Boundary != "" {
		if err := mw.SetBoundary(r.Boundary); err != nil {
			panic(err)
		}
	} else {
		r.Boundary = mw.Boundary()
	}
	rw, err := mw.CreatePart(r.RootPart.Header)
	if err != nil {
		panic(err)
	}
	rw.Write(r.RootPart.Body)
	for _, cp := range r.ChildrenParts {
		cw, err := mw.CreatePart(cp.Header)
		if err != nil {
			panic(err)
		}
		cw.Write(cp.Body)
	}
	r.WriteContentType(w)
	w.Write(b.Bytes())
	if err := mw.Close(); err != nil {
		panic(err)
	}
	return nil
}

// WriteContentType implements the callback method of Related type for gin data render
func (r Related) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	contentType := r.ContentType()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{contentType}
	}
}

// ContentType returns the format media-type of multipart entity r
func (r Related) ContentType() string {
	if r.Boundary == "" {
		panic("Require a boundary")
	}
	params := map[string]string{
		"boundary": r.Boundary,
	}
	if r.Type != "" {
		params["type"] = escapeQuotes(r.Type)
	}
	if r.Start != "" {
		params["start"] = r.Start
	}
	if r.StartInfo != "" {
		params["start-info"] = escapeQuotes(r.StartInfo)
	}
	return mime.FormatMediaType(relatedContentType, params)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
