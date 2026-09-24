package release

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
)

type githubReleaseResponse struct {
	TagName string  `json:"tag_name"`
	Name    *string `json:"name"`
	Body    *string `json:"body"`
	HTMLURL string  `json:"html_url"`
	Author  *struct {
		Login string `json:"login"`
	} `json:"author"`
	Prerelease bool   `json:"prerelease"`
	Message    string `json:"message"`
}

type releaseEventPayload struct {
	Release *githubReleaseResponse `json:"release"`
}

// Resolver loads release metadata from a tag (API) or the Actions event payload.
type Resolver struct {
	HTTPClient *http.Client
}

// Resolve returns release metadata for the given options.
func (r *Resolver) Resolve(tag, repository, githubToken, eventPath string) (model.ReleaseInfo, error) {
	if strings.TrimSpace(tag) != "" {
		return r.fetchByTag(strings.TrimSpace(tag), repository, githubToken)
	}
	return readFromEvent(eventPath)
}

func (r *Resolver) fetchByTag(tag, repository, githubToken string) (model.ReleaseInfo, error) {
	if strings.TrimSpace(githubToken) == "" {
		return model.ReleaseInfo{}, fmt.Errorf("github_token is required when tag is provided so the release can be fetched from the GitHub API")
	}
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return model.ReleaseInfo{}, err
	}

	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases/tags/%s",
		url.PathEscape(owner),
		url.PathEscape(repo),
		url.PathEscape(tag),
	)

	client := r.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return model.ReleaseInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "zoho-cliq-release-notifier")

	resp, err := client.Do(req)
	if err != nil {
		return model.ReleaseInfo{}, fmt.Errorf("failed to fetch release for tag %q in %s: %w", tag, repository, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return model.ReleaseInfo{}, fmt.Errorf("failed to fetch release for tag %q in %s: %w", tag, repository, err)
	}

	var data githubReleaseResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return model.ReleaseInfo{}, fmt.Errorf("failed to fetch release for tag %q in %s: invalid JSON (HTTP %d)", tag, repository, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		detail := data.Message
		if detail == "" {
			detail = string(body)
			if len(detail) > 200 {
				detail = detail[:200]
			}
		}
		return model.ReleaseInfo{}, fmt.Errorf("failed to fetch release for tag %q in %s: HTTP %d %s", tag, repository, resp.StatusCode, detail)
	}
	return normalize(data)
}

func readFromEvent(eventPath string) (model.ReleaseInfo, error) {
	if eventPath == "" {
		return model.ReleaseInfo{}, fmt.Errorf("missing release payload; trigger via release published or pass the tag input to fetch via API")
	}
	raw, err := os.ReadFile(eventPath)
	if err != nil {
		return model.ReleaseInfo{}, fmt.Errorf("missing release payload; trigger via release published or pass the tag input to fetch via API")
	}

	var payload releaseEventPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return model.ReleaseInfo{}, fmt.Errorf("malformed GitHub event payload at %s: %w", eventPath, err)
	}
	if payload.Release == nil || strings.TrimSpace(payload.Release.TagName) == "" {
		return model.ReleaseInfo{}, fmt.Errorf("missing release payload; trigger via release published or pass the tag input to fetch via API")
	}
	return normalize(*payload.Release)
}

func normalize(release githubReleaseResponse) (model.ReleaseInfo, error) {
	tagName := strings.TrimSpace(release.TagName)
	if tagName == "" {
		return model.ReleaseInfo{}, fmt.Errorf("release payload is missing tag_name")
	}
	name := ""
	if release.Name != nil {
		name = *release.Name
	}
	body := ""
	if release.Body != nil {
		body = *release.Body
	}
	author := ""
	if release.Author != nil {
		author = release.Author.Login
	}
	return model.ReleaseInfo{
		TagName:     tagName,
		Name:        name,
		Body:        body,
		HTMLURL:     release.HTMLURL,
		AuthorLogin: author,
		Prerelease:  release.Prerelease,
	}, nil
}

func splitRepository(repository string) (string, string, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repository %q; expected owner/repo", repository)
	}
	return parts[0], parts[1], nil
}
