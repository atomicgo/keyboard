#!/usr/bin/env bash

# ! ┌───────────────────────────────────────────────────────────────────┐
# ! │                                                                   │
# ! │                          IMPORTANT NOTE                           │
# ! │                                                                   │
# ! │   This file is synced with https://github.com/atomicgo/template   │
# ! │                                                                   │
# ! │        Please apply all changes to the template repository        │
# ! │                                                                   │
# ! └───────────────────────────────────────────────────────────────────┘

set -euo pipefail

# Colors
bold="\033[1m"
red="\033[31m"
green="\033[32m"
yellow="\033[33m"
blue="\033[34m"
cyan="\033[36m"
reset="\033[0m"

usage() {
	printf "%bUsage:%b %s --major|--minor|--patch\n" "$bold" "$reset" "$0"
}

release_type=""

case "${1:-}" in
	--major)
		release_type="major"
		;;
	--minor)
		release_type="minor"
		;;
	--patch)
		release_type="patch"
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		usage
		exit 1
		;;
esac

if [[ $# -ne 1 ]]; then
	usage
	exit 1
fi

if [[ "$release_type" == "major" ]]; then
	printf "\n%b%s%b\n" "$red$bold" "============================================================" "$reset"
	printf "%b%s%b\n" "$red$bold" "  BREAKING RELEASE WARNING" "$reset"
	printf "%b%s%b\n" "$yellow" "  This major release will be marked as a breaking release." "$reset"
	printf "%b%s%b\n\n" "$red$bold" "============================================================" "$reset"
fi

# Repository context
repo_root="$(git rev-parse --show-toplevel 2>/dev/null)" || {
	printf "%bError:%b this script must be run inside a git repository.\n" "$red$bold" "$reset" >&2
	exit 1
}

cd "$repo_root"

remote_url="$(git remote get-url origin 2>/dev/null || true)"
repo_name="$(basename "$repo_root")"

if [[ -z "$remote_url" ]]; then
	printf "%bError:%b an origin remote is required to push release tags.\n" "$red$bold" "$reset" >&2
	exit 1
fi

remote_path="${remote_url%.git}"
remote_path="${remote_path#git@*:}"
remote_path="${remote_path#*://*/}"
repo_name="$(basename "$remote_path")"

printf "%bFetching tags...%b\n" "$blue" "$reset"
git fetch --tags --quiet

# Find the latest stable vMAJOR.MINOR.PATCH tag.
latest_tag=""
while IFS= read -r tag; do
	if [[ "$tag" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
		latest_tag="$tag"
		break
	fi
done < <(git tag --list "v*" --sort=-v:refname)

if [[ -z "$latest_tag" ]]; then
	latest_tag="v0.0.0"
fi

if [[ "$latest_tag" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
	major="${BASH_REMATCH[1]}"
	minor="${BASH_REMATCH[2]}"
	patch="${BASH_REMATCH[3]}"
else
	printf "%bError:%b latest tag %b%s%b is not a valid vMAJOR.MINOR.PATCH tag.\n" "$red$bold" "$reset" "$yellow" "$latest_tag" "$reset" >&2
	exit 1
fi

case "$release_type" in
	major)
		major=$((major + 1))
		;;
	minor)
		minor=$((minor + 1))
		;;
	patch)
		patch=$((patch + 1))
		;;
esac

new_tag="v${major}.${minor}.${patch}"

if git rev-parse -q --verify "refs/tags/$new_tag" >/dev/null; then
	printf "%bError:%b tag %b%s%b already exists.\n" "$red$bold" "$reset" "$yellow" "$new_tag" "$reset" >&2
	exit 1
fi

printf "\n%bCurrent latest tag:%b %b%s%b\n" "$bold" "$reset" "$yellow" "$latest_tag" "$reset"
printf "%bNew %s tag:%b %b%s%b\n\n" "$bold" "$release_type" "$reset" "$green" "$new_tag" "$reset"

printf "%bRelease %b%s%b %b%s%b -> %b%s%b? [y/N]: %b" "$bold" "$cyan" "$repo_name" "$reset$bold" "$yellow" "$latest_tag" "$reset$bold" "$green" "$new_tag" "$reset" "$reset"
read -r answer

case "$answer" in
	y | Y | yes | YES)
		;;
	*)
		printf "%bRelease cancelled.%b\n" "$yellow" "$reset"
		exit 0
		;;
esac

printf "%bCreating tag %s...%b\n" "$blue" "$new_tag" "$reset"
git tag "$new_tag"

printf "%bPushing tag %s to origin...%b\n" "$blue" "$new_tag" "$reset"
git push origin "$new_tag"

printf "%bRelease tag %s pushed successfully.%b\n" "$green$bold" "$new_tag" "$reset"
