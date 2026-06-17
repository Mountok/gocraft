package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mountok/gocraft/internal/version"
)

const latestTagURL = "https://api.github.com/repos/Mountok/gocraft/tags?per_page=1"

type tag struct {
	Name string `json:"name"`
}

type Result struct {
	Current  string
	Latest   string
	Outdated bool
}

func Check(ctx context.Context) (Result, error) {
	current := version.Current()
	latest, err := Latest(ctx)
	if err != nil {
		return Result{Current: current}, err
	}

	return Result{
		Current:  current,
		Latest:   latest,
		Outdated: isOutdated(current, latest),
	}, nil
}

func Latest(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestTagURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gocraft")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("GitHub returned %s", resp.Status)
	}

	var tags []tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}
	if len(tags) == 0 || tags[0].Name == "" {
		return "", fmt.Errorf("no release tags found")
	}
	return tags[0].Name, nil
}

func WarnIfOutdated(output interface{ Write([]byte) (int, error) }) {
	if strings.TrimSpace(strings.ToLower(getenv("GOCRAFT_SKIP_UPDATE_CHECK"))) == "1" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := Check(ctx)
	if err != nil || !result.Outdated {
		return
	}

	fmt.Fprintf(output, "A new GoCraft version is available: %s -> %s\n", result.Current, result.Latest)
	fmt.Fprintln(output, "Update: curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/install.sh | sh")
}

var getenv = func(key string) string {
	return os.Getenv(key)
}

func isOutdated(current, latest string) bool {
	current = normalize(current)
	latest = normalize(latest)
	if current == "" || current == "dev" || current == latest {
		return false
	}
	return compareSemver(current, latest) < 0
}

func normalize(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimPrefix(value, "v")
	return value
}

func compareSemver(a, b string) int {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		av := part(aparts, i)
		bv := part(bparts, i)
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}
	return 0
}

func part(parts []string, index int) int {
	if index >= len(parts) {
		return 0
	}
	value := 0
	for _, r := range parts[index] {
		if r < '0' || r > '9' {
			break
		}
		value = value*10 + int(r-'0')
	}
	return value
}
