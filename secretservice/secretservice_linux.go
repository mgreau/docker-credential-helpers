package secretservice

/*
#cgo pkg-config: libsecret-1

#include "secretservice_linux.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/docker/docker-credential-helpers/credentials"
)

type secretservice struct{}

// New creates a new secretservice.
func New() credentials.Helper {
	return secretservice{}
}

// Add adds new credentials to the keychain.
func (h secretservice) Add(creds *credentials.Credentials) error {
	if creds == nil {
		return errors.New("missing credentials")
	}
	server := C.CString(creds.ServerURL)
	defer C.free(unsafe.Pointer(server))
	username := C.CString(creds.Username)
	defer C.free(unsafe.Pointer(username))
	secret := C.CString(creds.Secret)
	defer C.free(unsafe.Pointer(secret))

	if err := C.add(server, username, secret); err != nil {
		defer C.g_error_free(err)
		errMsg := (*C.char)(unsafe.Pointer(err.message))
		return errors.New(C.GoString(errMsg))
	}
	return nil
}

// Delete removes credentials from the keychain.
func (h secretservice) Delete(serverURL string) error {
	if serverURL == "" {
		return errors.New("missing server url")
	}
	server := C.CString(serverURL)
	defer C.free(unsafe.Pointer(server))

	if err := C.delete(server); err != nil {
		defer C.g_error_free(err)
		errMsg := (*C.char)(unsafe.Pointer(err.message))
		return errors.New(C.GoString(errMsg))
	}
	return nil
}

// Get returns the username and secret to use for a given registry server URL.
func (h secretservice) Get(serverURL string) (string, string, error) {
	if serverURL == "" {
		return "", "", errors.New("missing server url")
	}
	var username *C.char
	defer C.free(unsafe.Pointer(username))
	var secret *C.char
	defer C.free(unsafe.Pointer(secret))
	server := C.CString(serverURL)
	defer C.free(unsafe.Pointer(server))

	err := C.get(server, &username, &secret)
	if err != nil {
		defer C.g_error_free(err)
		errMsg := (*C.char)(unsafe.Pointer(err.message))
		return "", "", errors.New(C.GoString(errMsg))
	}
	user := C.GoString(username)
	pass := C.GoString(secret)
	if pass == "" {
		return "", "", credentials.ErrCredentialsNotFound
	}
	return user, pass, nil
}
