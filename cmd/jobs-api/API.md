FORMAT: 1A

# Swagger Petstore

This is a sample server Petstore server.  You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For this sample, you can use the api key `special-key` to test the authorization filters.

## Group pet

Everything about your Pets

### /v2/pet/{petId}/uploadImage

#### uploads an image [POST]

+ Parameters

  + petId (required) - ID of pet to update

+ Request (application/x-www-form-urlencoded)

  + Headers

          Accept: application/json

  + Attributes - Additional data to pass to server

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          {
            "type": "consequat nostrud cillum incidid",
            "message": "et Duis adipisicing",
            "code": 3586027
          }

  + Schema

          {
            "type": "object",
            "properties": {
              "code": {
                "type": "integer",
                "format": "int32"
              },
              "type": {
                "type": "string"
              },
              "message": {
                "type": "string"
              }
            }
          }

### /v2/pet

#### Add a new pet to the store [POST]

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "name": "nostrud consequat sint est velit",
            "photoUrls": [
              "qui mollit veniam",
              "quis pariatur",
              "laborum pariatur do tempor commodo",
              "est laboris aliqua veniam"
            ]
          }

  + Schema

          {
            "type": "object",
            "required": [
              "name",
              "photoUrls"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "category": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer",
                    "format": "int64"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "xml": {
                  "name": "Category"
                }
              },
              "name": {
                "type": "string",
                "example": "doggie"
              },
              "photoUrls": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "type": "string",
                  "xml": {
                    "name": "photoUrl"
                  }
                }
              },
              "tags": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "xml": {
                    "name": "tag"
                  },
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  }
                }
              },
              "status": {
                "type": "string",
                "description": "pet status in the store",
                "enum": [
                  "available",
                  "pending",
                  "sold"
                ]
              }
            }
          }

+ Response 405 (application/json)

  Invalid input

  + Body

#### Update an existing pet [PUT]

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "name": "laborum nisi adipisicing eu",
            "photoUrls": [
              "pariatur amet dolore deserunt occaecat",
              "ex in amet",
              "dolore id mollit proident aliqua"
            ]
          }

  + Schema

          {
            "type": "object",
            "required": [
              "name",
              "photoUrls"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "category": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer",
                    "format": "int64"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "xml": {
                  "name": "Category"
                }
              },
              "name": {
                "type": "string",
                "example": "doggie"
              },
              "photoUrls": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "type": "string",
                  "xml": {
                    "name": "photoUrl"
                  }
                }
              },
              "tags": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "xml": {
                    "name": "tag"
                  },
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  }
                }
              },
              "status": {
                "type": "string",
                "description": "pet status in the store",
                "enum": [
                  "available",
                  "pending",
                  "sold"
                ]
              }
            }
          }

+ Response 400 (application/json)

  Invalid ID supplied

  + Body

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "name": "ad sunt irure non dolor",
            "photoUrls": [
              "occaecat ea nostrud dolo",
              "minim et dolore pariatur"
            ]
          }

  + Schema

          {
            "type": "object",
            "required": [
              "name",
              "photoUrls"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "category": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer",
                    "format": "int64"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "xml": {
                  "name": "Category"
                }
              },
              "name": {
                "type": "string",
                "example": "doggie"
              },
              "photoUrls": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "type": "string",
                  "xml": {
                    "name": "photoUrl"
                  }
                }
              },
              "tags": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "xml": {
                    "name": "tag"
                  },
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  }
                }
              },
              "status": {
                "type": "string",
                "description": "pet status in the store",
                "enum": [
                  "available",
                  "pending",
                  "sold"
                ]
              }
            }
          }

+ Response 404 (application/json)

  Pet not found

  + Body

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "name": "culpa Duis",
            "photoUrls": [
              "ad irure commodo voluptate enim",
              "qui id"
            ],
            "status": "available"
          }

  + Schema

          {
            "type": "object",
            "required": [
              "name",
              "photoUrls"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "category": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer",
                    "format": "int64"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "xml": {
                  "name": "Category"
                }
              },
              "name": {
                "type": "string",
                "example": "doggie"
              },
              "photoUrls": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "type": "string",
                  "xml": {
                    "name": "photoUrl"
                  }
                }
              },
              "tags": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "xml": {
                    "name": "tag"
                  },
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  }
                }
              },
              "status": {
                "type": "string",
                "description": "pet status in the store",
                "enum": [
                  "available",
                  "pending",
                  "sold"
                ]
              }
            }
          }

