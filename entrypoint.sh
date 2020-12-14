#!/usr/bin/env bash

echo "start entrypoint.sh"

if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "Set the GITHUB_TOKEN."
  exit 1
fi

if [[ -z "$GITHUB_REPOSITORY" ]]; then
  echo "Set the GITHUB_REPOSITORY."
  exit 1
fi

if [[ -z "$GITHUB_EVENT_PATH" ]]; then
  echo "Set the GITHUB_EVENT_PATH."
  exit 1
fi

if [[ -z "$SLACK_WEB_HOOK" ]]; then
  echo "Set the SLACK_WEB_HOOK"
  exit 1
fi

NUMBER=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
GITHUB_SHA=$(cat $GITHUB_EVENT_PATH | jq -r .pull_request.head.sha)

export APPROVALS="$APPROVALS"
export PR_NUMBER="$NUMBER"
export GITHUB_TOKEN="$GITHUB_TOKEN"
export SLACK_WEB_HOOK="$SLACK_WEB_HOOK"

slack-notify-when-approved "$@"
