// Copyright 2019 Jensen Zhang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package multipart

import (
	"fmt"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderRelated(t *testing.T) {
	w := httptest.NewRecorder()
	rootType := "application/json"
	data := Related{
		Type: rootType,
		RootPart: Part{
			Header: textproto.MIMEHeader{
				"Content-Id":   []string{"root-resource"},
				"Content-Type": []string{rootType},
			},
			Body: []byte(`{"k": "v"}`),
		},
		ChildrenParts: []Part{
			Part{
				Header: textproto.MIMEHeader{
					"Content-Id":   []string{"related-resource-1"},
					"Content-Type": []string{"H/S"},
				},
				Body: []byte(`hello`),
			},
		},
	}
	// data.WriteContentType(w)
	// assert.Equal(t, "multipart/related", w.Header().Get("Content-Type"))

	err := data.Render(w)

	assert.NoError(t, err)
	fmt.Println(w.Header())
	fmt.Println(w.Body.String())
}
