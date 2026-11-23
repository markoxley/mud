#!/usr/bin/env bash
#
# release.sh — commit pending changes, tag the next release and push to origin.
#
# Usage:
#   scripts/release.sh [patch|minor|major|X.Y.Z|vX.Y.Z] [options]
#
# Bump level defaults to patch and is derived from the latest vX.Y.Z tag.
# An explicit version may be given instead (a leading "v" is added if absent).
#
# Options:
#   -m, --message MSG   Commit message for pending changes
#                       (default: "chore(release): vX.Y.Z")
#   -y, --yes           Non-interactive: skip the confirmation prompt (for CI)
#       --dry-run       Print what would happen without changing anything
#       --skip-tests    Skip the gofmt/build/vet/test verification
#   -h, --help          Show this help
#
# The script will:
#   1. Stamp CHANGELOG.md (renames "## Unreleased" to the dated version section)
#   2. Verify: gofmt -l, go build, go vet, go test -count=1
#   3. Commit any pending changes (including the changelog stamp)
#   4. Create an annotated tag and push the branch + tag to origin
#   5. Create a GitHub release via gh, if the CLI is installed

set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

BUMP="patch"
MSG=""
EXPLICIT=""
ASSUME_YES=0
DRY_RUN=0
SKIP_TESTS=0
REMOTE="${REMOTE:-origin}"

while [[ $# -gt 0 ]]; do
	case "$1" in
	patch | minor | major) BUMP="$1" ;;
	v[0-9]* | [0-9]*) EXPLICIT="$1" ;;
	-m | --message)
		MSG="${2:?--message requires an argument}"
		shift
		;;
	-y | --yes) ASSUME_YES=1 ;;
	--dry-run) DRY_RUN=1 ;;
	--skip-tests) SKIP_TESTS=1 ;;
	-h | --help)
		sed -n '2,28p' "$0"
		exit 0
		;;
	*)
		echo "error: unknown argument: $1" >&2
		exit 2
		;;
	esac
	shift
done

step() { printf '\n==> %s\n' "$*"; }
run() {
	if ((DRY_RUN)); then
		printf '[dry-run] %s\n' "$*"
	else
		"$@"
	fi
}

# --- 1. Sanity checks -------------------------------------------------------

step "Checking repository state"
BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$BRANCH" == "HEAD" ]]; then
	echo "error: detached HEAD — check out a branch first" >&2
	exit 1
fi
echo "branch: $BRANCH   remote: $REMOTE"

# --- 2. Determine next version ----------------------------------------------

step "Determining next version"
# 'v[0-9]*' excludes malformed tags such as 'v.0.0.7'.
LATEST="$(git tag -l 'v[0-9]*' | sort -V | tail -1)"
if [[ -z "$LATEST" ]]; then
	LATEST="v0.0.0"
	echo "no existing vX.Y.Z tags — starting from $LATEST"
fi
echo "latest tag: $LATEST"

if [[ -n "$EXPLICIT" ]]; then
	VER="v${EXPLICIT#v}"
else
	IFS=. read -r MAJ MIN PAT <<<"${LATEST#v}"
	case "$BUMP" in
	major)
		MAJ=$((MAJ + 1))
		MIN=0
		PAT=0
		;;
	minor)
		MIN=$((MIN + 1))
		PAT=0
		;;
	patch) PAT=$((PAT + 1)) ;;
	esac
	VER="v${MAJ}.${MIN}.${PAT}"
fi
echo "next version: $VER"

if git rev-parse -q --verify "refs/tags/$VER" >/dev/null 2>&1; then
	echo "error: tag $VER already exists locally" >&2
	exit 1
fi
# GIT_TERMINAL_PROMPT=0 avoids hanging on a credential prompt for private remotes.
if GIT_TERMINAL_PROMPT=0 git ls-remote --tags "$REMOTE" "refs/tags/$VER" 2>/dev/null | grep -q .; then
	echo "error: tag $VER already exists on $REMOTE" >&2
	exit 1
fi

# --- 3. Stamp the changelog --------------------------------------------------

if grep -q '^## Unreleased' CHANGELOG.md 2>/dev/null; then
	step "Stamping CHANGELOG.md: Unreleased -> $VER ($(date +%F))"
	if ((DRY_RUN)); then
		echo "[dry-run] would insert a dated '$VER' section after '## Unreleased'"
	else
		sed -i "s|^## Unreleased.*|## Unreleased\n\n## $VER - $(date +%Y-%m-%d)|" CHANGELOG.md
		git add CHANGELOG.md
	fi
fi

# --- 4. Verify ----------------------------------------------------------------

if ((SKIP_TESTS)); then
	echo "verification skipped (--skip-tests)"
else
	step "Verifying: gofmt, build, vet, test"
	if ((DRY_RUN)); then
		echo "[dry-run] gofmt -l . && go build ./... && go vet ./... && go test -count=1 ./..."
	else
		UNFMT="$(gofmt -l .)"
		if [[ -n "$UNFMT" ]]; then
			echo "error: unformatted files — run gofmt -w:" >&2
			echo "$UNFMT" >&2
			exit 1
		fi
		go build ./...
		go vet ./...
		go test -count=1 ./...
	fi
fi

# --- 5. Commit pending changes -------------------------------------------------

if [[ -n "$(git status --porcelain)" ]]; then
	step "Committing pending changes"
	git status --short
	run git add -A
	run git commit -m "${MSG:-chore(release): $VER}"
else
	echo "working tree clean — nothing to commit"
fi

# --- 6. Confirm ------------------------------------------------------------------

step "Release summary"
echo "version : $VER"
echo "branch  : $BRANCH -> $REMOTE"
echo "latest  : $LATEST"
if ((!ASSUME_YES && !DRY_RUN)); then
	read -r -p "Tag $VER and push to $REMOTE? [y/N] " REPLY
	[[ "$REPLY" =~ ^[Yy]$ ]] || {
		echo "aborted — nothing was pushed"
		exit 0
	}
fi

# --- 7. Tag, push, release --------------------------------------------------------

step "Tagging $VER"
run git tag -a "$VER" -m "Release $VER"

step "Pushing $BRANCH and $VER to $REMOTE"
run git push "$REMOTE" "$BRANCH"
run git push "$REMOTE" "refs/tags/$VER"

if command -v gh >/dev/null 2>&1; then
	step "Creating GitHub release"
	run gh release create "$VER" --generate-notes --title "$VER"
else
	echo "note: 'gh' CLI not found — tag pushed; create the release on GitHub if needed"
fi

step "Done — released $VER"
REPO_URL="$(git remote get-url "$REMOTE" 2>/dev/null | sed -e 's/\.git$//' -e 's|^git@\([^:]*\):|https://\1/|')"
[[ -n "$REPO_URL" ]] && echo "$REPO_URL/releases/tag/$VER"
