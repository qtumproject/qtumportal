package portal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/olebedev/emitter"

	"github.com/pkg/errors"
)

type authStore struct {
}

// AuthorizationState is the state of an Authorization
type AuthorizationState string

// s AuthorizationState json.Marshaler

// var foo json.Marshaler
// MarshalJSON() ([]byte, error)

const (
	AuthorizationPending AuthorizationState = "pending"
	// transition => {accepted, denied}

	AuthorizationAccepted AuthorizationState = "accepted"
	// transition => { consumed }

	AuthorizationDenied AuthorizationState = "denied"
	// transition => {}

	AuthorizationConsumed AuthorizationState = "consumed"
	// transition => {}

	// AuthorizationTimeout
)

type Authorization struct {
	ID        string             `json:"id"`
	State     AuthorizationState `json:"state"`
	Request   *jsonRPCRequest    `json:"request"`
	CreatedAt time.Time          `json:"createdAt"`
}

type authorizationStore struct {
	authorizaitons map[string]*Authorization

	events *emitter.Emitter

	mu sync.Mutex
}

func authChangeTopic(id string) string {
	return fmt.Sprintf("change:%s", id)
}

func newAuthorizationStore() *authorizationStore {
	return &authorizationStore{
		events:         &emitter.Emitter{},
		authorizaitons: make(map[string]*Authorization),
	}
}

// waitChange blocks until an authorization changes its state
func (s *authorizationStore) waitChange(ctx context.Context, id string) error {
	s.mu.Lock()

	auth, found := s.authorizaitons[id]
	if !found {
		s.mu.Unlock()
		return errors.New("not found")
	}

	if auth.State != AuthorizationPending {
		s.mu.Unlock()
		return errors.New("not pending")
	}

	// timeout
	topic := authChangeTopic(id)
	resolved := s.events.Once(topic)
	s.mu.Unlock()

	select {
	case <-resolved:
		return nil
	case <-ctx.Done():
		s.events.Off(topic, resolved)
	}

	return nil
}

func (s *authorizationStore) get(id string) (*Authorization, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	auth, found := s.authorizaitons[id]
	return auth, found
}

func (s *authorizationStore) allAuthorizations() []*Authorization {
	s.mu.Lock()
	defer s.mu.Unlock()

	var auths []*Authorization

	for _, auth := range s.authorizaitons {
		auths = append(auths, auth)
	}

	sort.Slice(auths, func(i, j int) bool {
		return auths[i].CreatedAt.After(auths[j].CreatedAt)
	})

	return auths
}

func (s *authorizationStore) pendingAuthorizations() []*Authorization {
	var auths []*Authorization

	for _, auth := range s.allAuthorizations() {
		if auth.State == AuthorizationPending {
			auths = append(auths, auth)
		}
	}

	return auths
}

func (s *authorizationStore) create(req *jsonRPCRequest) (*Authorization, error) {
	var buf [32]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return nil, err
	}

	id := base64.RawURLEncoding.EncodeToString(buf[:])

	auth := &Authorization{
		ID:        id,
		State:     AuthorizationPending,
		Request:   req,
		CreatedAt: time.Now(),
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

	auth.State = AuthorizationConsumed

	s.emitAuthChange(id)

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

	if auth.State != AuthorizationPending {
		return errors.New("Authorization not pending")
	}

	auth.State = AuthorizationAccepted

	s.emitAuthChange(id)

	return nil
}

func (s *authorizationStore) deny(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	auth, found := s.authorizaitons[id]

	if !found {
		return errors.New("Authorization not found")
	}

	if auth.State != AuthorizationPending {
		return errors.New("Authorization not pending")
	}

	auth.State = AuthorizationDenied

	s.emitAuthChange(id)

	return nil
}

func (s *authorizationStore) emitAuthChange(id string) {
	s.events.Emit(authChangeTopic(id))
}
