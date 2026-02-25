// Package generator provides diagram generation functionality.
package generator

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)

type CompressionLevel int

const (
	BestCompression      CompressionLevel = CompressionLevel(zlib.BestCompression)
	BestSpeedCompression CompressionLevel = CompressionLevel(zlib.BestSpeed)
	DefaultCompression   CompressionLevel = CompressionLevel(zlib.DefaultCompression)
	NoCompression        CompressionLevel = CompressionLevel(zlib.NoCompression)
	HuffmanOnly          CompressionLevel = CompressionLevel(zlib.HuffmanOnly)
)

var newZlibWriterLevel = zlib.NewWriterLevel

// CompressXML compresses XML data using zlib.
func CompressXML(xmlData []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := compressXMLToWriter(xmlData, &buf)
	if err != nil {
		return nil, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(encoded, buf.Bytes())
	return encoded, nil
}

func compressXMLToWriter(xmlData []byte, w io.Writer) error {
	writer, err := newZlibWriterLevel(w, zlib.DefaultCompression)
	if err != nil {
		return fmt.Errorf("failed to create zlib writer: %w", err)
	}

	_, err = writer.Write(xmlData)
	if err != nil {
		return fmt.Errorf("failed to write zlib data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return nil
}

// CompressXMLWithLevel compresses XML data using specified zlib compression level.
func CompressXMLWithLevel(xmlData []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	err := compressXMLWithLevelToWriter(xmlData, &buf, level)
	if err != nil {
		return nil, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(encoded, buf.Bytes())
	return encoded, nil
}

func compressXMLWithLevelToWriter(xmlData []byte, w io.Writer, level int) error {
	writer, err := newZlibWriterLevel(w, level)
	if err != nil {
		return fmt.Errorf("failed to create zlib writer: %w", err)
	}

	_, err = writer.Write(xmlData)
	if err != nil {
		return fmt.Errorf("failed to write zlib data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return nil
}

// CompressXMLWriter compresses XML data and writes to an io.Writer.
func CompressXMLWriter(xmlData []byte, w io.Writer) error {
	return compressXMLToWriter(xmlData, w)
}

// CompressAndEncode compresses XML data and returns base64 encoded string.
func CompressAndEncode(xmlData []byte) (string, error) {
	compressed, err := CompressXML(xmlData)
	if err != nil {
		return "", fmt.Errorf("failed to compress xml: %w", err)
	}
	return string(compressed), nil
}

// CompressAndEncodeWithLevel compresses XML data with specified level and returns base64 encoded string.
func CompressAndEncodeWithLevel(xmlData []byte, level int) (string, error) {
	compressed, err := CompressXMLWithLevel(xmlData, level)
	if err != nil {
		return "", fmt.Errorf("failed to compress xml: %w", err)
	}
	return string(compressed), nil
}