+ Response 405 (application/json)

  Validation exception

  + Body

#### Finds Pets by status [GET /v2/pet/findByStatus]

Multiple status values can be provided with comma separated strings

+ Parameters

  + status: available,pending,sold (required) - Status values that need to be considered for filter

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          [
            {
              "name": "pariatur et fugiat mollit sit",
              "photoUrls": [
                "dolor v",
                "ut",
                "ea"
              ]
            },
            {
              "name": "amet dolore laboris",
              "photoUrls": [
                "anim velit irure ea",
                "fugiat",
                "ut"
              ],
              "tags": [
                {
                  "id": 85444612,
                  "name": "enim labore dolor dolor"
                },
                {
                  "id": -72117236
                },
                {
                  "id": 28975743,
                  "name": "f"
                }
              ]
            },
            {
              "name": "proident",
              "photoUrls": [
                "cupidatat do pariatur tempor elit"
              ]
            },
            {
              "name": "eu",
              "photoUrls": [
                "sunt mollit pariatur sed aute"
              ]
            }
          ]

  + Schema

          {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "name",
                "photoUrls"
              ],
              "properties": {
                "id": {
                  "type": "integer",
                  "format": "int64"
                },
                "category": {
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  },
                  "xml": {
                    "name": "Category"
                  }
                },
                "name": {
                  "type": "string",
                  "example": "doggie"
                },
                "photoUrls": {
                  "type": "array",
                  "xml": {
                    "wrapped": true
                  },
                  "items": {
                    "type": "string",
                    "xml": {
                      "name": "photoUrl"
                    }
                  }
                },
                "tags": {
                  "type": "array",
                  "xml": {
                    "wrapped": true
                  },
                  "items": {
                    "xml": {
                      "name": "tag"
                    },
                    "type": "object",
                    "properties": {
                      "id": {
                        "type": "integer",
                        "format": "int64"
                      },
                      "name": {
                        "type": "string"
                      }
                    }
                  }
                },
                "status": {
                  "type": "string",
                  "description": "pet status in the store",
                  "enum": [
                    "available",
                    "pending",
                    "sold"
                  ]
                }
              },
              "xml": {
                "name": "Pet"
              }
            }
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid status value

  + Body

#### Finds Pets by tags [GET /v2/pet/findByTags]

Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.

+ Parameters

  + tags:  (required) - Tags to filter by

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          [
            {
              "name": "exercitation",
              "photoUrls": [
                "est dol",
                "magna aute eu"
              ]
            },
            {
              "name": "veniam ipsum anim",
              "photoUrls": [
                "ex qui et minim"
              ]
            }
          ]

  + Schema

          {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "name",
                "photoUrls"
              ],
              "properties": {
                "id": {
                  "type": "integer",
                  "format": "int64"
                },
                "category": {
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  },
                  "xml": {
                    "name": "Category"
                  }
                },
                "name": {
                  "type": "string",
                  "example": "doggie"
                },
                "photoUrls": {
                  "type": "array",
                  "xml": {
                    "wrapped": true
                  },
                  "items": {
                    "type": "string",
                    "xml": {
                      "name": "photoUrl"
                    }
                  }
                },
                "tags": {
                  "type": "array",
                  "xml": {
                    "wrapped": true
                  },
                  "items": {
                    "xml": {
                      "name": "tag"
                    },
                    "type": "object",
                    "properties": {
                      "id": {
                        "type": "integer",
                        "format": "int64"
                      },
                      "name": {
                        "type": "string"
                      }
                    }
                  }
                },
                "status": {
                  "type": "string",
                  "description": "pet status in the store",
                  "enum": [
                    "available",
                    "pending",
                    "sold"
                  ]
                }
              },
              "xml": {
                "name": "Pet"
              }
            }
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid tag value

  + Body

### /v2/pet/{petId}

#### Find pet by ID [GET]

Returns a single pet

