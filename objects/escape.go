//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package objects

import (
	"encoding/hex"
	"net/url"
	"strconv"
)

//ObjectTypeEscape escapes the text of types
func ObjectTypeEscape(objectType string) string {
	return url.QueryEscape(objectType)
}

//ObjectTypeUnescape unescapes the text of types
func ObjectTypeUnescape(s string) (string, error) {
	objectType, err := url.QueryUnescape(s)
	return objectType, err
}

//NameSegmentEscape escapes the text of NameSegment
func NameSegmentEscape(nameSegment string) string {
	return url.QueryEscape(nameSegment)
}

//NameSegmentUnescape escapes the text of NameSegment
func NameSegmentUnescape(s string) (string, error) {
	nameSegment, err := url.QueryUnescape(s)
	return nameSegment, err
}

//VersionEscape escapes the text of version
func VersionEscape(version int64) (versionstr string) {
	return strconv.FormatInt(version, 10)
}

//VersionUnescape escapes the text of version
func VersionUnescape(s string) (int64, error) {
	version, err := strconv.ParseInt(s, 10, 64)
	return version, err
}

//SignatureEscape escapes the text of signature
func SignatureEscape(signature []byte) string {
	return hex.EncodeToString(signature)
}

//SignatureUnescape escapes the text of signature
func SignatureUnescape(s string) ([]byte, error) {
	signature, err := hex.DecodeString(s)
	return signature, err
}
