package njir

import (
	"net/http"
	"os"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	StoryPointField = "customfield_10033"
)

func NewClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_EMAIL"),
		Password: os.Getenv("JIRA_API_TOKEN"),
	}

	// Attach OTEL transport to the underlying transport used by BasicAuth
	tp.Transport = otelhttp.NewTransport(
		http.DefaultTransport,
	)

	httpClient := tp.Client()

	client, err := jira.NewClient(httpClient, os.Getenv("JIRA_BASE_URL"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetIssueStoryPoint(fields *jira.IssueFields) int {
	sp := 0
	v := fields.Unknowns[StoryPointField]
	if v != nil {
		switch x := v.(type) {
		case int:
			sp = x
		case float64:
			sp = int(x)
		case string:
			if f, err := strconv.ParseFloat(x, 64); err == nil {
				sp = int(f)
			}
		}
	}
	return sp
}
