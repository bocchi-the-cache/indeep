package sigv4

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/awsutl"
)

const (
	TimeKey              = "X-Amz-Date"
	SigningDateKeyPrefix = "AWS4"
)

type checker struct{ tenants api.Tenants }

func New(tenants api.Tenants) api.SigChecker { return &checker{tenants} }

func (c *checker) CheckSigV4(r *http.Request) (bool, error) {
	auth, err := awsutl.NewAuthorization(r)
	if err != nil {
		return false, err
	}

	sk, err := c.tenants.SecretKey(auth.Credential.AccessKey)
	if err != nil {
		return false, err
	}

	var canonicalHeaderList []string
	for _, key := range auth.SignedHeaders {
		canonicalHeaderList = append(
			canonicalHeaderList,
			fmt.Sprintf("%s:%s\n", key, strings.TrimSpace(r.Header.Get(key))),
		)
	}
	canonicalHeaders := strings.Join(canonicalHeaderList, "")

	canonicalRequest := strings.Join(
		[]string{
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			canonicalHeaders,
			strings.Join(auth.SignedHeaders, ":"),
			auth.ContentHash,
		},
		"\n",
	)
	h := sha256.New()
	if _, err := h.Write([]byte(canonicalRequest)); err != nil {
		return false, err
	}
	hashedCanonicalRequest := hex.EncodeToString(h.Sum(nil))

	date := auth.Credential.Date
	stringToSign := strings.Join(
		[]string{
			awsutl.AuthScheme,
			r.Header.Get(TimeKey),
			strings.Join(
				[]string{
					date,
					auth.Credential.Region,
					auth.Credential.Service,
					auth.Credential.Suffix,
				},
				"/",
			),
			hashedCanonicalRequest,
		},
		"\n",
	)

	dateKey := hmacSHA256([]byte(SigningDateKeyPrefix+string(sk)), date)
	dateRegionKey := hmacSHA256(dateKey, auth.Credential.Region)
	dateRegionServiceKey := hmacSHA256(dateRegionKey, auth.Credential.Service)
	signingKey := hmacSHA256(dateRegionServiceKey, auth.Credential.Suffix)
	sig := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	return sig == auth.Signature, nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
