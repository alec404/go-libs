package net

import (
	"errors"
	"strings"
	"syscall"
)

// Returns if the given err is "connection reset by peer" error.
func IsConnectionReset(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.ECONNRESET)
	}
	return false
}

// Returns if the given err is "http2: client connection lost" error.
func IsHTTP2ConnectionLost(err error) bool {
	return err != nil && strings.Contains(err.Error(), "http2: client connection lost")
}

// Returns if the given err is  "http2: client connection force closed via ClientConn.Close" error.
func IsHTTP2ConnectionForceClosed(err error) bool {
	return err != nil && strings.Contains(err.Error(), "http2: client connection force closed via ClientConn.Close")
}

// Returns if the given err is "connection refused" error
func IsConnectionRefused(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.ECONNREFUSED)
	}
	return false
}

func IsBrokenPipe(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.EPIPE)
	}
	return false
}
