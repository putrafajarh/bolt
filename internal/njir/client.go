package njir

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
)

var (
	StoryPointField = "customfield_10033"
)

type NjirClient struct {
	client *jira.Client
}

func NewNjirClient() (*NjirClient, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &NjirClient{
		client: client,
	}, nil
}

// SearchIssuesWithTransitions search issues with appended transitions
func (c *NjirClient) SearchIssuesWithTransitions(ctx context.Context, jql string) ([]jira.Issue, error) {
	issues, _, err := c.client.Issue.SearchV2JQLWithContext(ctx, jql, &jira.SearchOptionsV2{
		Expand:       "names",
		Fields:       []string{"assignee", "project", "reporter", "status", "creator", StoryPointField, "summary", "created", "updated", "sprint", "transitions"},
		MaxResults:   200,
		FieldsByKeys: true,
	})
	if err != nil {
		return nil, err
	}

	// Fetch transitions concurrently before transaction
	g, ctx := errgroup.WithContext(ctx)
	for i := range issues {
		g.Go(func() error {
			transitions, _, err := c.client.Issue.GetTransitionsWithContext(ctx, issues[i].ID)
			if err != nil {
				return err
			}
			issues[i].Transitions = transitions
			return nil
		})
	}
	g.Wait()

	return issues, err
}

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
	if fields == nil {
		return 0
	}

	sp := 0
	v := fields.Unknowns[StoryPointField]
	if v != nil {
		switch x := v.(type) {
		case int:
			sp = x
		case int64:
			sp = int(x)
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
