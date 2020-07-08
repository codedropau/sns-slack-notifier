package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliSlackToken        = kingpin.Flag("token", "Slack token").Envar("SLACK_TOKEN").String()
	cliSlackChannel      = kingpin.Flag("channel-id", "Slack channel ID").Envar("SLACK_CHANNEL_ID").String()
	cliSlackMessageColor = kingpin.Flag("color", "Slack mesage color").Envar("SLACK_MESSAGE_COLOR").Default(ColourGood).String()
)

const (
	// ColourGood is the slack colour for good.
	ColourGood = "good"

	// ColourDanger is the slack colour for danger.
	ColourDanger = "danger"

	// ColourWarning is the slack colour for warning.
	ColourWarning = "warning"
)

func main() {
	kingpin.Parse()
	lambda.Start(HandleRequest)
}

// HandleRequest contains the code which will be executed.
func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	slackAPI := slack.New(*cliSlackToken)
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		attachment := slack.Attachment{
			Color: *cliSlackMessageColor,
			Text: snsRecord.Message,
			Footer: fmt.Sprintf(":skpr: %s Source: %s MessageID: %s Topic: %s", snsRecord.Timestamp.Format(time.UnixDate), record.EventSource, snsRecord.MessageID, snsRecord.TopicArn),
		}
		_, _, err := slackAPI.PostMessage(*cliSlackChannel, slack.MsgOptionText(snsRecord.Subject, false), slack.MsgOptionAttachments(attachment))
		if err != nil {
			return err
		}
		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
	}

	return nil
}
