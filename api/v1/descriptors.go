package v1

import (
	"net/http"
	"regexp"

	"github.com/danielkrainas/csense/api/describe"
	"github.com/danielkrainas/csense/api/errcode"
)

var (
	IDRegex = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[1][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)

	versionHeader = describe.Parameter{
		Name:        "cSense-API-Version",
		Type:        "string",
		Description: "The build version of the cSense API server.",
		Format:      "<version>",
		Examples:    []string{"0.0.0-dev"},
	}

	hostHeader = describe.Parameter{
		Name:        "Host",
		Type:        "string",
		Description: "",
		Format:      "<hostname>",
		Examples:    []string{"api.csense.io"},
	}

	hookIDParameter = describe.Parameter{
		Name:        "hook_id",
		Type:        "string",
		Description: "Identifier for organization",
		Format:      IDRegex.String(),
		Required:    true,
	}

	jsonContentLengthHeader = describe.Parameter{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "Length of the JSON body.",
		Format:      "<length>",
	}

	zeroContentLengthHeader = describe.Parameter{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "The 'Content-Length' header must be zero and the body must be empty.",
		Format:      "0",
	}

	hookNotFoundResp = describe.Response{
		Name:        "Hook Unknown Error",
		StatusCode:  http.StatusNotFound,
		Description: "The hook is not known to the server.",
		Headers: []describe.Parameter{
			versionHeader,
			jsonContentLengthHeader,
		},
		Body: describe.Body{
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		},
		ErrorCodes: []errcode.ErrorCode{
			ErrorCodeHookUnknown,
		},
	}
)

var (
	errorsBody = `{
	"errors:" [
	    {
            "code": <error code>,
            "message": <error message>,
            "detail": ...
        },
        ...
    ]
}`

	hookBody = `{

}`

	hooksBody = `[
` + hookBody + `, ...
]`
)

var API = struct {
	Routes []describe.Route `json:"routes"`
}{
	Routes: routeDescriptors,
}

var routeDescriptors = []describe.Route{
	{
		Name:        RouteNameBase,
		Path:        "/v1",
		Entity:      "Base",
		Description: "Base V1 API route, can be used for lightweight health and version check.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Check that the server supports the cSense V1 API.",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "The API implements the V1 protocol and is accessible.",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.Response{
							{
								Description: "The API does not support the V1 protocol.",
								StatusCode:  http.StatusNotFound,
								Headers: []describe.Parameter{
									versionHeader,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameHooks,
		Path:        "/v1/hooks",
		Entity:      "[]Hook",
		Description: "Route to retrieve the list of active hooks and create new ones.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Get all hooks",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "All hooks returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      hooksBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "PUT",
				Description: "Create a hook",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							hookIDParameter,
						},

						Successes: []describe.Response{
							{
								Description: "Hook created",
								StatusCode:  http.StatusCreated,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      hookBody,
								},
							},
						},

						Failures: []describe.Response{},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameHook,
		Path:        "/v1/hooks/{hook_id:" + IDRegex.String() + "}",
		Entity:      "Hook",
		Description: "Route to remove, retrieve, and modify an existing hook.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Get a single hook",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							hookIDParameter,
						},

						Successes: []describe.Response{
							{
								Description: "The hook was returned successfully.",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      hookBody,
								},
							},
						},

						Failures: []describe.Response{
							hookNotFoundResp,
						},
					},
				},
			},
			{
				Method:      "POST",
				Description: "Modify an existing hook",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							hookIDParameter,
						},

						Successes: []describe.Response{
							{
								Description: "The hook was removed successfully.",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.Response{
							hookNotFoundResp,
						},
					},
				},
			},
			{
				Method:      "DELETE",
				Description: "Remove a hook.",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							hookIDParameter,
						},

						Successes: []describe.Response{
							{
								Description: "The hook was removed successfully.",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.Response{
							hookNotFoundResp,
						},
					},
				},
			},
		},
	},
}

var routeDescriptorsMap map[string]describe.Route

func init() {
	routeDescriptorsMap = make(map[string]describe.Route, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		routeDescriptorsMap[descriptor.Name] = descriptor
	}
}
