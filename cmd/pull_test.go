package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v60/github"
)

func TestPullPaginationWithMockServer(t *testing.T) {
	generateRepoJSON := func(startID int) string {
		result := "["
		for i := 0; i < 100; i++ {
			id := startID + i
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf(`{"id": %d, "name": "repo%d", "clone_url": "https://github.com/test/repo%d.git", "owner": {"login": "test"}}`, id, id, id)
		}
		result += "]"
		return result
	}

	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		switch requestCount {
		case 1:
			w.Header().Set("Link", `<http://test.test?page=2>; rel="next"`)
			fmt.Fprint(w, generateRepoJSON(1))
		case 2:
			w.Header().Set("Link", `<http://test.test?page=3>; rel="next"`)
			fmt.Fprint(w, generateRepoJSON(101))
		case 3:
			w.Header().Set("Link", `<http://test.test?page=4>; rel="next"`)
			fmt.Fprint(w, generateRepoJSON(201))
		case 4:
			w.Header().Set("Link", `<http://test.test?page=5>; rel="next"`)
			fmt.Fprint(w, generateRepoJSON(301))
		case 5:
			fmt.Fprint(w, generateRepoJSON(401))
		}
	}))
	defer server.Close()

	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	var allRepos []*github.Repository
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Type: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(context.Background(), opts)
		if err != nil {
			t.Fatalf("Error listing repositories: %v", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	if len(allRepos) != 500 {
		t.Errorf("Expected 500 repos, got %d", len(allRepos))
	}

	if requestCount != 5 {
		t.Errorf("Expected 5 API requests, got %d", requestCount)
	}

	t.Logf("✓ Paginazione funziona: %d request, %d repo totali", requestCount, len(allRepos))
}
