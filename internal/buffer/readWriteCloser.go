package buffer

import "bytes"

// ReadWriteCloser ReadWriteCloser
type ReadWriteCloser struct {
	WBuf *bytes.Buffer
	RBuf *bytes.Buffer
}

// Read read
func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.RBuf.Read(p)
}

// Write write
func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.WBuf.Write(p)
}

// Close close
func (rwc *ReadWriteCloser) Close() error {
	rwc.RBuf.Reset()
	rwc.WBuf.Reset()
	return nil
}
