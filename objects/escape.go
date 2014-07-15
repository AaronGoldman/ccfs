package objects

import (
	"encoding/hex"
	"net/url"
	"strconv"
)

func ObjectTypeEscape(objectType string) string {
	return url.QueryEscape(objectType)
}
func ObjectTypeUnescape(s string) (string, error) {
	objectType, err := url.QueryUnescape(s)
	return objectType, err
}

func NameSegmentEscape(nameSegment string) string {
	return url.QueryEscape(nameSegment)
}
func NameSegmentUnescape(s string) (string, error) {
	nameSegment, err := url.QueryUnescape(s)
	return nameSegment, err
}

func VersionEscape(version int64) (versionstr string) {
	return strconv.FormatInt(version, 10)
}
func VersionUnescape(s string) (int64, error) {
	version, err := strconv.ParseInt(s, 10, 64)
	return version, err
}

func SignatureEscape(signature []byte) string {
	return hex.EncodeToString(signature)
}
func SignatureUnescape(s string) ([]byte, error) {
	signature, err := hex.DecodeString(s)
	return signature, err
}
