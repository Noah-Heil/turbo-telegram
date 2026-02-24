// Package generator provides diagram generation functionality.
package generator

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)

// CompressXML compresses XML data using zlib.
func CompressXML(xmlData []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)

	_, err := writer.Write(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to write zlib data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %w", err)
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(encoded, buf.Bytes())

	return encoded, nil
}

// CompressXMLWriter compresses XML data and writes to an io.Writer.
func CompressXMLWriter(xmlData []byte, w io.Writer) error {
	writer := zlib.NewWriter(w)
	defer func() { _ = writer.Close() }()

	_, err := writer.Write(xmlData)
	if err != nil {
		return fmt.Errorf("failed to write zlib data: %w", err)
	}
	return nil
}

// CompressAndEncode compresses XML data and returns base64 encoded string.
func CompressAndEncode(xmlData []byte) (string, error) {
	compressed, err := CompressXML(xmlData)
	if err != nil {
		return "", fmt.Errorf("failed to compress xml: %w", err)
	}
	return string(compressed), nil
}
