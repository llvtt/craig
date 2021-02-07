package slack

import (
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

func copyHeaders(req *events.APIGatewayProxyRequest) (header http.Header) {
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

func HttpRequest(req *events.APIGatewayProxyRequest) (request *http.Request, err error) {
	var verifier slack.SecretsVerifier
	headers := copyHeaders(req)

	verifier, err = slack.NewSecretsVerifier(headers, slackSigningSecret)
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
