package hooks

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/hooks/formatting"
)

type Shooter interface {
	Fire(ctx context.Context, r *v1.Reaction) error
}

type MockShooter struct{}

func (s *MockShooter) Fire(ctx context.Context, r *v1.Reaction) error {
	context.GetLogger(ctx).Warnf("fired event for container %q and hook %q", r.Container.Name, r.Hook.Name)
	return nil
}

type LiveShooter struct {
	HttpClient *http.Client
}

func (s *LiveShooter) Fire(ctx context.Context, r *v1.Reaction) error {
	var body []byte
	var err error
	var bodyType string

	switch r.Hook.Format {
	case v1.FormatJSON:
		body, bodyType, err = formatting.JSON(r)
	case v1.FormatSlackJSON:
		body, bodyType, err = formatting.Slack(r)
	default:
		return fmt.Errorf("body format %q unsupported", r.Hook.Format)
	}

	if err != nil {
		return fmt.Errorf("error formatting body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, r.Hook.Url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("Content-Length", fmt.Sprint(len(body)))
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("couldn't execute request: %v", err)
	}

	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return fmt.Errorf("unexpected response status for hook shot: %d", resp.StatusCode)
	}

	return nil
}
