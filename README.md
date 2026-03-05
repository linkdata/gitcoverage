[![coverage](https://github.com/linkdata/gitcoverage/blob/main/coverage_badge_animated.svg)](#)

# gitcoverage

Generate code coverage badge and push it and optional HTML report to the 'coverage' branch.

This action has no dependencies except for `git`, a `bash` shell and common *nix command line utilities
`awk`, `sed` and GNU coreutils (`mkdir, cp, rm, ls, cat, echo, printf`).

It supports Linux/macOS runners and Windows runners with Bash tooling (Git Bash/WSL-enabled images such as
`windows-2025`).

Requires **Git 2.15.0 or newer** (the action fails fast on older versions).

## Usage

You need to have given write permissions for the for the workflow.
If the 'coverage' branch does not exist, it will be created as an orphan (without main repo history).
The action creates bot commits with signing disabled (`commit.gpgsign=false`) for compatibility with runners that enforce local signing config but have no key.
If your `coverage` branch requires signed commits, configure signing keys on the runner or relax that branch rule.
Reference the generated badge in your README.md like this:

```md
[![coverage](https://github.com/USERNAME/REPO/blob/coverage/BRANCH/badge.svg)](#)
```

If you submitted a detailed HTML report of the coverage to the action, replace the '#' with:

`https://html-preview.github.io/?url=https://github.com/USERNAME/REPO/blob/coverage/BRANCH/report.html`

### Inputs

- `coverage` (required): Coverage percentage (for example `83` or `83%`).
- `report` (optional): Path to an HTML report file to publish as `report.html`.
- `branch` (optional): Source branch override. Recommended for tag-triggered workflows where multiple branches may contain the same tag commit.
  Also recommended for very large or restricted repos to avoid scanning all remote branches during tag-triggered branch resolution.
  On Windows runners, the action applies a strict compatibility filter and requires branch names to match `[A-Za-z0-9._/+-]+`.

## Examples

Inside your .github/workflows/workflow.yml file:

```yml
permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: linkdata/gitcoverage@v6
        with:
          coverage: "83%"
          report:   "coveragereport.html.out"
```

More complete example using Go:

```yml
permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        go:
          - "stable"
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}

      - name: Go Generate
        run: go generate ./...

      - name: Go Test
        run: go test -coverprofile=coverage ./...

      - name: Go Build
        run: go build .

      - name: Calculate code coverage
        id: coverage
        run: |
          echo "COVERAGE=$(go tool cover -func=coverage | tail -n 1 | tr -s '\t' | cut -f 3)" >> $GITHUB_OUTPUT
          go tool cover -html=coverage -o=coveragereport.html

      - name: Publish code coverage badge (and optional report)
        uses: linkdata/gitcoverage@v6
        with:
          coverage: ${{ steps.coverage.outputs.coverage }}
          report:   "coveragereport.html"
```

Tag workflow example with explicit source branch:

```yml
- name: Publish code coverage badge from tag build
  uses: linkdata/gitcoverage@v6
  with:
    coverage: "91%"
    branch:   "release/1.x"
```
