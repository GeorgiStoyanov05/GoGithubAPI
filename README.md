# GitHub User Stats Reporter (Go)

A small Go CLI that reads a list of GitHub usernames and prints a side-by-side statistics report by calling the public GitHub REST API.

## What it does

Given a text file (one username per line), the tool:
- Fetches user info: `GET https://api.github.com/users/{username}`
- Fetches repositories: `GET https://api.github.com/users/{username}/repos`
- Fetches language usage per repo: `GET https://api.github.com/repos/{username}/{repo}/languages`
- Parses JSON responses with `encoding/json` (`json.Unmarshal`) into Go structs
- Prints a console report as structured tables so multiple users can be compared easily

## Report includes

- Basic user info (e.g. login, name, public stats)
- Total number of repositories
- Total followers
- Total forks across all repositories
- Programming language distribution (based on the `/languages` endpoint values)
- Activity by year (derived from repo creation and last update timestamps)

## Requirements

- Go (recommended: 1.20+)
- Internet access (public GitHub API)

## Quick start
Run (file path as the first argument):
```bash
go run ./cmd/main ./test.txt
```

Example `test.txt`:
```txt
torvalds
golang
kubernetes
```

## Optional: higher API rate limits

GitHub applies stricter rate limits to unauthenticated requests. If your report hits limits (especially when pulling repo languages), use a token:

```bash
export GITHUB_TOKEN=YOUR_TOKEN_HERE
./gh-report ./usernames.txt
```

*(Token should have minimal scopes; public read is enough for public data.)*

## Output

The program prints a console report (tables) with one row per user, plus breakdown sections (languages + yearly activity).  
Exact formatting/columns depend on the implementation, but the intent is: quick comparison across users.

## Notes

- This tool uses **public** GitHub API endpoints and only reads **public** data.
- GitHub paginates some endpoints (like repo lists). If you see fewer repos than expected for large accounts, make sure pagination/per-page handling is enabled in the client.
- Language stats come from the GitHub `/languages` endpoint (values represent language “size” in bytes for that repo).

