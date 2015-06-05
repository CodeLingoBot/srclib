package graph

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// MakeURI converts a repository clone URL, such as
// "git://github.com/user/repo.git", to a normalized URI string, such
// as "github.com/user/repo" lexically. MakeURI panics if there is an
// error, and should only be used if cloneURL is a correctly-formed
// URL. It is a wrapper around TryMakeURI.
func MakeURI(cloneURL string) string {
	uri, err := TryMakeURI(cloneURL)
	if err != nil {
		panic(err)
	}
	return uri
}

// TryMakeURI converts a repository clone URL, such as
// "git://github.com/user/repo.git", to a normalized URI string, such
// as "github.com/user/repo" lexically. TryMakeURI returns an error if
// cloneURL is empty or malformed.
func TryMakeURI(cloneURL string) (string, error) {
	if cloneURL == "" {
		return "", errors.New("MakeURI: empty clone URL")
	}

	url, err := url.Parse(cloneURL)
	if err != nil {
		return "", err
	} else if url.Path == "" || url.Path == "/" {
		return "", fmt.Errorf("determining URI from repo clone URL failed: missing path from URL (%q)", cloneURL)
	} else if url.Host == "" && (url.Path[0] == '/' || !strings.Contains(strings.Trim(url.Path, "/"), "/")) {
		// We ensure our Path doesn't look like the output of TryMakeURI
		// so that the output of this function is a fixed point.
		// ie TryMakeURI("github.com/user/repo") == ("github.com/user/repo", nil),
		// not an error.
		return "", fmt.Errorf("determining URI from repo clone URL failed: missing host from URL (%q)", cloneURL)
	}

	path := strings.TrimSuffix(url.Path, ".git")
	path = filepath.Clean(path)
	path = strings.TrimSuffix(path, "/")
	return strings.ToLower(url.Host) + path, nil
}

// URIEqual returns true if a and b are equal, based on a case insensitive
// comparison (strings.EqualFold).
func URIEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
