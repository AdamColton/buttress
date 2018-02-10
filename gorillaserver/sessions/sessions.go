package sessions

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Store struct {
	sessions.Store
}

func (s *Store) Get(r *http.Request, name string) (*Session, error) {
	sess, err := s.Store.Get(r, name)
	if err != nil || sess == nil {
		return nil, err
	}
	return &Session{sess}, err
}

func (s *Store) New(r *http.Request, name string) (*Session, error) {
	sess, err := s.Store.New(r, name)
	if err != nil || sess == nil {
		return nil, err
	}
	return &Session{sess}, err
}

// Session extends Gorilla sessions with some helper methods.
type Session struct {
	*sessions.Session
}

func NewCookieStore(keyPairs ...[]byte) *Store {
	return &Store{sessions.NewCookieStore(keyPairs...)}
}

func NewFilesystemStore(path string, keyPairs ...[]byte) *Store {
	return &Store{sessions.NewFilesystemStore(path, keyPairs...)}
}
