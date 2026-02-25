package generator

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"testing"

	"diagram-gen/internal/testutil"
)

type errWriter struct{}

func (e errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

type closeFailWriter struct {
	writes int
}

func lockCompression(t *testing.T) {
	t.Helper()
	testutil.LockGlobal()
	t.Cleanup(testutil.UnlockGlobal)
}

func (c *closeFailWriter) Write(p []byte) (int, error) {
	c.writes++
	if c.writes > 1 {
		return 0, errors.New("close write failed")
	}
	return len(p), nil
}

func TestCompressXMLWriterWriteError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	err := compressXMLToWriter([]byte("test"), errWriter{})
	if err == nil {
		t.Fatal("expected error from writer")
	}
}

func TestCompressXMLWriterCloseError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	writer := &closeFailWriter{}
	err := compressXMLToWriter([]byte("test"), writer)
	if err == nil {
		t.Fatal("expected error from close")
	}
}

func TestCompressXMLWithLevelCreateError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	prev := newZlibWriterLevel
	newZlibWriterLevel = func(_ io.Writer, _ int) (*zlib.Writer, error) {
		return nil, fmt.Errorf("create failed")
	}
	t.Cleanup(func() { newZlibWriterLevel = prev })

	_, err := CompressXMLWithLevel([]byte("test"), 1)
	if err == nil {
		t.Fatal("expected error from writer creation")
	}
}

func TestCompressXMLCreateError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	prev := newZlibWriterLevel
	newZlibWriterLevel = func(_ io.Writer, _ int) (*zlib.Writer, error) {
		return nil, fmt.Errorf("create failed")
	}
	t.Cleanup(func() { newZlibWriterLevel = prev })

	_, err := CompressXML([]byte("test"))
	if err == nil {
		t.Fatal("expected error from writer creation")
	}
}

func TestCompressAndEncodeWithLevelCreateError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	prev := newZlibWriterLevel
	newZlibWriterLevel = func(_ io.Writer, _ int) (*zlib.Writer, error) {
		return nil, fmt.Errorf("create failed")
	}
	t.Cleanup(func() { newZlibWriterLevel = prev })

	_, err := CompressAndEncodeWithLevel([]byte("test"), 1)
	if err == nil {
		t.Fatal("expected error from writer creation")
	}
}

func TestCompressAndEncodeCreateError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	prev := newZlibWriterLevel
	newZlibWriterLevel = func(_ io.Writer, _ int) (*zlib.Writer, error) {
		return nil, fmt.Errorf("create failed")
	}
	t.Cleanup(func() { newZlibWriterLevel = prev })

	_, err := CompressAndEncode([]byte("test"))
	if err == nil {
		t.Fatal("expected error from writer creation")
	}
}

func TestCompressXMLWithLevelToWriterWriteError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	err := compressXMLWithLevelToWriter([]byte("test"), errWriter{}, 1)
	if err == nil {
		t.Fatal("expected error from writer")
	}
}

func TestCompressXMLWithLevelToWriterCloseError(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	writer := &closeFailWriter{}
	err := compressXMLWithLevelToWriter([]byte("test"), writer, 1)
	if err == nil {
		t.Fatal("expected error from close")
	}
}

func TestCompressXMLWriterSuccess(t *testing.T) {
	t.Parallel()
	lockCompression(t)
	var buf bytes.Buffer
	err := compressXMLToWriter([]byte("test"), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected buffer to be written")
	}
}
