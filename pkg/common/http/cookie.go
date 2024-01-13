package http

import (
	"net/http"
)

const (
	ErrFDuplicateCookie = "duplicate cookie: %s"
)

// SetCookie sets a cookie.
// Unfortunately, Go language is so ugly that it does not allow to check
// whether the specified cookie is already set or not.
func SetCookie(rw http.ResponseWriter, cookie *http.Cookie) (err error) {
	http.SetCookie(rw, cookie)
	return nil
}

//----------------------------------------------------------------------------//
// Q:	What do you like in Go language ?
//
// A:	Go language is a good example showing what a programming language
//		should NOT be. Moreover, by looking at how people use this language one
//		can learn their character and attitude to this world.
//----------------------------------------------------------------------------//
