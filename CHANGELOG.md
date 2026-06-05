# Changelog

## v11

- Skip coverage badge publishing by default on `pull_request` and `pull_request_target` events. Skipped runs emit a notice and complete successfully without pushing.
- Add `run-on-pull-request` to let trusted workflows opt into attempting PR-context publishing.
- Update workflow examples to keep CI running on PRs while publishing badges only from default-branch push events.
