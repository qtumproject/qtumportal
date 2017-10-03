package portal

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationStore(t *testing.T) {
	is := assert.New(t)
	s := newAuthorizationStore()

	req := &jsonRPCRequest{
		Method: "foo",
		Params: []byte("[1,2,3]"),
	}

	badreq1 := &jsonRPCRequest{
		Method: "foo2",
		Params: []byte("[1,2,3]"),
	}

	badreq2 := &jsonRPCRequest{
		Method: "foo",
		Params: []byte("[1,2]"),
	}

	req2 := &jsonRPCRequest{
		Method: "foo",
		Params: []byte("[1,2,3]"),
	}

	auth, err := s.create(req)
	is.NoError(err)

	encodedLength := base64.RawURLEncoding.EncodedLen(32)
	is.Equal(encodedLength, len(auth.ID), "Should generate 32 bytes ID")

	id := auth.ID

	is.False(s.verify(id, badreq2), "Should fail if request not accepted")
	err = s.accept(id)
	is.NoError(err)

	is.False(s.verify("no such id", req2), "Should fail if id is not found")
	is.False(s.verify(id, badreq1), "Should fail if Method is different")
	is.False(s.verify(id, badreq2), "Should fail if Params is different")

	is.True(s.verify(id, req2), "Should succeed if id and request are the same")
	is.False(s.verify(id, req2), "Should succeed only once for an authorization")

	// spew.Dump(auth)
}

func TestAuthorizationStoreDeny(t *testing.T) {
	is := assert.New(t)
	s := newAuthorizationStore()

	req := &jsonRPCRequest{
		Method: "foo",
		Params: []byte("[1,2,3]"),
	}

	auth, err := s.create(req)
	is.NoError(err)
	is.True(s.exists(auth.ID), "Authorization should exist for ID")

	err = s.deny(auth.ID)
	is.NoError(err)
	is.False(s.exists(auth.ID), "Authorization should not exist for ID after denial")
}
