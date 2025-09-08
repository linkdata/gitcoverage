# gitcoverage

Generate code coverage badge and push it and optional HTML report to the 'coverage' branch.
You need to have given write permissions for the for the workflow.
Reference the generated badge in your README.md like this:

```md
[![coverage](https://github.com/USERNAME/REPO/blob/coverage/BRANCH/badge.svg)](#)
```

If you have a detailed HTML report of the coverage (here called report.html), replace the '#' with:

`https://htmlpreview.github.io/?https://github.com/USERNAME/REPO/blob/coverage/BRANCH/report.html`

# Examples

Inside your .github/workflows/workflow.yml file:

```yml
permissions:
  contents: write

jobs:
  build:
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: linkdata/gitcoverage
        with:
          coverage: "83%"
          report:   "coveragereport.out"
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
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}

      - name: Go generate
        run: go generate ./...

      - name: Test
        run: go test -coverprofile=coverage.out ./...

      - name: Coverage
        id: coverage
        run: |
          echo "COVERAGE=$(go tool cover -func=coverage.out | tail -n 1 | tr -s '\t' | cut -f 3)" >> $GITHUB_OUTPUT
          go tool cover -html=coverage.out -o=coveragereport.out

      - name: Publish badge (and optional report)
        uses: linkdata/gitcoverage
        with:
          coverage: ${{ steps.coverage.outputs.coverage }}
          report:   "coveragereport.out"
```
