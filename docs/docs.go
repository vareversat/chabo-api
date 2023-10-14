// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "email": "dev@vareversat.fr"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/vareversat/chabo-api/blob/main/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/forecasts": {
            "get": {
                "description": "Fetch all existing forecasts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Forecasts"
                ],
                "summary": "Get all foracasts",
                "parameters": [
                    {
                        "type": "string",
                        "format": "date-time",
                        "description": "The date to filter from (RFC3339)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "format": "int",
                        "default": 10,
                        "description": "Set the limit of the queried results",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "format": "int",
                        "default": 0,
                        "description": "Set the offset of the queried results",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "boat",
                            "maintenance"
                        ],
                        "type": "string",
                        "description": "The closing reason",
                        "name": "reason",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The boat name of the event",
                        "name": "boat",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "leaving_bordeaux",
                            "entering_in_bordeaux"
                        ],
                        "type": "string",
                        "description": "The boat maneuver of the event",
                        "name": "maneuver",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "UTC",
                        "description": "Timezone to format the date related fields (TZ identifier)",
                        "name": "Timezone",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domains.ForecastsResponse"
                        }
                    },
                    "400": {
                        "description": "Some params are missing and/or not properly formatted fror the requests",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "An error occured on the server side",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/forecasts/refresh": {
            "post": {
                "description": "Get, format et populate database with the data from the OpenData API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Forecasts"
                ],
                "summary": "Refresh the forecasts with the ones from the OpenData API",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domains.Refresh"
                        }
                    },
                    "429": {
                        "description": "Too many attempt to refresh",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "An error occured on the server side",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/forecasts/{id}": {
            "get": {
                "description": "Fetch a forecast by his unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Forecasts"
                ],
                "summary": "Get a foracast",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The forecast ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "UTC",
                        "description": "Timezone to format the date related fields (TZ identifier)",
                        "name": "Timezone",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domains.ForecastResponse"
                        }
                    },
                    "400": {
                        "description": "Some params are missing and/or not properly formatted fror the requests",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "The ID does not match any forecast",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "An error occured on the server side",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/refresh/last": {
            "get": {
                "description": "Get the last trace of refresh action on POST /forecasts/refresh",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Refreshes"
                ],
                "summary": "Get the last refresh action",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domains.Refresh"
                        }
                    },
                    "404": {
                        "description": "No previous refresh action exists",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "An error occured on the server side",
                        "schema": {
                            "$ref": "#/definitions/domains.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/system/healthcheck": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Get the status of the API",
                "responses": {
                    "200": {
                        "description": "The api is healthy",
                        "schema": {
                            "$ref": "#/definitions/domains.SystemHealthNOK"
                        }
                    },
                    "503": {
                        "description": "The api is unhealthy",
                        "schema": {
                            "$ref": "#/definitions/domains.SystemHealthOK"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domains.Boat": {
            "type": "object",
            "properties": {
                "approximative_crossing_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-05-25T00:53:16.535668Z"
                },
                "maneuver": {
                    "$ref": "#/definitions/domains.BoatManeuver"
                },
                "name": {
                    "type": "string",
                    "example": "EUROPA 2"
                }
            }
        },
        "domains.BoatManeuver": {
            "type": "string",
            "enum": [
                "leaving_bordeaux",
                "entering_in_bordeaux"
            ],
            "x-enum-varnames": [
                "Leaving",
                "Entering"
            ]
        },
        "domains.ClosingReason": {
            "type": "string",
            "enum": [
                "boat",
                "maintenance"
            ],
            "x-enum-varnames": [
                "BoatReason",
                "Maintenance"
            ]
        },
        "domains.ClosingType": {
            "type": "string",
            "enum": [
                "two_way",
                "one_way"
            ],
            "x-enum-varnames": [
                "TwoWay",
                "OneWay"
            ]
        },
        "domains.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "error in params"
                }
            }
        },
        "domains.Forecast": {
            "type": "object",
            "properties": {
                "boats": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domains.Boat"
                    }
                },
                "circulation_closing_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-05-25T00:53:16.535668Z"
                },
                "circulation_reopening_date": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-05-25T00:53:16.535668Z"
                },
                "closing_duration_ns": {
                    "type": "integer",
                    "example": 4980000000000
                },
                "closing_reason": {
                    "$ref": "#/definitions/domains.ClosingReason"
                },
                "closing_type": {
                    "$ref": "#/definitions/domains.ClosingType"
                },
                "id": {
                    "type": "string",
                    "example": "63a6430fc07ff1d895c9555ef2ef6e41c1e3b1f5"
                }
            }
        },
        "domains.ForecastResponse": {
            "type": "object",
            "properties": {
                "forecast": {
                    "$ref": "#/definitions/domains.Forecast"
                },
                "timezone": {
                    "type": "string",
                    "example": "UTC"
                }
            }
        },
        "domains.ForecastsResponse": {
            "type": "object",
            "properties": {
                "_links": {
                    "type": "array",
                    "items": {}
                },
                "forecasts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domains.Forecast"
                    }
                },
                "hits": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "timezone": {
                    "type": "string",
                    "example": "UTC"
                }
            }
        },
        "domains.Refresh": {
            "type": "object",
            "properties": {
                "duration_ns": {
                    "type": "integer",
                    "example": 348872934
                },
                "item_count": {
                    "type": "integer",
                    "example": 10
                },
                "timestamp": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-05-25T00:53:16.535668Z"
                }
            }
        },
        "domains.SystemHealthNOK": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "system is not running properly"
                }
            }
        },
        "domains.SystemHealthOK": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "system is running properly"
                }
            }
        }
    },
    "externalDocs": {
        "description": "All data comes from the Bordeaux Open Data API",
        "url": "https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Chabo API - The Chaban-Delmas bridge API",
	Description:      "You can get every info you need about all the events of the Chaban-Delmas bridge in Bordeaux, France",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
