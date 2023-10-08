package awsutl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/httputl"
)

const (
	AlgorithmKey = "X-Amz-Algorithm"
	AuthScheme   = "AWS4-HMAC-SHA256"

	ContentHashKey   = "X-Amz-Content-Sha256"
	EmptyContentHash = "UNSIGNED-PAYLOAD"

	QueryPrefix = "X-Amz-"

	CredentialEntryKey = "Credential"
	CredentialQueryKey = QueryPrefix + CredentialEntryKey

	SignedHeadersEntryKey = "SignedHeaders"
	SignedHeadersQueryKey = QueryPrefix + SignedHeadersEntryKey

	SignatureEntryKey = "Signature"
	SignatureQueryKey = QueryPrefix + SignatureEntryKey
)

var (
	ErrUnknownAuthScheme          = errors.New("unknown authorization scheme")
	ErrEmptyCredential            = errors.New("empty AWS credential")
	ErrEmptySignedHeaders         = errors.New("empty AWS signed headers")
	ErrEmptySignature             = errors.New("empty AWS signature")
	ErrEmptyContentHash           = errors.New("empty AWS content hash")
	ErrIncompleteCredentialFields = errors.New("incomplete AWS credential fields")
)

type Authorization struct {
	Scheme        string
	Credential    *Credential
	SignedHeaders []string
	Signature     string
	ContentHash   string
}

func NewAuthorization(r *http.Request) (*Authorization, error) {
	auth, err := newAuthorization(r)
	if err != nil {
		return nil, err
	}
	var errs []error
	if auth.Credential == nil {
		errs = append(errs, ErrEmptyCredential)
	}
	if auth.SignedHeaders == nil {
		errs = append(errs, ErrEmptySignedHeaders)
	}
	if auth.Signature == "" {
		errs = append(errs, ErrEmptySignature)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return auth, nil
}

func newAuthorization(r *http.Request) (*Authorization, error) {
	if query := r.URL.Query(); query.Get(AlgorithmKey) == AuthScheme {
		return newQueryAuthorization(query)
	}

	httpAuth, err := httputl.NewAuthorization(r.Header)
	if err != nil {
		return nil, err
	}
	auth, err := newHeaderAuthorization(httpAuth)
	if err != nil {
		return nil, err
	}

	contentHash := r.Header.Get(ContentHashKey)
	if contentHash == "" {
		return nil, ErrEmptyContentHash
	}
	auth.ContentHash = contentHash

	return auth, nil
}

func newHeaderAuthorization(a *httputl.Authorization) (*Authorization, error) {
	auth := &Authorization{Scheme: a.Scheme}
	for _, rawEntry := range strings.Split(a.Credential, ",") {
		entryItems := strings.SplitN(rawEntry, "=", 2)
		if len(entryItems) != 2 {
			continue
		}
		key := entryItems[0]
		value := entryItems[1]
		switch key {
		case CredentialEntryKey:
			cred, err := NewCredential(value)
			if err != nil {
				return nil, err
			}
			auth.Credential = cred
		case SignedHeadersEntryKey:
			auth.SignedHeaders = strings.Split(strings.ToLower(value), ";")
		case SignatureEntryKey:
			auth.Signature = value
		}
	}
	if s := auth.Scheme; s != AuthScheme {
		return nil, fmt.Errorf("%w: scheme=%s", ErrUnknownAuthScheme, s)
	}
	return auth, nil
}

func newQueryAuthorization(query url.Values) (*Authorization, error) {
	cred, err := NewCredential(query.Get(CredentialQueryKey))
	if err != nil {
		return nil, err
	}
	return &Authorization{
		Scheme:        AuthScheme,
		Credential:    cred,
		SignedHeaders: strings.Split(query.Get(SignedHeadersQueryKey), ":"),
		Signature:     query.Get(SignatureQueryKey),
		ContentHash:   EmptyContentHash,
	}, nil
}

type Credential struct {
	AccessKey api.AccessKey
	Date      string
	Region    string
	Service   string
	Suffix    string
}

func NewCredential(raw string) (*Credential, error) {
	const fieldsN = 5
	fields := strings.SplitN(raw, "/", fieldsN)
	if l := len(fields); l != fieldsN {
		return nil, fmt.Errorf("%w: len=%d", ErrIncompleteCredentialFields, l)
	}
	return &Credential{
		AccessKey: api.AccessKey(fields[0]),
		Date:      fields[1],
		Region:    fields[2],
		Service:   fields[3],
		Suffix:    fields[4],
	}, nil
}
