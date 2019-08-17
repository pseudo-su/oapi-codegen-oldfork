// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codegen

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

// GenerateInlinedSpec generates a gzipped, base64 encoded JSON representation of the
// swagger definition, which we embed inside the generated code.
func GenerateInlinedSpec(t *template.Template, swagger *openapi3.Swagger) (string, error) {
	// Marshal to json
	encoded, err := swagger.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("error marshaling swagger: %s", err)
	}

	// gzip
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", fmt.Errorf("error creating gzip compressor: %s", err)
	}
	_, err = zw.Write(encoded)
	if err != nil {
		return "", fmt.Errorf("error gzipping swagger file: %s", err)
	}
	err = zw.Close()
	if err != nil {
		return "", fmt.Errorf("error gzipping swagger file: %s", err)
	}
	str := base64.StdEncoding.EncodeToString(buf.Bytes())

	var parts []string
	const width = 80

	// Chop up the string into an array of strings.
	for len(str) > width {
		part := str[0:width]
		parts = append(parts, part)
		str = str[width:]
	}
	if len(str) > 0 {
		parts = append(parts, str)
	}

	// Generate inline code.
	buf.Reset()
	w := bufio.NewWriter(&buf)
	err = t.ExecuteTemplate(w, "inline.tmpl", parts)
	if err != nil {
		return "", fmt.Errorf("error generating inlined spec: %s", err)
	}
	err = w.Flush()
	if err != nil {
		return "", fmt.Errorf("error flushing output buffer for inlined spec: %s", err)
	}
	return buf.String(), nil
}

// GenerateInlinedSpecUI generates a gzipped, base64 encoded HTML template of
// swagger ui, which we embed inside the generated code.
func GenerateInlinedSpecUI(t *template.Template, swagger *openapi3.Swagger) (string, error) {
	// load the HTML swagger ui template
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := t.ExecuteTemplate(w, "swagger.html.tmpl", nil)
	if err != nil {
		return "", fmt.Errorf("error loading swagger ui template: %s", err)
	}
	swaggerUI := buf.Bytes()
	buf.Reset()

	// gzip
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", fmt.Errorf("error creating gzip compressor: %s", err)
	}
	_, err = zw.Write(swaggerUI)
	if err != nil {
		return "", fmt.Errorf("error gzipping swagger file: %s", err)
	}
	err = zw.Close()
	if err != nil {
		return "", fmt.Errorf("error gzipping swagger file: %s", err)
	}
	str := base64.StdEncoding.EncodeToString(buf.Bytes())

	var parts []string
	const width = 80

	// Chop up the string into an array of strings.
	for len(str) > width {
		part := str[0:width]
		parts = append(parts, part)
		str = str[width:]
	}
	if len(str) > 0 {
		parts = append(parts, str)
	}

	// Generate inline code.
	buf.Reset()
	w = bufio.NewWriter(&buf)
	err = t.ExecuteTemplate(w, "inline-ui.tmpl", parts)
	if err != nil {
		return "", fmt.Errorf("error generating inlined spec ui: %s", err)
	}
	err = w.Flush()
	if err != nil {
		return "", fmt.Errorf("error flushing output buffer for inlined spec ui: %s", err)
	}
	return buf.String(), nil
}
