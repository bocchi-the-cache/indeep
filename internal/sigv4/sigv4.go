package sigv4

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	sigv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/awsutl"
	"github.com/bocchi-the-cache/indeep/internal/httputl"
)

const (
	AuthScheme = "AWS4-HMAC-SHA256"

	ContentHashKey = "X-Amz-Content-Sha256"
	TimeKey        = "X-Amz-Date"
	TimeFormat     = "20060102T150405Z"
)

type Input struct {
	Auth        *awsutl.Authorization
	Credentials aws.Credentials
	ContentHash string
	Time        time.Time
}

func NewInput(tenants api.Tenants, r *http.Request) (*Input, error) {
	auth, err := httputl.NewAuthorization(r)
	if err != nil {
		return nil, err
	}
	if s := auth.Scheme; s != AuthScheme {
		return nil, fmt.Errorf("%w: scheme=%s", api.ErrUnknownAuthScheme, s)
	}
	awsAuth, err := awsutl.NewAuthorization(auth)
	if err != nil {
		return nil, err
	}

	sk, err := tenants.SecretKey(awsAuth.Credential.AccessKey)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(TimeFormat, r.Header.Get(TimeKey))
	if err != nil {
		return nil, err
	}

	return &Input{
		Auth: awsAuth,
		Credentials: aws.Credentials{
			AccessKeyID:     string(awsAuth.Credential.AccessKey),
			SecretAccessKey: string(sk),
		},
		ContentHash: r.Header.Get(ContentHashKey),
		Time:        t,
	}, nil
}

type checker struct{ tenants api.Tenants }

func New(tenants api.Tenants) api.SigV4Checker { return &checker{tenants} }

func (c *checker) CheckSigV4(r *http.Request) error {
	input, err := NewInput(c.tenants, r)
	if err != nil {
		return err
	}
	return sigv4.NewSigner().SignHTTP(
		context.Background(),
		input.Credentials,
		r,
		input.ContentHash,
		input.Auth.Credential.Service,
		input.Auth.Credential.Region,
		input.Time,
	)
}