+ Parameters

  + petId (required) - ID of pet to return

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          {
            "name": "deserunt fugiat Duis ea",
            "photoUrls": [
              "elit minim ut consequat",
              "exercitation laborum mollit"
            ],
            "status": "sold",
            "category": {
              "name": "do tempor magna",
              "id": -29166945
            }
          }

  + Schema

          {
            "type": "object",
            "required": [
              "name",
              "photoUrls"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "category": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer",
                    "format": "int64"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "xml": {
                  "name": "Category"
                }
              },
              "name": {
                "type": "string",
                "example": "doggie"
              },
              "photoUrls": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "type": "string",
                  "xml": {
                    "name": "photoUrl"
                  }
                }
              },
              "tags": {
                "type": "array",
                "xml": {
                  "wrapped": true
                },
                "items": {
                  "xml": {
                    "name": "tag"
                  },
                  "type": "object",
                  "properties": {
                    "id": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "name": {
                      "type": "string"
                    }
                  }
                }
              },
              "status": {
                "type": "string",
                "description": "pet status in the store",
                "enum": [
                  "available",
                  "pending",
                  "sold"
                ]
              }
            }
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid ID supplied

  + Body

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 404 (application/json)

  Pet not found

  + Body

#### Updates a pet in the store with form data [POST]

+ Parameters

  + petId (required) - ID of pet that needs to be updated

+ Request (application/x-www-form-urlencoded)

  + Headers

          Accept: application/json

  + Attributes - Updated name of the pet

  + Body

+ Response 405 (application/json)

  Invalid input

  + Body

#### Deletes a pet [DELETE]

+ Parameters

  + petId (required) - Pet id to delete

+ Request

  + Headers

          Accept: application/json
          api_key: 

  + Body

+ Response 400 (application/json)

  Invalid ID supplied

  + Body

+ Request

  + Headers

          Accept: application/json
          api_key: 

  + Body

+ Response 404 (application/json)

  Pet not found

  + Body

## Group store

Access to Petstore orders

### /v2/store/order

#### Place an order for a pet [POST]

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "shipDate": "3005-09-11T01:00:30.700Z"
          }

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "petId": {
                "type": "integer",
                "format": "int64"
              },
              "quantity": {
                "type": "integer",
                "format": "int32"
              },
              "shipDate": {
                "type": "string",
                "format": "date-time"
              },
              "status": {
                "type": "string",
                "description": "Order Status",
                "enum": [
                  "placed",
                  "approved",
                  "delivered"
                ]
              },
              "complete": {
                "type": "boolean"
              }
            }
          }

+ Response 200 (application/json)

  successful operation

  + Body

          {
            "complete": false,
            "petId": 26815841
          }

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "petId": {
                "type": "integer",
                "format": "int64"
              },
              "quantity": {
                "type": "integer",
                "format": "int32"
              },
              "shipDate": {
                "type": "string",
                "format": "date-time"
              },
              "status": {
                "type": "string",
                "description": "Order Status",
                "enum": [
                  "placed",
                  "approved",
                  "delivered"
                ]
              },
              "complete": {
                "type": "boolean"
              }
            }
          }

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {
            "shipDate": "3255-05-06T18:01:43.313Z"
          }

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "petId": {
                "type": "integer",
                "format": "int64"
              },
              "quantity": {
                "type": "integer",
                "format": "int32"
              },
              "shipDate": {
                "type": "string",
                "format": "date-time"
              },
              "status": {
                "type": "string",
                "description": "Order Status",
                "enum": [
                  "placed",
                  "approved",
                  "delivered"
                ]
              },
              "complete": {
                "type": "boolean"
              }
            }
          }

+ Response 400 (application/json)

  Invalid Order

  + Body

### /v2/store/order/{orderId}

#### Find purchase order by ID [GET]

For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions

+ Parameters

  + orderId (required) - ID of pet that needs to be fetched

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          {}

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "petId": {
                "type": "integer",
                "format": "int64"
              },
              "quantity": {
                "type": "integer",
                "format": "int32"
              },
              "shipDate": {
                "type": "string",
                "format": "date-time"
              },
              "status": {
                "type": "string",
                "description": "Order Status",
                "enum": [
                  "placed",
                  "approved",
                  "delivered"
                ]
              },
              "complete": {
                "type": "boolean"
              }
            }
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid ID supplied

  + Body

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 404 (application/json)

  Order not found

  + Body

#### Delete purchase order by ID [DELETE]

For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors

+ Parameters

  + orderId (required) - ID of the order that needs to be deleted

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid ID supplied

  + Body

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 404 (application/json)

  Order not found

  + Body

### /v2/store/inventory

#### Returns pet inventories by status [GET]

