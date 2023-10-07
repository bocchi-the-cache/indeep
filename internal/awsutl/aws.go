package awsutl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/httputl"
)

const (
	ShortTimeFormat = "20060102"

	CredentialEntryKey    = "Credential"
	SignedHeadersEntryKey = "SignedHeaders"
	SignatureEntryKey     = "Signature"
)

var (
	ErrEmptyCredential            = errors.New("empty AWS credential")
	ErrEmptySignedHeaders         = errors.New("empty AWS signed headers")
	ErrEmptySignature             = errors.New("empty AWS signature")
	ErrIncompleteCredentialFields = errors.New("incomplete AWS credential fields")
	ErrCredentialParseTime        = errors.New("parse AWS credential time error")
)

type Authorization struct {
	Scheme        string
	Credential    *Credential
	SignedHeaders []string
	Signature     string
}

func NewAuthorization(a *httputl.Authorization) (*Authorization, error) {
	ret := &Authorization{Scheme: a.Scheme}
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
			ret.Credential = cred
		case SignedHeadersEntryKey:
			ret.SignedHeaders = strings.Split(value, ";")
		case SignatureEntryKey:
			ret.Signature = value
		}
	}
	if ret.Credential == nil {
		return nil, ErrEmptyCredential
	}
	if ret.SignedHeaders == nil {
		return nil, ErrEmptySignedHeaders
	}
	if ret.Signature == "" {
		return nil, ErrEmptySignature
	}
	return ret, nil
}

type Credential struct {
	AccessKey api.AccessKey
	Date      time.Time
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

	rawDate := fields[1]
	t, err := time.Parse(ShortTimeFormat, rawDate)
	if err != nil {
		return nil, fmt.Errorf("%w: err=%v", ErrCredentialParseTime, err)
	}

	return &Credential{
		AccessKey: api.AccessKey(fields[0]),
		Date:      t,
		Region:    fields[2],
		Service:   fields[3],
		Suffix:    fields[4],
	}, nil
}
