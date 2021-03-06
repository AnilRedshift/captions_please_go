package handle_command

import (
	"context"
	"fmt"

	"github.com/AnilRedshift/captions_please_go/internal/api/common"
	"github.com/AnilRedshift/captions_please_go/pkg/message"
	"github.com/AnilRedshift/captions_please_go/pkg/twitter"
	"github.com/sirupsen/logrus"
)

func Help(ctx context.Context, tweet *twitter.Tweet) common.ActivityResult {
	result := _reply(ctx, tweet, message.HelpMessage(ctx))
	if result.Err != nil {
		logrus.Info(fmt.Sprintf("%s: Replying with the help message failed with %v", tweet.Id, result.Err))
	}
	return common.ActivityResult{Tweet: tweet, Action: "reply with help"}
}

func Unknown(ctx context.Context, tweet *twitter.Tweet) common.ActivityResult {
	result := _reply(ctx, tweet, message.UnknownCommandMessage(ctx))
	if result.Err != nil {
		logrus.Info(fmt.Sprintf("%s: Replying with the unknown message failed with %v", tweet.Id, result.Err))
	}
	return common.ActivityResult{Tweet: tweet, Action: "reply with unknown message"}

}
