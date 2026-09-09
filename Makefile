SHELL := /bin/bash

GO_SUBMODULES := $(shell find . \
	\( -path './.git' -o -path '*/vendor' \) -prune -o \
	-type f -name 'go.mod' -print \
	| sed -e 's#^\./##' -e 's#/go.mod$$##' \
	| grep -v '^go.mod$$' \
	| sort)

.PHONY: modules
modules:
	@echo "Detected Go submodules:"
	@for module in $(GO_SUBMODULES); do \
		echo "  - $$module"; \
	done

.PHONY: release
release:
	@set -euo pipefail; \
	\
	BUMP="$(filter patch minor major,$(MAKECMDGOALS))"; \
	\
	if [ -z "$$BUMP" ]; then \
		echo "Usage: make release [patch|minor|major]"; \
		exit 1; \
	fi; \
	\
	if [ "$$(echo "$$BUMP" | wc -w | tr -d ' ')" -ne 1 ]; then \
		echo "Specify exactly one release type: patch, minor, or major"; \
		exit 1; \
	fi; \
	\
	git fetch --tags; \
	\
	CURRENT_VERSION="$$(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-version:refname | head -n 1)"; \
	if [ -z "$$CURRENT_VERSION" ]; then \
		CURRENT_VERSION="v0.0.0"; \
	fi; \
	\
	VERSION="$${CURRENT_VERSION#v}"; \
	IFS='.' read -r MAJOR MINOR PATCH <<< "$$VERSION"; \
	\
	case "$$BUMP" in \
		major) \
			MAJOR=$$((MAJOR + 1)); \
			MINOR=0; \
			PATCH=0; \
			;; \
		minor) \
			MINOR=$$((MINOR + 1)); \
			PATCH=0; \
			;; \
		patch) \
			PATCH=$$((PATCH + 1)); \
			;; \
	esac; \
	\
	NEXT_VERSION="v$${MAJOR}.$${MINOR}.$${PATCH}"; \
	MESSAGE="Release $$NEXT_VERSION"; \
	BRANCH="$$(git branch --show-current)"; \
	\
	if [ -z "$$BRANCH" ]; then \
		echo "Cannot release from a detached HEAD"; \
		exit 1; \
	fi; \
	\
	TAGS=("$$NEXT_VERSION"); \
	for module in $(GO_SUBMODULES); do \
		TAGS+=("$${module}/$${NEXT_VERSION}"); \
	done; \
	\
	for tag in "$${TAGS[@]}"; do \
		if git show-ref --verify --quiet "refs/tags/$$tag"; then \
			echo "Tag $$tag already exists"; \
			exit 1; \
		fi; \
	done; \
	\
	echo "Releasing $$CURRENT_VERSION -> $$NEXT_VERSION"; \
	echo ""; \
	echo "Tags:"; \
	printf '  %s\n' "$${TAGS[@]}"; \
	echo ""; \
	\
	git add -A; \
	\
	if ! git diff --cached --quiet; then \
		git commit -m "$$MESSAGE"; \
	else \
		echo "No changes to commit; tagging current HEAD"; \
	fi; \
	\
	for tag in "$${TAGS[@]}"; do \
		git tag -a "$$tag" -m "$$MESSAGE"; \
	done; \
	\
	git push --atomic origin "$$BRANCH" "$${TAGS[@]}"; \
	\
	echo ""; \
	echo "Released $$NEXT_VERSION"

.PHONY: patch
patch:
	@:

.PHONY: minor
minor:
	@:

.PHONY: major
major:
	@:
