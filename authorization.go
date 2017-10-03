package portal

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"sync"

	"github.com/pkg/errors"
)

type authStore struct {
}

//go:generate stringer -type=AuthorizationState

// AuthorizationState is the state of an Authorization
type AuthorizationState int

// s AuthorizationState json.Marshaler

// var foo json.Marshaler
// MarshalJSON() ([]byte, error)

const (
	AuthorizationPending AuthorizationState = iota
	AuthorizationAccepted
	// AuthorizationDenied
	// AuthorizationTimeout
)

type Authorization struct {
	ID      string
	State   AuthorizationState
	Request *jsonRPCRequest
}

type authorizationStore struct {
	authorizaitons map[string]*Authorization

	mu sync.Mutex
}

func newAuthorizationStore() *authorizationStore {
	return &authorizationStore{
		authorizaitons: make(map[string]*Authorization),
	}
}

func (s *authorizationStore) create(req *jsonRPCRequest) (*Authorization, error) {
	var buf [32]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return nil, err
	}

	id := base64.RawURLEncoding.EncodeToString(buf[:])

	auth := &Authorization{
		ID:      id,
		State:   AuthorizationPending,
		Request: req,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO automatic removal after 5 minutes
	s.authorizaitons[id] = auth

	return auth, nil
}

// consume confirms that an RPC request had been accepted by user.
// It RPC request had been accepted, it returns true exactly once.
func (s *authorizationStore) verify(id string, req *jsonRPCRequest) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	auth, found := s.authorizaitons[id]

	if !found {
		return false
	}

	if auth.State != AuthorizationAccepted {
		return false
	}

	same := auth.Request.Method == req.Method && bytes.Equal(auth.Request.Params, req.Params)

	if !same {
		return false
	}

	delete(s.authorizaitons, id)

	return true
}

func (s *authorizationStore) exists(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.authorizaitons[id]
	return found
}

func (s *authorizationStore) accept(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	auth, found := s.authorizaitons[id]

	if !found {
		return errors.New("Authorization not found")
	}

	auth.State = AuthorizationAccepted

	return nil
}

func (s *authorizationStore) deny(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.authorizaitons[id]

	if !found {
		return errors.New("Authorization not found")
	}

	delete(s.authorizaitons, id)

	return nil
}

// denyAuthorization(jsonRequest)
// acceptAuthorization(jsonRequest)
