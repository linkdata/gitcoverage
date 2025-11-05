[![coverage](https://github.com/linkdata/gitcoverage/blob/main/coverage_badge_animated.svg)](#)

# gitcoverage

Generate code coverage badge and push it and optional HTML report to the 'coverage' branch.

This action has no dependencies except for `git`, the `bash` shell and common *nix command line utilities
`awk`, `sed` and GNU coreutils (`mkdir, cp, rm, ls, cat, echo, printf`). This means it won't run on Windows
runners; use `if: runner.os != 'Windows'` to exclude those in the workflow.

## Usage

You need to have given write permissions for the for the workflow.
If the 'coverage' branch does not exist, it will be created as an orphan (without main repo history).
Reference the generated badge in your README.md like this:

```md
[![coverage](https://github.com/USERNAME/REPO/blob/coverage/BRANCH/badge.svg)](#)
```

If you submitted a detailed HTML report of the coverage to the action, replace the '#' with:

`https://htmlpreview.github.io/?https://github.com/USERNAME/REPO/blob/coverage/BRANCH/report.html`

## Examples

Inside your .github/workflows/workflow.yml file:

```yml
permissions:
  contents: write

jobs:
  build:
    steps:
      - uses: actions/checkout@v4
      - uses: linkdata/gitcoverage@v3
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
        if: runner.os != 'Windows'
        id: coverage
        run: |
          echo "COVERAGE=$(go tool cover -func=coverage | tail -n 1 | tr -s '\t' | cut -f 3)" >> $GITHUB_OUTPUT
          go tool cover -html=coverage -o=coveragereport.html

      - name: Publish code coverage badge (and optional report)
        if: runner.os != 'Windows'
        uses: linkdata/gitcoverage@v3
        with:
          coverage: ${{ steps.coverage.outputs.coverage }}
          report:   "coveragereport.html"
```
