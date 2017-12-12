// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON, FlatSwaggerJSON json.RawMessage

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
    "title": "Account Private API",
    "contact": {
      "name": "mars"
    },
    "version": "v1"
  },
  "basePath": "/api-private/v1/accounts",
  "paths": {
    "/login": {
      "post": {
        "operationId": "Login",
        "parameters": [
          {
            "type": "string",
            "name": "name",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "password",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      }
    },
    "/logout": {
      "post": {
        "operationId": "Logout",
        "parameters": [
          {
            "type": "string",
            "name": "jwt",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      }
    },
    "/smsCode": {
      "post": {
        "operationId": "SmsCode",
        "parameters": [
          {
            "type": "string",
            "name": "scene",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "captchaId",
            "in": "query"
          },
          {
            "type": "string",
            "name": "captchaCode",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "ok"
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      }
    },
    "/smsLogin": {
      "post": {
        "summary": "sms login",
        "operationId": "SmsLogin",
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "smsCode",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/smsSignup": {
      "post": {
        "summary": "sms signup",
        "operationId": "SmsSignup",
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "smsCode",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "password",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "$ref": "#/responses/ErrorResponse"
          }
        }
      }
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
            "description": "Errors",
            "type": "array",
            "items": {
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
              }
            }
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
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
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
    "title": "Account Private API",
    "contact": {
      "name": "mars"
    },
    "version": "v1"
  },
  "basePath": "/api-private/v1/accounts",
  "paths": {
    "/login": {
      "post": {
        "operationId": "Login",
        "parameters": [
          {
            "type": "string",
            "name": "name",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "password",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/loginDefaultBody"
            }
          }
        }
      }
    },
    "/logout": {
      "post": {
        "operationId": "Logout",
        "parameters": [
          {
            "type": "string",
            "name": "jwt",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok"
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/logoutDefaultBody"
            }
          }
        }
      }
    },
    "/smsCode": {
      "post": {
        "operationId": "SmsCode",
        "parameters": [
          {
            "type": "string",
            "name": "scene",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "captchaId",
            "in": "query"
          },
          {
            "type": "string",
            "name": "captchaCode",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "ok"
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/smsCodeDefaultBody"
            }
          }
        }
      }
    },
    "/smsLogin": {
      "post": {
        "summary": "sms login",
        "operationId": "SmsLogin",
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "smsCode",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/smsSignup": {
      "post": {
        "summary": "sms signup",
        "operationId": "SmsSignup",
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "smsCode",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "password",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/smsSignupDefaultBody"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "loginDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/loginDefaultBodyErrors"
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
    "loginDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/loginDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "loginDefaultBodyErrorsItems": {
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
    },
    "logoutDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/logoutDefaultBodyErrors"
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
    "logoutDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/logoutDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "logoutDefaultBodyErrorsItems": {
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
    },
    "smsCodeDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/smsCodeDefaultBodyErrors"
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
    "smsCodeDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/smsCodeDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "smsCodeDefaultBodyErrorsItems": {
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
    },
    "smsSignupDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/smsSignupDefaultBodyErrors"
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
    "smsSignupDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/smsSignupDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "smsSignupDefaultBodyErrorsItems": {
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
            "description": "Errors",
            "type": "array",
            "items": {
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
              }
            }
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
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}`))
}