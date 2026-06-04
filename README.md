[![coverage](https://github.com/linkdata/gitcoverage/blob/main/coverage_badge_animated.svg)](#)

# gitcoverage

Generate code coverage badge and push it and optional HTML report to the 'gitcoverage' branch.

This action has no dependencies except for `git`, a `bash` shell and common *nix command line utilities
`awk`, `sed` and GNU coreutils (`mkdir, cp, rm, ls, cat, echo, printf`).

It supports Linux/macOS runners and Windows runners with Bash tooling (Git Bash/WSL-enabled images such as
`windows-2025`).

Requires **Git 2.15.0 or newer** (the action fails fast on older versions).

## Usage

You need to have given write permissions for the for the workflow job that runs this action.
If the 'gitcoverage' branch does not exist, it will be created as an orphan (without main repo history).
The action creates bot commits with signing disabled (`commit.gpgsign=false`) for compatibility with runners that enforce local signing config but have no key.
If your `gitcoverage` branch requires signed commits, configure signing keys on the runner or relax that branch rule.
By default the action pushes using the credentials that `actions/checkout` persists.
Pass the `token` input to authenticate explicitly instead, which lets you check out with `persist-credentials: false`.
When credentials are still persisted (`persist-credentials: true`), those take precedence over the `token` input.
Reference the generated badge in your README.md like this:

```md
[![coverage](https://github.com/USERNAME/REPO/blob/gitcoverage/BRANCH/badge.svg)](#)
```

If you submitted a detailed HTML report of the coverage to the action, replace the '#' with:

`https://html-preview.github.io/?url=https://github.com/USERNAME/REPO/blob/gitcoverage/BRANCH/report.html`

### Inputs

- `coverage` (required): Coverage percentage (for example `83` or `83%`).
- `report` (optional): Path to an HTML report file to publish as `report.html`.
- `token` (optional): GitHub token used to push updates to the `gitcoverage` branch.
  When set, the action configures a temporary Git credential helper, so you can check out with `actions/checkout` and `persist-credentials: false`.
  When omitted, the action uses the credentials persisted by `actions/checkout`.
- `branch` (optional): Source branch override. Recommended for tag-triggered workflows where multiple branches may contain the same tag commit.
  Also recommended for very large or restricted repos to avoid scanning all remote branches during tag-triggered branch resolution.
  On Windows runners, the action applies a strict compatibility filter and requires branch names to match `[A-Za-z0-9._/+-]+`.
  This filter does not reject Windows-reserved path components such as `CON`, `NUL`, `AUX`, `COM1`, or `LPT9`; avoid those names on Windows runners.

## Examples

Inside your .github/workflows/workflow.yml file:

```yml
permissions:
  contents: read

jobs:
  coverage:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v6
        with:
          persist-credentials: false
      - uses: linkdata/gitcoverage@v10
        with:
          coverage: "83%"
          report:   "coveragereport.html.out"
          token:    ${{ github.token }}
```

More complete example using Go:

```yml
name: build

permissions:
  contents: read

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    outputs:
      coverage: ${{ steps.coverage.outputs.coverage }}
    steps:
      - uses: actions/checkout@v6
        with:
          persist-credentials: false

      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: stable

      - name: Generate
        run: go generate ./...

      - name: Go vet
        run: go vet ./...

      - name: Check gofmt
        run: |
          unformatted="$(gofmt -l .)"
          if [ -n "$unformatted" ]; then
            echo "The following files are not gofmt-formatted:"
            echo "$unformatted"
            exit 1
          fi

      - name: Staticcheck
        run: |
          go install honnef.co/go/tools/cmd/staticcheck@latest
          staticcheck ./...

      - name: golangci-lint
        uses: golangci/golangci-lint-action@latest

      - name: Run Gosec Security Scanner
        uses: securego/gosec@v2.26.1
        with:
          args: ./...

      - name: Test
        run: go test -race -coverprofile=coverage.out ./...

      - name: Calculate code coverage
        id: coverage
        run: |
          echo "coverage=$(go tool cover -func=coverage.out | tail -n 1 | tr -s '\t' | cut -f 3)" >> "$GITHUB_OUTPUT"
          go tool cover -html=coverage.out -o=coveragereport.html.out

      - name: Upload code coverage
        uses: actions/upload-artifact@v7
        with:
          name: coverage
          path: coveragereport.html.out
          retention-days: 1

      - name: Go report card
        uses: creekorful/goreportcard-action@v1.0
        continue-on-error: true

      - name: Build
        run: go build -v ./...

  coverage:
    needs: build
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v6
        with:
          persist-credentials: false

      - name: Download code coverage
        uses: actions/download-artifact@v8
        with:
          name: coverage

      - name: Publish code coverage badge (and optional report)
        uses: linkdata/gitcoverage@v10
        with:
          coverage: ${{ needs.build.outputs.coverage }}
          report:   "coveragereport.html.out"
          token:    ${{ github.token }}
```

Tag workflow example with explicit source branch:

```yml
- name: Publish code coverage badge from tag build
  uses: linkdata/gitcoverage@v10
  with:
    coverage: "91%"
    branch:   "release/1.x"
    token:    ${{ github.token }}
```
