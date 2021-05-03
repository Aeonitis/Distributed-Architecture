package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	// encoding is the encoding that we use to persist record sizes & index entries
	encoding = binary.BigEndian
)

const (
	// lenWidth defines the number of bytes used to store the record’s length.
	lenWidth = 8
)

// store a simple wrapper around a file with two APIs to append & read bytes to and from the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// newStore creates a store for the given file
func newStore(f *os.File) (*store, error) {

	// os.Stat to get file’s current size, in case of re-creating the store from
	//a file that has existing data, if e.g. our service had restarted.
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append persists the given bytes to the store.
func (s *store) Append(p []byte) (numBytes uint64, position uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	position = s.size

	// Write length of the record so when I read the record, I know how many bytes to read.
	if err := binary.Write(s.buf, encoding, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	// Write to buffered writer instead of directly to file to reduce number of system calls and improve performance.
	numberOfWrittenBytes, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	numberOfWrittenBytes += lenWidth
	s.size += uint64(numberOfWrittenBytes)
	return uint64(numberOfWrittenBytes), position, nil
}

// Read returns the record stored at the given position
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First flush the writer buffer, in case I try to read a record that the buffer has not flushed to disk yet.
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// Find out how many bytes we have to read to get the whole record
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	// Fetch and return the record.
	// Compiler allocates byte slices that don’t escape the functions they’re declared in on the stack.
	// A value escapes when it lives beyond the lifetime of the function call e.g. if we return the value.
	record := make([]byte, encoding.Uint64(size))
	if _, err := s.File.ReadAt(record, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return record, nil
}

// ReadAt reads len(b) bytes into b, beginning at the off offset in the store’s file
func (s *store) ReadAt(b []byte, offset int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	// Adapter of io.ReaderAt function
	return s.File.ReadAt(b, offset)
}

// Close persists any buffered data before closing the file
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
