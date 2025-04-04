package session_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/harrybrwn/x/session"
)

func Example() {
	type data struct {
		ID, Name string
	}
	session.RegisterSerializable(&data{})
	store := session.NewMemStore[data](time.Second)
	sessions := session.NewManager("cookies", store)

	http.Handle("/session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := sessions.Get(r)
		if err != nil {
			slog.Error("failed to get session", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, "hello %s, you're user number %s", s.Value.Name, s.Value.ID)
	}))

	http.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("password") != "passw0rd1" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s := sessions.NewSession(&data{ID: "1", Name: "Jimmy"})
		err := s.Save(r.Context())
		if err != nil {
			slog.Error("failed to save session", "error", err)
			w.WriteHeader(500)
			return
		}
		http.SetCookie(w, s.Cookie())
		w.WriteHeader(200)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/login?password=passw0rd1", nil)
	http.DefaultServeMux.ServeHTTP(rec, req)
	cookie := rec.Result().Header.Get("Set-Cookie")

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/session", nil)
	req.Header.Set("Cookie", cookie)
	http.DefaultServeMux.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// Output:
	// hello Jimmy, you're user number 1
}

// func ExampleMiddleware() {
// 	storage := session.NewMemStore[int](time.Minute)
// 	sessions := session.NewManager("user", storage)
// 	logging := func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			s, err := sessions.Get(r)
// 			if err == nil && s.Value != nil {
// 				slog.Info("request from registered user", "id", *s.Value)
// 			}
// 			h.ServeHTTP(w, r)
// 		})
// 	}
// 	http.Handle("/stuff", logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// stuff...
// 	})))
// }
