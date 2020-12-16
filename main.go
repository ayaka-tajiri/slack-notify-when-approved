package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	EnvApprovals        = "APPROVALS"
	EnvPrNumber         = "PR_NUMBER"
	EnvGithubRepository = "GITHUB_REPOSITORY"
	EnvGithubToken      = "GITHUB_TOKEN"
	EnvGithubSha        = "GITHUB_SHA"
	EnvGithubRef        = "GITHUB_REF"
	EnvGithubEventName  = "GITHUB_EVENT_NAME"
	EnvGithubWorkflow   = "GITHUB_WORKFLOW"
	EnvGithubActor      = "GITHUB_ACTOR"
	EnvSlackWebHook     = "SLACK_WEB_HOOK"
	EnvSlackTitle       = "SLACK_TITLE"
	EnvSlackMessage     = "SLACK_MESSAGE"
	EnvSlackColor       = "SLACK_COLOR"
	EnvSlackUserName    = "SLACK_USERNAME"
	EnvSlackFooter      = "SLACK_FOOTER"
	EnvSlackIcon        = "SLACK_ICON"
	EnvSlackEmoji       = "SLACK_EMOJI"
	EnvSlackChannel     = "SLACK_CHANNEL"
)

type Reviewer struct {
	State string `json:"state"`
}

type SlackMessage struct {
	Text        string       `json:"text"`
	UserName    string       `json:"username"`
	IconUrl     string       `json:"icon_url"`
	IconEmoji   string       `json:"icon_emoji"`
	Channel     string       `json:"channel"`
	UnfurlLinks string       `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Color      string  `json:"color"`
	Pretext    string  `json:"pretext"`
	AuthorName string  `json:"author_name"`
	AuthorLink string  `json:"author_link"`
	AuthorIcon string  `json:"author_icon"`
	Footer     string  `json:"footer"`
	Fields     []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func main() {
	// check approvals
	strTargetApprovals := os.Getenv(EnvApprovals)
	targetApproval, err := strconv.Atoi(strTargetApprovals)
	if err != nil {
		log.Fatal(err)
	}
	if targetApproval != approvalCount() {
		os.Exit(0)
	}

	// send slack notify
	slackMessage := &SlackMessage{}
	slackMessage.slackNotify()
}

func approvalCount() int {
	githubRepository := os.Getenv(EnvGithubRepository)
	githubToken := os.Getenv(EnvGithubToken)
	prNumber := os.Getenv(EnvPrNumber)
	reviewerUri := "https://api.github.com"
	reviewerUrl := reviewerUri + "/repos/" + githubRepository + "/pulls/" + prNumber + "/reviews?per_page=100"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", reviewerUrl, nil)
	req.Header.Add("Authorization", "token "+githubToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var reviewers []Reviewer
	if err := json.Unmarshal(body, &reviewers); err != nil {
		log.Fatal(err)
	}

	totalApprovals := 0
	for _, reviewer := range reviewers {
		if reviewer.State == "APPROVED" {
			totalApprovals++
		}
	}

	return totalApprovals
}

func (slackMessage *SlackMessage) slackNotify() {
	longSha := os.Getenv(EnvGithubSha)
	commitSha := longSha[0:6]
	fields := []Field{
		{
			Title: "PullRequest URL",
			Value: "<https://github.com/" + os.Getenv(EnvGithubRepository) + "/pull/" + os.Getenv(EnvPrNumber)>",
			Short: true/,
		},
		{
			Title: "Actions URL",
			Value: "<https://github.com/" + os.Getenv(EnvGithubRepository) + "/commit/" + os.Getenv(EnvGithubSha) +
				"/checks|" + os.Getenv(EnvGithubWorkflow) + ">",
			Short: true,
		},
		{
			Title: "Commit",
			Value: "<https://github.com/" + os.Getenv(EnvGithubRepository) + "/commit/" + os.Getenv(EnvGithubSha) +
				"|" + commitSha + ">",
			Short: true,
		},
		{
			Title: getEnv(EnvSlackTitle, "message"),
			Value: getEnv(EnvSlackMessage, "Approved "+os.Getenv(EnvApprovals)+" user"),
			Short: false,
		},
	}

	slackMessage.UserName = os.Getenv(EnvSlackUserName)
	slackMessage.IconUrl = os.Getenv(EnvSlackIcon)
	slackMessage.IconEmoji = os.Getenv(EnvSlackEmoji)
	slackMessage.Channel = os.Getenv(EnvSlackChannel)
	slackMessage.Attachments = []Attachment{
		{
			Color:      getEnv(EnvSlackColor, "good"),
			Footer:     getEnv(EnvSlackFooter, "<https://github.com/ayaka-tajiri/slack-notify-when-approved|Powered By Ayaka's GitHub Actions Library>"),
			Fields:     fields,
		},
	}

	slackMessageByte, _ := json.Marshal(slackMessage)
	slackWebHook := os.Getenv(EnvSlackWebHook)
	_, err := http.Post(slackWebHook, "application/json", bytes.NewBuffer(slackMessageByte))
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
