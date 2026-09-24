package inputs

import (
	"fmt"
	"os"
	"strings"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
)

// Load reads GitHub Action inputs from INPUT_* env vars (Docker/composite convention)
// with fallbacks to standard GitHub Actions environment variables.
func Load() (model.Inputs, error) {
	webhookURL := strings.TrimSpace(getInput("webhook_url"))
	tag := strings.TrimSpace(getInput("tag"))
	notifyPrerelease, err := parseBool(getInput("notify_prerelease"), false)
	if err != nil {
		return model.Inputs{}, err
	}

	githubToken := strings.TrimSpace(getInput("github_token"))
	if githubToken == "" {
		githubToken = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}

	repository := strings.TrimSpace(getInput("repository"))
	if repository == "" {
		repository = strings.TrimSpace(os.Getenv("GITHUB_REPOSITORY"))
	}
	if repository == "" {
		return model.Inputs{}, fmt.Errorf("repository is required (set the repository input or GITHUB_REPOSITORY)")
	}

	return model.Inputs{
		WebhookURL:       webhookURL,
		Tag:              tag,
		NotifyPrerelease: notifyPrerelease,
		GitHubToken:      githubToken,
		Repository:       repository,
		EventPath:        strings.TrimSpace(os.Getenv("GITHUB_EVENT_PATH")),
	}, nil
}

func getInput(name string) string {
	// GitHub Actions maps inputs to INPUT_<NAME> with uppercase and non-alnum → _.
	key := "INPUT_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
	return os.Getenv(key)
}

func parseBool(value string, defaultValue bool) (bool, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return defaultValue, nil
	}
	switch normalized {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean input value: %q", value)
	}
}
