# slack-notify-when-approved

This Github Action sends a Slack notification when the number of Approve reaches a certain number.

## Usage
```
on: pull_request_review
name: Slack Notify approved pull requests
jobs:
  slackNotifyWhenApproved:
    name: Slack Notify When approved
    runs-on: ubuntu-latest
    steps:
    - name: Slack Notify When approved
      uses: ayaka-tajiri/slack-notify-when-approved@main
      env:
        APPROVALS: 0
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        SLACK_WEB_HOOK: ${{ secrets.SLACK_WEB_HOOK }}
```

## env
- APPROVALS
- GITHUB_TOKEN
- SLACK_WEB_HOOK
- SLACK_TITLE
- SLACK_MESSAGE
- SLACK_COLOR
- SLACK_USERNAME
- SLACK_FOOTER
- SLACK_ICON
- SLACK_EMOJI
- SLACK_CHANNEL

## LISENCE
