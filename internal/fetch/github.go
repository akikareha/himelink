package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GitHubFetchRepoInfo(baseURL, owner, repo string) (RepoInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", baseURL, owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return RepoInfo{}, err
	}
	defer resp.Body.Close()

	var info RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return RepoInfo{}, err
	}

	return info, nil
}

func GitHubFetchRaw(baseURL, owner, repo, branch, path string) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/%s/%s/%s/%s",
		baseURL, owner, repo, branch, path,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
