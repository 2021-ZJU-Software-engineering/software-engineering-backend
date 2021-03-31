// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/case": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "New a case",
                "parameters": [
                    {
                        "description": "patient ID, doctor ID, department name and other case details",
                        "name": "caseDetail",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.ReturnedData"
                        }
                    }
                }
            }
        },
        "/case/department": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get lastest case by department name",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "patient ID",
                        "name": "patientID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "department name",
                        "name": "department",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/case/{caseid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get case by a case-id list",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "description": "case IDs",
                        "name": "caseIDList",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Update case",
                "parameters": [
                    {
                        "description": "case ID and updated details",
                        "name": "caseDetail",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/case/{caseid}/prescription": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get prescriptions by case id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "case ID",
                        "name": "caseID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/pkg_cases.Prescription"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "New a prescrition",
                "parameters": [
                    {
                        "description": "case ID and prescription details",
                        "name": "prescriptionDetail",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/pkg_cases.Prescription"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/case/{caseid}/prescription/{prescriptionid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get prescription by prescription id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "prescription ID",
                        "name": "prescriptionID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/pkg_cases.Prescription"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/case/{currentcaseid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get prevoius cases by current case ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "current case ID",
                        "name": "caseID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.ReturnedData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "respond to a ping request from client",
                "produces": [
                    "application/json"
                ],
                "summary": "Test server up statue",
                "responses": {
                    "200": {
                        "description": "Good, server is up",
                        "schema": {
                            "$ref": "#/definitions/api.ReturnedData"
                        }
                    }
                }
            }
        },
        "/{patientid}/case": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get the lastest case ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "patient ID",
                        "name": "patientID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "/{patientid}/case/department": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Case"
                ],
                "summary": "Get a case ID list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "patient ID",
                        "name": "patientID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "department name",
                        "name": "department",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "allOf": [
                                    {
                                        "$ref": "#/definitions/api.ReturnedData"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            "data": {
                                                "type": "array",
                                                "items": {
                                                    "type": "integer"
                                                }
                                            }
                                        }
                                    }
                                ]
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ReturnedData": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Anything you want to pass to the frontend, but make it simple and necessary\nIf there's nothing to return, this field will be omitted",
                    "type": "object"
                },
                "status": {
                    "description": "A simple string indicating the status.\nIs it ok, or some error occurs? If so, what is the error?\nIt should be \"ok\" is everything goes fine",
                    "type": "string"
                }
            }
        },
        "github.com_AsterNighT_software-engineering-backend_pkg_cases.Case": {
            "type": "object",
            "properties": {
                "complaint": {
                    "description": "Use urls to locate pictures",
                    "type": "string"
                },
                "department": {
                    "type": "string"
                },
                "diagnosis": {
                    "type": "string"
                },
                "doctorID": {
                    "type": "integer"
                },
                "id": {
                    "description": "Every object should have ID",
                    "type": "integer"
                },
                "pastHistory": {
                    "type": "string"
                },
                "patientID": {
                    "description": "A has many relationship should be on this",
                    "type": "integer"
                },
                "prescriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Prescription"
                    }
                },
                "previousCase": {
                    "description": "Previous case (the lastest one). If there is none prevous case, set nil",
                    "$ref": "#/definitions/github.com_AsterNighT_software-engineering-backend_pkg_cases.Case"
                },
                "previousCaseID": {
                    "type": "integer"
                }
            }
        },
        "github.com_AsterNighT_software-engineering-backend_pkg_cases.Prescription": {
            "type": "object",
            "properties": {
                "caseID": {
                    "type": "integer"
                },
                "details": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "pkg_cases.Case": {
            "type": "object",
            "properties": {
                "complaint": {
                    "description": "Use urls to locate pictures",
                    "type": "string"
                },
                "department": {
                    "type": "string"
                },
                "diagnosis": {
                    "type": "string"
                },
                "doctorID": {
                    "type": "integer"
                },
                "id": {
                    "description": "Every object should have ID",
                    "type": "integer"
                },
                "pastHistory": {
                    "type": "string"
                },
                "patientID": {
                    "description": "A has many relationship should be on this",
                    "type": "integer"
                },
                "prescriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/pkg_cases.Prescription"
                    }
                },
                "previousCase": {
                    "description": "Previous case (the lastest one). If there is none prevous case, set nil",
                    "$ref": "#/definitions/pkg_cases.Case"
                },
                "previousCaseID": {
                    "type": "integer"
                }
            }
        },
        "pkg_cases.Prescription": {
            "type": "object",
            "properties": {
                "caseID": {
                    "type": "integer"
                },
                "details": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:12448",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Swagger Example API",
	Description: "This is a sample server Petstore server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
