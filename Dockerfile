FROM golang:1.15.6-alpine3.12

LABEL "com.github.actions.icon"="send"
LABEL "com.github.actions.color"="purple"
LABEL "com.github.actions.name"="Slack Notify Approved"
LABEL "com.github.actions.description"="This action will send notification to Slack When pull request is approved."
LABEL version="1.0.0"

WORKDIR ${GOPATH}/src/github.com/ayaka-tajiri/slack-notify-when-approved
COPY main.go ${GOPATH}/src/github.com/ayaka-tajiri/slack-notify-when-approved

RUN go get -v ./...
RUN go build -o /go/bin/slack-notify-when-papproved .

RUN apk update \
	&& apk upgrade \
	&& apk add \
	jq \
	bash

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
