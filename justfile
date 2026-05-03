set dotenv-load
set quiet


# Keep this at the top
_default:
  just --list
  echo '`just --choose` to select a task to run interactively'

# Default target (just run `make`). Incremental builds for dev
build *GO_BUILD_ARGS:
  #!/usr/bin/env bash
  go mod tidy
  go build ./... {{GO_BUILD_ARGS}}

# Run all tests
test *GO_TEST_ARGS:
  #!/usr/bin/env bash
  go test ./... {{GO_TEST_ARGS}}

# Lint and format all files
format:
  #!/usr/bin/env bash
  go fmt ./...

# Delete build output and crud that that piles up in the workspace
clean:
  #!/usr/bin/env bash
  rm -rf .build/

# Run clean first and then delete dependencies and lockfile(s)
purge: clean
  #!/usr/bin/env bash
  rm go.sum
