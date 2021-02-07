package slack

import (
	"context"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var (
	slackSigningSecret string
)

func init() {
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

type Slacker struct {
	Client   *slack.Client
	Channel  string
	verifier *slack.SecretsVerifier
}

func NewSlacker(channel string) *Slacker {
	return &Slacker{slack.New(os.Getenv("SLACK_ACCESS_TOKEN")), channel, nil}
}

func (s *Slacker) ParseCommand(req *events.APIGatewayProxyRequest) (*slack.SlashCommand, error) {
	httpReq, err := s.toHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	cmd, err := slack.SlashCommandParse(httpReq)
	if s.verifier == nil {
		return nil, errors.New("no verifier")
	}
	if err := s.verifier.Ensure(); err != nil {
		return nil, err
	}

	return &cmd, nil
}

func (s *Slacker) PostMessage(ctx context.Context, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	_, _, err := s.Client.PostMessageContext(ctx, s.Channel, slack.MsgOptionText(message, false))
	return err
}

func (s *Slacker) copyHeaders(req *events.APIGatewayProxyRequest) (header http.Header) {
	header = make(http.Header)
	for name, values := range req.Headers {
		splitValues := strings.Split(values, ",")
		for _, value := range splitValues {
			header[name] = append(header[name], value)
		}
	}
	for name, values := range req.MultiValueHeaders {
		header[name] = values
	}
	return
}

func (s *Slacker) toHTTPRequest(req *events.APIGatewayProxyRequest) (request *http.Request, err error) {
	var verifier slack.SecretsVerifier
	headers := s.copyHeaders(req)

	verifier, err = slack.NewSecretsVerifier(headers, slackSigningSecret)
	s.verifier = &verifier
	if err != nil {
		return
	}

	request, err = http.NewRequest(
		req.HTTPMethod,
		req.Path,
		ioutil.NopCloser(io.TeeReader(strings.NewReader(req.Body), &verifier)),
	)
	if request != nil {
		request.Header = headers
	}

	return
}
