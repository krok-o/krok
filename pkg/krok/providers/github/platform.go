package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	ggithub "github.com/google/go-github/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/webhooks.v5/github"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
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
	Hostname string
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger                zerolog.Logger
	PlatformTokenProvider providers.PlatformTokenProvider
	AuthProvider          providers.RepositoryAuth
}

// Github is a github based platform implementation.
type Github struct {
	Config
	Dependencies

	// Used for testing the CreateHook call. There probably is a better way to do this...
	repoMock GoogleGithubRepoService
}

// NewGithubPlatformProvider creates a new hook platform provider for Github.
func NewGithubPlatformProvider(cfg Config, deps Dependencies) *Github {
	return &Github{Config: cfg, Dependencies: deps}
}

var _ providers.Platform = &Github{}

// ValidateRequest will take a hook and verify it being a valid hook request according to
// Github's rules.
func (g *Github) ValidateRequest(ctx context.Context, req *http.Request, repoID int) error {
	req.Header.Set("Content-type", "application/json")
	defer req.Body.Close()

	repoAuth, err := g.AuthProvider.GetRepositoryAuth(ctx, repoID)
	if err != nil {
		g.Logger.Debug().Err(err).Msg("Failed to get Repository Auth information.")
		return err
	}

	// Get the secret from the repo auth provider?
	hook, _ := github.New(github.Options.Secret(repoAuth.Secret))
	h, err := hook.Parse(req,
		github.CheckRunEvent,
		github.CheckSuiteEvent,
		github.CommitCommentEvent,
		github.CreateEvent,
		github.DeleteEvent,
		github.DeploymentEvent,
		github.DeploymentStatusEvent,
		github.ForkEvent,
		github.GollumEvent,
		github.InstallationEvent,
		github.InstallationRepositoriesEvent,
		github.IntegrationInstallationRepositoriesEvent,
		github.IssueCommentEvent,
		github.IssuesEvent,
		github.LabelEvent,
		github.MemberEvent,
		github.MembershipEvent,
		github.MetaEvent,
		github.MilestoneEvent,
		github.OrganizationEvent,
		github.OrgBlockEvent,
		github.PageBuildEvent,
		github.PingEvent,
		github.ProjectCardEvent,
		github.ProjectColumnEvent,
		github.ProjectEvent,
		github.PublicEvent,
		github.PullRequestEvent,
		github.PullRequestReviewCommentEvent,
		github.PullRequestReviewEvent,
		github.ReleaseEvent,
		github.RepositoryEvent,
		github.RepositoryVulnerabilityAlertEvent,
		github.SecurityAdvisoryEvent,
		github.StatusEvent)
	if err != nil {
		g.Logger.Debug().Err(err).Msg("Failed to parse github event.")
		return err
	}
	switch h.(type) {
	case github.PingPayload:
		g.Logger.Debug().Msg("All good, send back ping.")
		return nil
	}
	return nil
}

// GoogleGithubRepoService is an interface defining the Wrapper Interface
// needed to test the github client.
type GoogleGithubRepoService interface {
	CreateHook(ctx context.Context, owner, repo string, hook *ggithub.Hook) (*ggithub.Hook, *ggithub.Response, error)
}

// GoogleGithubClient is a client that has the ability to replace the actual
// git client.
type GoogleGithubClient struct {
	Repositories GoogleGithubRepoService
	*ggithub.Client
}

// NewGoogleGithubClient creates a wrapper around the github client. This is
// needed in order to decouple gaia from github client to be
// able to unit test createGithubWebhook and ultimately have
// the ability to replace github with anything else.
func NewGoogleGithubClient(httpClient *http.Client, repoMock GoogleGithubRepoService) GoogleGithubClient {
	if repoMock != nil {
		return GoogleGithubClient{
			Repositories: repoMock,
		}
	}
	githubClient := ggithub.NewClient(httpClient)

	return GoogleGithubClient{
		Repositories: githubClient.Repositories,
	}
}

// CreateHook can create a hook for the Github platform.
func (g *Github) CreateHook(ctx context.Context, repo *models.Repository) error {
	log := g.Logger.With().Str("repo", repo.Name).Strs("events", repo.Events).Logger()
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
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	config := make(map[string]interface{})
	config["url"] = repo.UniqueURL
	config["secret"] = repo.Auth.Secret
	config["content_type"] = "json"

	// figure out a way to mock this nicely later on.
	githubClient := NewGoogleGithubClient(tc, g.repoMock)
	repoName := path.Base(repo.URL)
	repoName = strings.TrimSuffix(repoName, ".git")
	// var repoLocation string
	re := regexp.MustCompile("^(https|git)(://|@)([^/:]+)[/:]([^/:]+)/(.+)$")
	m := re.FindAllStringSubmatch(repo.URL, -1)
	if m == nil {
		log.Debug().Str("repo_name", repoName).Str("url", repo.URL).Msg("Failed to extract url parameters.")
		return errors.New("failed to extract url parameters from git url")
	}
	repoUser := m[0][4]
	hook, resp, err := githubClient.Repositories.CreateHook(context.Background(), repoUser, repoName, &ggithub.Hook{
		Events: repo.Events,
		Name:   ggithub.String("web"),
		Active: ggithub.Bool(true),
		Config: config,
	})
	if err != nil {
		log.Debug().Err(err).Msg("CreateHook failed.")
		return err
	}
	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		log.Error().Msg("invalid status code")
		return fmt.Errorf("invalid status code %d received from hook creation", resp.StatusCode)
	}
	g.Logger.Debug().Str("name", *hook.Name).Msg("Hook with name successfully created.")
	return nil
}
