package formatting

import (
	"encoding/json"
	"fmt"

	"github.com/danielkrainas/csense/api/v1"
)

func Slack(r *v1.Reaction) ([]byte, string, error) {
	m := &message{
		Attachments: []*attachment{
			{
				Fallback:   fmt.Sprintf("%s on %s", r.Container.Name, r.Host.Hostname),
				Pretext:    fmt.Sprintf("Container %s on %s for %q", r.Container.State, r.Host.Hostname, r.Hook.Name),
				MarkdownIn: []string{"pretext"},
				Color:      "#394D54",
				Title:      fmt.Sprintf("Container %s on %s for %q", r.Container.State, r.Host.Hostname, r.Hook.Name),
				Timestamp:  r.Timestamp,
				Fields: []*field{
					{
						Title: "Host",
						Value: r.Host.Hostname,
						Short: true,
					},
					{
						Title: "State",
						Value: string(r.Container.State),
						Short: true,
					},
					{
						Title: "Container",
						Value: r.Container.Name,
						Short: false,
					},
					{
						Title: "Image",
						Value: r.Container.ImageName,
						Short: len(r.Container.ImageName) > 20,
					},
				},
			},
		},
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, "", err
	}

	return b, "application/json", nil
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type attachment struct {
	Fallback   string   `json:"fallback"`
	Pretext    string   `json"pretext"`
	Color      string   `json:"color"`
	Title      string   `json:"title"`
	MarkdownIn []string `json:"mrkdwn_in"`
	Fields     []*field `json:"fields"`
	Timestamp  int64    `json:"ts"`
}

type message struct {
	Attachments []*attachment `json:"attachments"`
}

/*
{
    "attachments": [
        {
            "fallback": "Container 'value' *Running* on *dan-nixlaptop*",
			"pretext": "The hook `Watcher1` was triggered by a container",
			"mrkdwn_in": ["pretext"],
            "color": "#394D54",
            "title": "Container /docker/container-nasdfmsdfasdasdgasd3423413515",
            "fields": [
                {
                    "title": "Image",
                    "value": "reg/image:3.3.3",
                    "short": true
                },
				{
					"title": "Host",
					"value": "dan-nixlaptop",
					"short": true
				}
            ],
            "ts": 123456789
        }
    ]
}
*/
