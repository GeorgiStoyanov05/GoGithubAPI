package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const userURL = "https://api.github.com/users/"
const userRepoURL = "https://api.github.com/repos/"

func GetUser(username string) (User, error) {
	token := os.Getenv("GITHUB_TOKEN")
	req, err := http.NewRequest(http.MethodGet, userURL+username, nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "my-github-api-app")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserRepos(username string) ([]UserRepo, error) {
	token := os.Getenv("GITHUB_TOKEN")
	req, err := http.NewRequest(http.MethodGet, userURL+username+"/repos", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "my-github-api-app")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var repos []UserRepo
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func GetRepoLanguages(username string, repoName string) (map[string]int, error) {
	token := os.Getenv("GITHUB_TOKEN")
	req, err := http.NewRequest(http.MethodGet, userRepoURL+username+"/"+repoName+"/languages", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "my-github-api-app")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	langs := make(map[string]int)
	err = json.NewDecoder(res.Body).Decode(&langs)
	if err != nil {
		return nil, err
	}
	return langs, nil
}

func GetLangDistribution(username string) (map[string]float32, error) {
	totalUsages := 0
	languages := make(map[string]float32)
	repos, err := GetUserRepos(username)
	if err != nil {
		return nil, err
	}
	for i := range len(repos) {
		repoLangs, err := GetRepoLanguages(username, repos[i].Name)
		if err != nil {
			return nil, err
		}
		for k, v := range repoLangs {
			totalUsages += v
			languages[k] += float32(v)
		}
	}
	for k, _ := range languages {
		languages[k] = (languages[k] / float32(totalUsages)) * 100
	}

	return languages, nil
}

func GetRepoForks(username string) (map[string]int, error) {
	forksCount := make(map[string]int)
	repos, err := GetUserRepos(username)
	if err != nil {
		return nil, err
	}
	for i := range len(repos) {
		_, ok := forksCount[repos[i].Name]
		if !ok {
			forksCount[repos[i].Name] = 0
		}
		forksCount[repos[i].Name] = repos[i].Forks
	}

	return forksCount, nil
}

func GetUserActivity(username string) (map[int][]int, error) {
	res := make(map[int][]int)
	repos, err := GetUserRepos(username)
	if err != nil {
		return nil, err
	}
	for i := range len(repos) {
		creationYear := repos[i].CreatedAt.Year()
		lastUpdateYear := repos[i].UpdatedAt.Year()
		_, ok := res[creationYear]
		if !ok {
			res[creationYear] = make([]int, 3)
		}
		_, ok = res[lastUpdateYear]
		if !ok {
			res[lastUpdateYear] = make([]int, 3)
		}
		res[creationYear][0]++
		res[creationYear][2]++
		res[lastUpdateYear][1]++
		res[lastUpdateYear][2]++
	}

	return res, nil
}

func GetUserReport(username string) (UserReport, error) {
	user, err := GetUser(username)

	if err != nil {
		return UserReport{}, err
	}

	var report UserReport
	report.Name = user.Name
	report.Username = user.Login
	report.Bio = user.Bio
	report.Followers = user.Followers
	report.Following = user.Following
	report.PublicRepos = user.PublicRepos
	report.LanguageDistribution, err = GetLangDistribution(user.Login)

	if err != nil {
		return UserReport{}, err
	}
	report.RepoForks, err = GetRepoForks(user.Login)
	if err != nil {
		return UserReport{}, err
	}
	report.Activity, err = GetUserActivity(user.Login)
	if err != nil {
		return UserReport{}, err
	}
	return report, nil
}
