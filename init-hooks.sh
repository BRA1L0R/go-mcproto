#!/bin/bash

echo "
echo \"Running pre-commit hook\"

#golangcgi-lint
golangcgi-lint run ./...

#unit testing
go test -v ./...
" > .git/hooks/pre-commit

chmod +x .git/hooks/pre-commit