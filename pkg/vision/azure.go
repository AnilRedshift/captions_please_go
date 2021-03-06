package vision

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnilRedshift/captions_please_go/pkg/message"
	"github.com/AnilRedshift/captions_please_go/pkg/structured_error"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.1/computervision"
	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

type azure struct {
	client        computervision.BaseClient
	supportedTags []language.Tag
}

func NewAzureVision(computerVisionKey string) Describer {
	client := computervision.New("https://captionspleasecomputervision.cognitiveservices.azure.com")
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)
	supportedTags := make([]language.Tag, len(languageMapping))
	i := 0
	for tag := range languageMapping {
		supportedTags[i] = tag
		i++
	}
	return &azure{client: client, supportedTags: supportedTags}
}

var languageMapping = map[language.Tag]string{
	language.English:           "en",
	language.Spanish:           "es",
	language.Japanese:          "ja",
	language.Portuguese:        "pt",
	language.SimplifiedChinese: "zh",
}

func (a *azure) Describe(ctx context.Context, url string) ([]VisionResult, structured_error.StructuredError) {
	var result []VisionResult
	var err error
	imageURL := computervision.ImageURL{URL: &url}
	tag, wrongLangErr := message.GetCompatibleLanguage(ctx, a.supportedTags)
	if wrongLangErr != nil {
		logrus.Debug("Azure cannot produce descriptions in the desired language")
		tag = language.English
	}
	var description computervision.ImageDescription
	description, err = a.client.DescribeImage(ctx, imageURL, nil, languageMapping[tag], nil)
	logDebugJSON(description)
	if err == nil && description.Captions != nil {
		result = make([]VisionResult, 0, len(*description.Captions))
		for i, caption := range *description.Captions {
			if caption.Confidence != nil && caption.Text != nil {
				result = result[:len(result)+1]
				result[i] = VisionResult{Text: *caption.Text, Confidence: float32(*caption.Confidence)}
			}
		}
		logDebugJSON(result)
	} else {
		logrus.Debug(fmt.Sprintf("azure describe returned error %v", err))
	}

	if err == nil {
		err = wrongLangErr
	}
	return result, structured_error.Wrap(err, structured_error.DescribeError)
}

func (a *azure) GetOCR(ctx context.Context, url string) (*OCRResult, structured_error.StructuredError) {
	var ocr *OCRResult
	imageURL := computervision.ImageURL{URL: &url}
	result, err := a.client.RecognizePrintedText(ctx, true, imageURL, computervision.OcrLanguagesUnk)
	builder := strings.Builder{}
	if err == nil && result.Regions != nil {
		for _, region := range *result.Regions {
			for _, line := range *region.Lines {
				for _, word := range *line.Words {
					builder.WriteString(*word.Text + " ")
				}
				builder.WriteString(" ")
			}
			builder.WriteString("\n\n")
		}

		ocrLanguage := OCRLanguage{Tag: language.English, Confidence: 0.0}
		if result.Language != nil {
			tag, parseErr := language.Parse(*result.Language)
			if parseErr == nil {
				ocrLanguage = OCRLanguage{Tag: tag, Confidence: 1.0}
			}
		}
		ocr = &OCRResult{
			Text:     builder.String(),
			Language: ocrLanguage,
		}
	}
	return ocr, structured_error.Wrap(err, structured_error.OCRError)
}
