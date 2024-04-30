package deepinhelp

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/pluginhelp"
	"k8s.io/test-infra/prow/plugins"
)

const pluginName = "deepin-help"

var (
	helpRe = regexp.MustCompile(`(?mi)^/help\s*$`)
)

type prCommands struct {
	prCommnadsSummary string
}

func (ig prCommands) helpMsg() string {
	return fmt.Sprintf(`
### deepin pr commands help details.

%s`, ig.prCommnadsSummary)
}

func init() {
	plugins.RegisterGenericCommentHandler(pluginName, handleGenericComment, helpProvider)
}

func helpProvider(config *plugins.Configuration, _ []config.OrgRepo) (*pluginhelp.PluginHelp, error) {
	// The Config field is omitted because this plugin is not configurable.
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The help plugin provides commands for deepin github pr",
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/?)",
		Description: "Show deepin github pr commands",
		Featured:    false,
		WhoCanUse:   "Anyone can trigger this command on a PR.",
		Examples:    []string{"/?"},
	})
	return pluginHelp, nil
}

type githubClient interface {
	BotUserChecker() (func(candidate string) bool, error)
	CreateComment(owner, repo string, number int, comment string) error
}

func handleGenericComment(pc plugins.Agent, e github.GenericCommentEvent) error {
	cfg := pc.PluginConfig

	ig := prCommands{
		prCommnadsSummary: cfg.DeepinHelp.HelpCommandsSummary,
	}
	return handle(pc.GitHubClient, pc.Logger, &e, ig)
}

func handle(gc githubClient, log *logrus.Entry, e *github.GenericCommentEvent, ig prCommands) error {
	// Only consider open issues and new comments.
	if !e.IsPR || e.Action != github.GenericCommentActionCreated {
		return nil
	}

	org := e.Repo.Owner.Login
	repo := e.Repo.Name
	commentAuthor := e.User.Login

	// If PR does not have the help label and we're asking it to be added,
	// add the label
	if helpRe.MatchString(e.Body) {
		if err := gc.CreateComment(org, repo, e.Number, plugins.FormatResponseRaw(e.Body, e.HTMLURL, commentAuthor, ig.helpMsg())); err != nil {
			log.WithError(err).Errorf("Failed to create comment \"%s\".", ig.helpMsg())
		}
	}

	return nil
}
