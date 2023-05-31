package middleware

import (
	"net/http"
)

type AdminMiddleware struct {
	xid string
}

func NewAdminMiddleware(xid string) *AdminMiddleware {
	return &AdminMiddleware{
		xid: xid,
	}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xid := r.Header.Get("xid")
		if xid != m.xid {
			return
		}
		next(w, r)
	}
}
