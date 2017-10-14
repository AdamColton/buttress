// Package gorillauser provides a set of tools for using gorilla sessions to
// store a user. It uses bcrypt for passwords.
package gorillauser

import (
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var StoreName = "gorillauser"

type User interface {
	GetHashedPassword() []byte
	SetHashedPassword([]byte)
	GetID() []byte
	GetName() string
}

type UserStore interface {
	GetByName(name string) (User, error)
	GetByID(id []byte) (User, error)
	Create(name string, hashedPassword []byte) (User, error)
}

type GorillaUserStore struct {
	UserStore
	sessions.Store
	Field struct {
		Name            string
		Password        string
		ConfirmPassword string
		MinNameLen      int
		MinPasswordLen  int
	}
}

func (gus *GorillaUserStore) Get(r *http.Request) (*sessions.Session, error) {
	s, err := gus.Store.Get(r, StoreName)
	if err == nil && s == nil {
		s, err = gus.Store.New(r, StoreName)
	}
	return s, err
}

func (gus *GorillaUserStore) Login(name, password string, r *http.Request, w http.ResponseWriter) (User, error) {
	u, err := gus.GetByName(name)
	if u == nil || err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(u.GetHashedPassword(), []byte(password))
	if err != nil {
		return nil, err
	}

	s, err := gus.Get(r)
	if err != nil {
		return nil, err
	}

	return u, save(u, s, r, w)
}

func save(u User, s *sessions.Session, r *http.Request, w http.ResponseWriter) error {
	s.Values["name"] = u.GetName()
	s.Values["id"] = u.GetID()
	err := s.Save(r, w)
	return err
}

func (gus *GorillaUserStore) Create(name, password string, r *http.Request, w http.ResponseWriter) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u, err := gus.UserStore.Create(name, hashedPassword)
	if err != nil {
		return u, err
	}

	s, err := gus.Get(r)
	if err != nil {
		return u, err
	}

	return u, save(u, s, r, w)
}

func (gus *GorillaUserStore) Auth(r *http.Request) ([]byte, error) {
	s, err := gus.Store.Get(r, StoreName)
	if err != nil || s == nil {
		return nil, err
	}

	id, ok := s.Values["id"].([]byte)
	if !ok {
		return nil, nil
	}

	return id, nil
}

func (gus *GorillaUserStore) AuthRoute(loggedInHandler, notLoggedInHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id, _ := gus.Auth(r); id == nil {
			notLoggedInHandler(w, r)
		} else {
			loggedInHandler(w, r)
		}
	}
}

func (gus *GorillaUserStore) SignOut(r *http.Request, w http.ResponseWriter) error {
	s, err := gus.Store.Get(r, StoreName)
	if err != nil || s == nil {
		return err
	}
	delete(s.Values, "name")
	delete(s.Values, "id")
	return s.Save(r, w)
}

func (gus *GorillaUserStore) CreateHandler(success, failure http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := gus.Get(r)
		if err != nil {
			failure(w, r)
			return
		}

		r.ParseForm()
		var errStr string
		name := r.FormValue(gus.Field.Name)
		if len([]rune(name)) < gus.Field.MinNameLen {
			errStr = "Name too short"
		}
		password := r.FormValue(gus.Field.Password)
		if errStr == "" && len([]rune(password)) < gus.Field.MinPasswordLen {
			errStr = "Password too short"
		}
		confirm := r.FormValue(gus.Field.ConfirmPassword)
		if errStr == "" && password != confirm {
			errStr = "Passwords don't match"
		}
		if errStr != "" {
			s.Values["error"] = errStr
			s.Save(r, w)
			failure(w, r)
			return
		}

		u, err := gus.Create(name, password, r, w)
		if err != nil || u == nil {
			if err != nil {
				s.Values["error"] = err.Error()
				s.Save(r, w)
			}
			failure(w, r)
		} else {
			delete(s.Values, "error")
			success(w, r)
		}
	}
}

func (gus *GorillaUserStore) LoginHandler(success, failure http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue(gus.Field.Name)
		password := r.FormValue(gus.Field.Password)

		u, err := gus.Login(name, password, r, w)
		if err != nil || u == nil {
			failure(w, r)
		} else {
			success(w, r)
		}
	}
}

func (gus *GorillaUserStore) SignOutHandler(after http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gus.SignOut(r, w)
		after(w, r)
	}
}
