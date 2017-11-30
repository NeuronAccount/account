// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Accounts",
    "contact": {
      "name": "mars"
    },
    "version": "v1"
  },
  "basePath": "/api/v1/oauth2",
  "paths": {
    "/token": {
      "post": {
        "operationId": "OAuth2Token",
        "security": [
          {
            "Basic": []
          }
        ],
        "parameters": [
          {
            "type": "string",
            "name": "grant_type",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "code",
            "in": "query"
          },
          {
            "type": "string",
            "name": "response_type",
            "in": "query"
          },
          {
            "type": "string",
            "name": "redirect_uri",
            "in": "query"
          },
          {
            "type": "string",
            "name": "state",
            "in": "query"
          },
          {
            "type": "string",
            "name": "client_id",
            "in": "query"
          },
          {
            "type": "string",
            "name": "refresh_token",
            "in": "query"
          },
          {
            "type": "string",
            "name": "scope",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "$ref": "#/definitions/AccessToken"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/oAuth2TokenDefaultBody"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "AccessToken": {
      "type": "object",
      "properties": {
        "access_token": {
          "type": "string"
        },
        "expires_in": {
          "type": "integer",
          "format": "int32"
        },
        "refresh_token": {
          "type": "string"
        },
        "scope": {
          "type": "string"
        },
        "token_type": {
          "type": "string"
        }
      }
    },
    "oAuth2TokenDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/oAuth2TokenDefaultBodyErrors"
        },
        "message": {
          "description": "Error message",
          "type": "string"
        },
        "status": {
          "type": "string",
          "format": "int32",
          "default": "Http status"
        }
      },
      "x-go-gen-location": "operations"
    },
    "oAuth2TokenDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/oAuth2TokenDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "oAuth2TokenDefaultBodyErrorsItems": {
      "type": "object",
      "properties": {
        "code": {
          "description": "error code",
          "type": "string"
        },
        "field": {
          "description": "field name",
          "type": "string"
        },
        "message": {
          "description": "error message",
          "type": "string"
        }
      },
      "x-go-gen-location": "operations"
    }
  },
  "responses": {
    "ErrorResponse": {
      "description": "Error response",
      "schema": {
        "type": "object",
        "properties": {
          "code": {
            "description": "Error code",
            "type": "string"
          },
          "errors": {
            "$ref": "#/definitions/oAuth2TokenDefaultBodyErrors"
          },
          "message": {
            "description": "Error message",
            "type": "string"
          },
          "status": {
            "type": "string",
            "format": "int32",
            "default": "Http status"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "Basic": {
      "type": "basic"
    }
  }
}`))
}