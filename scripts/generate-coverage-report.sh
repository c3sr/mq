#!/usr/bin/env sh

# This is a gross hack to install gocovmerge without adding it to go.mod.
# go 1.15 doesn't provide a more convenient way to do this, that I can find.
cd /tmp && \
  go get -t -v github.com/wadey/gocovmerge/... && \
  cd -

go clean -testcache && go test -coverprofile cover.out ./...
scripts/run-integration-tests.sh

gocovmerge cover.out integration.out > merged.out
go tool cover -html=merged.out -o coverage.html
rm cover.out integration.out merged.out