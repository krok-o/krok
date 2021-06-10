package gitlab

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
	ggitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/go-playground/webhooks.v5/gitlab"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger                zerolog.Logger
	PlatformTokenProvider providers.PlatformTokenProvider
	AuthProvider          providers.RepositoryAuth
	UUIDGenerator         providers.UUIDGenerator
}

// Gitlab is a gitlab based platform implementation.
type Gitlab struct {
	Dependencies

	// Test client for the gitlab client.
	httpClient *http.Client
}

// NewGitlabPlatformProvider creates a new hook platform provider for Gitlab.
func NewGitlabPlatformProvider(deps Dependencies) *Gitlab {
	return &Gitlab{Dependencies: deps}
}

var _ providers.Platform = &Gitlab{}

// ValidateRequest will take a hook and verify it being a valid hook request according to Gitlab's rules.
func (g *Gitlab) ValidateRequest(ctx context.Context, req *http.Request, repoID int) error {
	req.Header.Set("Content-type", "application/json")
	repoAuth, err := g.AuthProvider.GetRepositoryAuth(ctx, repoID)
	if err != nil {
		g.Logger.Debug().Err(err).Msg("Failed to get Repository Auth information.")
		return err
	}

	if repoAuth == nil {
		g.Logger.Debug().Msg("Auth is not present.")
		return errors.New("no auth specified")
	}

	hook, _ := gitlab.New(gitlab.Options.Secret(repoAuth.Secret))
	_, err = hook.Parse(req,
		gitlab.BuildEvents,
		gitlab.CommentEvents,
		gitlab.ConfidentialIssuesEvents,
		gitlab.IssuesEvents,
		gitlab.JobEvents,
		gitlab.MergeRequestEvents,
		gitlab.PipelineEvents,
		gitlab.PushEvents,
		gitlab.SystemHookEvents,
		gitlab.TagEvents,
		gitlab.WikiPageEvents)
	if err != nil {
		g.Logger.Debug().Err(err).Msg("Failed to parse gitlab event.")
		return err
	}
	return nil
}

// GetEventID Based on the platform, retrieve the ID of the event.
func (g *Gitlab) GetEventID(ctx context.Context, r *http.Request) (string, error) {
	return g.UUIDGenerator.Generate()
}

// CreateHook can create a hook for the Gitlab platform.
func (g *Gitlab) CreateHook(ctx context.Context, repo *models.Repository) error {
	log := g.Logger.With().Str("unique_url", repo.UniqueURL).Str("repo", repo.Name).Strs("events", repo.Events).Logger()
	token, err := g.PlatformTokenProvider.GetTokenForPlatform(repo.VCS)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get platform token.")
		return err
	}
	if repo.Auth == nil {
		log.Error().Msg("No auth provided for the repository.")
		return errors.New("no auth provided with the repository")
	}
	if repo.Auth.Secret == "" {
		log.Error().Msg("No secret provided for the repository.")
		return errors.New("no secret provided to create a hook")
	}
	if len(repo.Events) == 0 {
		log.Error().Msg("No events provided to subscribe to.")
		return errors.New("no events provided to subscribe to")
	}
	if repo.UniqueURL == "" {
		log.Error().Msg("Unique callback url is empty.")
		return errors.New("unique callback url is empty")
	}
	if repo.GitLab.GetProjectID() == nil {
		log.Error().Msg("Project ID must not be empty for a gitlab repository.")
	}

	var opts []ggitlab.ClientOptionFunc
	if g.httpClient != nil {
		opts = append(opts, ggitlab.WithHTTPClient(g.httpClient))
	}
	git, err := ggitlab.NewClient(token, opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create gitlab client.")
		return errors.New("failed to create gitlab client")
	}
	hookOpts := &ggitlab.AddProjectHookOptions{
		Token: &repo.Auth.Secret,
		URL:   &repo.UniqueURL,
	}

	for _, event := range repo.Events {
		switch event {
		case "ConfidentialNoteEvents":
			hookOpts.ConfidentialIssuesEvents = ggitlab.Bool(true)
		case "PushEvents":
			hookOpts.PushEvents = ggitlab.Bool(true)
		case "IssuesEvents":
			hookOpts.IssuesEvents = ggitlab.Bool(true)
		case "ConfidentialIssuesEvents":
			hookOpts.ConfidentialIssuesEvents = ggitlab.Bool(true)
		case "MergeRequestsEvents":
			hookOpts.MergeRequestsEvents = ggitlab.Bool(true)
		case "TagPushEvents":
			hookOpts.TagPushEvents = ggitlab.Bool(true)
		case "NoteEvents":
			hookOpts.NoteEvents = ggitlab.Bool(true)
		case "JobEvents":
			hookOpts.JobEvents = ggitlab.Bool(true)
		case "PipelineEvents":
			hookOpts.PipelineEvents = ggitlab.Bool(true)
		case "WikiPageEvents":
			hookOpts.WikiPageEvents = ggitlab.Bool(true)
		case "DeploymentEvents":
			hookOpts.DeploymentEvents = ggitlab.Bool(true)
		case "ReleasesEvents":
			hookOpts.ReleasesEvents = ggitlab.Bool(true)
		default:
			return fmt.Errorf("invalid event type %q", event)
		}
	}
	// TODO: Project ID can either be a name or an ID. Consider trying to match that. I can store
	// an integer serialized to bytes in the DB. https://golang.org/pkg/encoding/binary/#example_Write
	pid := *repo.GitLab.GetProjectID()
	hook, response, err := git.Projects.AddProjectHook(pid, hookOpts)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create hook.")
		log.Debug().Int("code", response.StatusCode).Msg("Status code of the response.")
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to read response body.")
			return err
		}
		log.Debug().Str("body", string(body)).Msg("The body of the response.")
		return err
	}
	log.Debug().Str("url", hook.URL).Msg("Successfully created hook for gitlab!")
	return nil
}