Returns a map of status codes to quantities

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          {}

  + Schema

          {
            "type": "object",
            "additionalProperties": {
              "type": "integer",
              "format": "int32"
            }
          }

## Group user

Operations about user

### /v2/user/createWithArray

#### Creates list of users with given input array [POST]

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          [
            {
              "password": "in Excepteur",
              "lastName": "cillum ut cupidatat"
            },
            {
              "password": "tempor",
              "email": "sunt",
              "userStatus": 7504521
            },
            {},
            {
              "email": "adipisicing",
              "firstName": "elit",
              "password": "cillum esse reprehenderit magna do"
            },
            {}
          ]

  + Schema

          {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "integer",
                  "format": "int64"
                },
                "username": {
                  "type": "string"
                },
                "firstName": {
                  "type": "string"
                },
                "lastName": {
                  "type": "string"
                },
                "email": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                },
                "phone": {
                  "type": "string"
                },
                "userStatus": {
                  "type": "integer",
                  "format": "int32",
                  "description": "User Status"
                }
              },
              "xml": {
                "name": "User"
              }
            }
          }

### /v2/user/createWithList

#### Creates list of users with given input array [POST]

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          [
            {},
            {},
            {},
            {}
          ]

  + Schema

          {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "integer",
                  "format": "int64"
                },
                "username": {
                  "type": "string"
                },
                "firstName": {
                  "type": "string"
                },
                "lastName": {
                  "type": "string"
                },
                "email": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                },
                "phone": {
                  "type": "string"
                },
                "userStatus": {
                  "type": "integer",
                  "format": "int32",
                  "description": "User Status"
                }
              },
              "xml": {
                "name": "User"
              }
            }
          }

### /v2/user/{username}

#### Get user by user name [GET]

+ Parameters

  + username (required) - The name that needs to be fetched. Use user1 for testing.

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

          {
            "userStatus": -92133763,
            "id": -31525395
          }

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "username": {
                "type": "string"
              },
              "firstName": {
                "type": "string"
              },
              "lastName": {
                "type": "string"
              },
              "email": {
                "type": "string"
              },
              "password": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              },
              "userStatus": {
                "type": "integer",
                "format": "int32",
                "description": "User Status"
              }
            }
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid username supplied

  + Body

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 404 (application/json)

  User not found

  + Body

#### Updated user [PUT]

This can only be done by the logged in user.

+ Parameters

  + username (required) - name that need to be updated

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {}

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "username": {
                "type": "string"
              },
              "firstName": {
                "type": "string"
              },
              "lastName": {
                "type": "string"
              },
              "email": {
                "type": "string"
              },
              "password": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              },
              "userStatus": {
                "type": "integer",
                "format": "int32",
                "description": "User Status"
              }
            }
          }

+ Response 400 (application/json)

  Invalid user supplied

  + Body

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {}

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "username": {
                "type": "string"
              },
              "firstName": {
                "type": "string"
              },
              "lastName": {
                "type": "string"
              },
              "email": {
                "type": "string"
              },
              "password": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              },
              "userStatus": {
                "type": "integer",
                "format": "int32",
                "description": "User Status"
              }
            }
          }

+ Response 404 (application/json)

  User not found

  + Body

#### Delete user [DELETE]

This can only be done by the logged in user.

+ Parameters

  + username (required) - The name that needs to be deleted

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid username supplied

  + Body

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 404 (application/json)

  User not found

  + Body

#### Logs user into the system [GET /v2/user/login]

+ Parameters

  + username (required) - The user name for login

  + password (required) - The password for login in clear text

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Headers

          X-Expires-After: 
          X-Rate-Limit: 

  + Body

          qui ut reprehenderit

  + Schema

          {
            "type": "string"
          }

+ Request

  + Headers

          Accept: application/json

  + Body

+ Response 400 (application/json)

  Invalid username/password supplied

  + Body

### /v2/user/logout

#### Logs out current logged in user session [GET]

### /v2/user

#### Create user [POST]

This can only be done by the logged in user.

+ Request (application/json)

  + Headers

          Accept: application/json

  + Body

          {}

  + Schema

          {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64"
              },
              "username": {
                "type": "string"
              },
              "firstName": {
                "type": "string"
              },
              "lastName": {
                "type": "string"
              },
              "email": {
                "type": "string"
              },
              "password": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              },
              "userStatus": {
                "type": "integer",
                "format": "int32",
                "description": "User Status"
              }
            }
          }

