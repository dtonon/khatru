package khatru

import (
	"hash/maphash"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/nbd-wtf/go-nostr"
)

const (
	AUTH_CONTEXT_KEY = iota
	WS_KEY
)

var nip20prefixmatcher = regexp.MustCompile(`^\w+: `)

func pointerHasher[V any](_ maphash.Seed, k *V) uint64 {
	return uint64(uintptr(unsafe.Pointer(k)))
}

func isOlder(previous, next *nostr.Event) bool {
	return previous.CreatedAt < next.CreatedAt ||
		(previous.CreatedAt == next.CreatedAt && previous.ID > next.ID)
}

func isAuthRequired(msg string) bool {
	idx := strings.IndexByte(msg, ':')
	return msg[0:idx] == "auth-required"
}

func getServiceBaseURL(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		if host == "localhost" {
			proto = "http"
		} else if strings.Index(host, ":") != -1 {
			// has a port number
			proto = "http"
		} else if _, err := strconv.Atoi(strings.ReplaceAll(host, ".", "")); err == nil {
			// it's a naked IP
			proto = "http"
		} else {
			proto = "https"
		}
	}
	return proto + "://" + host
}