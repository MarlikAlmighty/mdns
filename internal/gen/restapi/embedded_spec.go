// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Rest API for mDNS Server",
    "title": "mDNS API",
    "contact": {
      "email": "cryptocoin62@gmail.com"
    },
    "version": "1.0.0"
  },
  "host": "localhost",
  "basePath": "/",
  "paths": {
    "/dns": {
      "get": {
        "tags": [
          "list"
        ],
        "summary": "Show all dns records",
        "operationId": "show_dns_records",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_records"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "put": {
        "consumes": [
          "application/json; charset=utf-8"
        ],
        "produces": [
          "application/json; charset=utf-8"
        ],
        "tags": [
          "update"
        ],
        "summary": "Update dns entry",
        "operationId": "update_dns_entry",
        "parameters": [
          {
            "name": "update",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json; charset=utf-8"
        ],
        "produces": [
          "application/json; charset=utf-8"
        ],
        "tags": [
          "add"
        ],
        "summary": "Add dns entry",
        "operationId": "add_dns_entry",
        "parameters": [
          {
            "name": "add",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "delete"
        ],
        "summary": "Delete dns entry",
        "operationId": "delete_dns_entry",
        "parameters": [
          {
            "name": "delete",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      }
    },
    "/dns/{id}": {
      "get": {
        "tags": [
          "show"
        ],
        "summary": "List one dns entry",
        "operationId": "list_one__dns_entry",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "answer": {
      "type": "object",
      "properties": {
        "Code": {
          "type": "integer",
          "format": "uint32"
        },
        "Message": {
          "type": "string"
        }
      }
    },
    "dns_entry": {
      "type": "object",
      "properties": {
        "acme": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "dkim": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "domain": {
          "type": "string"
        },
        "ips": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "ipv4": {
          "type": "string"
        },
        "ipv6": {
          "type": "string"
        },
        "private_key": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "public_key": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dns_records": {
      "type": "object",
      "additionalProperties": {
        "$ref": "#/definitions/dns_entry"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Rest API for mDNS Server",
    "title": "mDNS API",
    "contact": {
      "email": "cryptocoin62@gmail.com"
    },
    "version": "1.0.0"
  },
  "host": "localhost",
  "basePath": "/",
  "paths": {
    "/dns": {
      "get": {
        "tags": [
          "list"
        ],
        "summary": "Show all dns records",
        "operationId": "show_dns_records",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_records"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "put": {
        "consumes": [
          "application/json; charset=utf-8"
        ],
        "produces": [
          "application/json; charset=utf-8"
        ],
        "tags": [
          "update"
        ],
        "summary": "Update dns entry",
        "operationId": "update_dns_entry",
        "parameters": [
          {
            "name": "update",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json; charset=utf-8"
        ],
        "produces": [
          "application/json; charset=utf-8"
        ],
        "tags": [
          "add"
        ],
        "summary": "Add dns entry",
        "operationId": "add_dns_entry",
        "parameters": [
          {
            "name": "add",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "delete"
        ],
        "summary": "Delete dns entry",
        "operationId": "delete_dns_entry",
        "parameters": [
          {
            "name": "delete",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      }
    },
    "/dns/{id}": {
      "get": {
        "tags": [
          "show"
        ],
        "summary": "List one dns entry",
        "operationId": "list_one__dns_entry",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/dns_entry"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/answer"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "answer": {
      "type": "object",
      "properties": {
        "Code": {
          "type": "integer",
          "format": "uint32"
        },
        "Message": {
          "type": "string"
        }
      }
    },
    "dns_entry": {
      "type": "object",
      "properties": {
        "acme": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "dkim": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "domain": {
          "type": "string"
        },
        "ips": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "ipv4": {
          "type": "string"
        },
        "ipv6": {
          "type": "string"
        },
        "private_key": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "public_key": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dns_records": {
      "type": "object",
      "additionalProperties": {
        "$ref": "#/definitions/dns_entry"
      }
    }
  }
}`))
}
