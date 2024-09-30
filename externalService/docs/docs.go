// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const Title = "External Service API Docs"

const Version = "1.0"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/info": {
            "get": {
                "description": "Get information about a song based on group name and song name.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "info"
                ],
                "summary": "Get song info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the music group",
                        "name": "group_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the song",
                        "name": "song_name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with song info",
                        "schema": {
                            "$ref": "#/definitions/info.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request error response",
                        "schema": {
                            "$ref": "#/definitions/info.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "info.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
