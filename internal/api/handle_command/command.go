package handle_command

import (
	"context"
	"fmt"

	"github.com/AnilRedshift/captions_please_go/internal/api/common"
	"github.com/AnilRedshift/captions_please_go/internal/api/replier"
	"github.com/sirupsen/logrus"
)

func Command(ctx context.Context, message string, job common.ActivityJob) common.ActivityResult {
	command := parseCommand(message)
	logrus.Debug(fmt.Sprintf("%s: Directive %s for language %v", job.Tweet.Id, command.directive, command.tag))
	ctx = replier.WithLanguage(ctx, command.tag)

	switch command.directive {
	case autoDirective:
		return HandleAuto(ctx, job.Tweet)
	case altTextDirective:
		return HandleAltText(ctx, job.Tweet)
	case ocrDirective:
		return HandleOCR(ctx, job.Tweet)
	case describeDirective:
		return HandleDescribe(ctx, job.Tweet)
	case helpDirective:
		fallthrough
	default:
		return Help(ctx, job.Tweet)
	}
}