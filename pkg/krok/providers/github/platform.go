package github

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/webhooks.v5/github"
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

// ValidateRequest will take a hook and verify it being a valid hook request according to
// Github's rules.
func (g *Github) ValidateRequest(ctx context.Context, req *http.Request) error {
	req.Header.Set("Content-type", "application/json")
	defer req.Body.Close()

	// Get the secret from the repo auth provider?
	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))
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
		log.Debug().Msg("All good, send back ping.")
		return nil
	}
	return nil
}
