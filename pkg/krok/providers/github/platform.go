package github

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

// Hook represent a github based webhook context.
type Hook struct {
	Signature string
	Event     string
	ID        string
	Payload   []byte
}

// Repository contains information about the repository. All we care about
// here are the possible urls for identification.
type Repository struct {
	GitURL  string `json:"git_url"`
	SSHURL  string `json:"ssh_url"`
	HTMLURL string `json:"html_url"`
}

// Payload contains information about the event like, user, commit id and so on.
// All we care about for the sake of identification is the repository.
type Payload struct {
	Repo Repository `json:"repository"`
}

// Config has the configuration options for the plugins.
type Config struct {
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
}

// Github is a github based platform implementation.
type Github struct {
	Config
	Dependencies
}

// NewGithubPlatformProvider creates a new hook platform provider for Github.
func NewGithubPlatformProvider(cfg Config, deps Dependencies) (*Github, error) {
	return &Github{Config: cfg, Dependencies: deps}, nil
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	_, _ = computed.Write(body)
	return computed.Sum(nil)
}

func verifySignature(secret []byte, signature string, body []byte) bool {
	signaturePrefix := "sha1="
	signatureLength := 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	_, _ = hex.Decode(actual, []byte(signature[5:]))
	expected := signBody(secret, body)
	return hmac.Equal(expected, actual)
}

func checkHeaders(req *http.Request) (Hook, error) {
	h := Hook{}

	if h.Signature = req.Header.Get("x-hub-signature"); len(h.Signature) == 0 {
		return Hook{}, errors.New("no signature")
	}

	if h.Event = req.Header.Get("x-github-event"); len(h.Event) == 0 {
		return Hook{}, errors.New("no event")
	}

	if h.Event != "push" {
		if h.Event == "ping" {
			return Hook{Event: "ping"}, nil
		}
		return Hook{}, errors.New("invalid event")
	}

	if h.ID = req.Header.Get("x-github-delivery"); len(h.ID) == 0 {
		return Hook{}, errors.New("no event id")
	}
	return h, nil
}

// ValidateRequest will take a hook and verify it being a valid hook request according to
// Github's rules.
func (g *Github) ValidateRequest(ctx context.Context, req *http.Request) error {
	req.Header.Set("Content-type", "application/json")
	defer req.Body.Close()

	h, err := checkHeaders(req)
	if err != nil {
		return err
	}
	if h.Event == "ping" {
		return err
	}

	p := Payload{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	h.Payload = body
	if err := json.Unmarshal(h.Payload, &p); err != nil {
		return err
	}

	// TODO: this secret needs to be configured somewhere... somehow. In vault, but do we attach it
	// to a repo or not? I guess a token in vault is fine...
	if !verifySignature(nil, h.Signature, h.Payload) {
		return err
	}
	return nil
}
