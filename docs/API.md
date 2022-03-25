FORMAT: 1A

# Jobs API

[Jobs API](https://github.com/NFT-com/indexer) is a REST API serving the data related to discovery and parsing of NFT data.

## Group discovery

Discovery jobs

### /discoveries

#### Create a new discovery job [POST]

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
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "addresses": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 200 (application/json)

  Discovery job created

  + Body

            {}

  + Schema

            {
              "type": "object"
            }

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "addresses": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 400 (application/json)

  Invalid job

  + Body

#### Find Discoveries jobs [GET /discoveries{?status}]

+ Parameters

  + status: created,queued,processing,failed,finished,canceled - Job Status, empty for all status

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

            [
              {},
              {
                "addresses": [
                  "ullamco",
                  "pariatur in laborum",
                  "ipsum magna adipisicing exercitation commodo"
                ]
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
                    "type": "string",
                    "format": "uuid"
                  },
                  "chain_url": {
                    "type": "string",
                    "format": "url"
                  },
                  "chain_type": {
                    "type": "string",
                    "example": "web3"
                  },
                  "block_number": {
                    "type": "string",
                    "example": "12345"
                  },
                  "addresses": {
                    "type": "array",
                    "items": {
                      "type": "string"
                    }
                  },
                  "standard_type": {
                    "type": "string",
                    "example": "ERC721"
                  },
                  "status": {
                    "type": "string",
                    "description": "Discovery Status",
                    "default": "created",
                    "enum": [
                      "created",
                      "queued",
                      "processing",
                      "failed",
                      "finished",
                      "canceled"
                    ]
                  }
                }
              }
            }

### /discoveries/{id}

#### Find discovery job by ID [GET]

Returns a single discovery job

+ Parameters

  + id (required) - ID of the discovery job

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  Successful operation

  + Body

            {}

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "addresses": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 404 (application/json)

  Discovery job not found

  + Body

#### Updates a discovery job [PATCH]

+ Parameters

  + id (required) - ID of the discovery job that needs to be updated

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {
              "status": "created"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 200 (application/json)

  Successful operation

  + Body

            {
              "chain_type": "consectetur eiusmod incididunt",
              "block_number": "dolore incididunt"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "addresses": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {
              "status": "processing"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 400 (application/json)

  Invalid discovery job status

  + Body

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {
              "status": "canceled"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 404 (application/json)

  Discovery job not found

  + Body

### /discoveries/{id}/requeue

#### Recreates a discovery job [POST]

+ Parameters

  + id (required) - ID of discovery job to requeue

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  Successful operation

  + Body

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 404 (application/json)

  Discovery job not found

  + Body

## Group parse

Parsing jobs

### /parsings

#### Create a new parsing job [POST]

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {
              "block_number": "ex",
              "event_type": "sint",
              "chain_type": "proident Lorem quis dolore"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "address": {
                  "type": "string",
                  "example": "0x2685d224956b311c8729f1ad72c9cacd9f6e8f56"
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "event_type": {
                  "type": "string",
                  "example": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 200 (application/json)

  Parsing job created

  + Body

            {}

  + Schema

            {
              "type": "object"
            }

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "address": {
                  "type": "string",
                  "example": "0x2685d224956b311c8729f1ad72c9cacd9f6e8f56"
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "event_type": {
                  "type": "string",
                  "example": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 400 (application/json)

  Invalid parsing job

  + Body

#### Find parsing jobs [GET /parsings{?status}]

+ Parameters

  + status: created,queued,processing,failed,finished,canceled - Job Status, empty for all status

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  Successful operation

  + Body

            [
              {
                "status": "failed"
              },
              {
                "status": "finished"
              },
              {
                "status": "created"
              },
              {
                "status": "failed",
                "address": "ut dolore anim in"
              },
              {
                "standard_type": "ex dolore",
                "address": "proident exercitation aliqua ipsum in"
              }
            ]

  + Schema

            {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "string",
                    "format": "uuid"
                  },
                  "chain_url": {
                    "type": "string",
                    "format": "url"
                  },
                  "chain_type": {
                    "type": "string",
                    "example": "web3"
                  },
                  "block_number": {
                    "type": "string",
                    "example": "12345"
                  },
                  "address": {
                    "type": "string",
                    "example": "0x2685d224956b311c8729f1ad72c9cacd9f6e8f56"
                  },
                  "standard_type": {
                    "type": "string",
                    "example": "ERC721"
                  },
                  "event_type": {
                    "type": "string",
                    "example": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
                  },
                  "status": {
                    "type": "string",
                    "description": "Discovery Status",
                    "default": "created",
                    "enum": [
                      "created",
                      "queued",
                      "processing",
                      "failed",
                      "finished",
                      "canceled"
                    ]
                  }
                }
              }
            }

### /parsings/{id}

#### Find parsing job by ID [GET]

Returns a single parsing job

+ Parameters

  + id (required) - ID of the parsing job

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  Successful operation

  + Body

            {}

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "address": {
                  "type": "string",
                  "example": "0x2685d224956b311c8729f1ad72c9cacd9f6e8f56"
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "event_type": {
                  "type": "string",
                  "example": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 404 (application/json)

  Parsing job not found

  + Body

#### Updates a parsing job [PATCH]

+ Parameters

  + id (required) - ID of the parsing job that needs to be updated

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {}

  + Schema

            {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 200 (application/json)

  Successful operation

  + Body

            {
              "standard_type": "anim culpa consequat ut minim",
              "address": "commodo ex cillum"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string",
                  "format": "uuid"
                },
                "chain_url": {
                  "type": "string",
                  "format": "url"
                },
                "chain_type": {
                  "type": "string",
                  "example": "web3"
                },
                "block_number": {
                  "type": "string",
                  "example": "12345"
                },
                "address": {
                  "type": "string",
                  "example": "0x2685d224956b311c8729f1ad72c9cacd9f6e8f56"
                },
                "standard_type": {
                  "type": "string",
                  "example": "ERC721"
                },
                "event_type": {
                  "type": "string",
                  "example": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
                },
                "status": {
                  "type": "string",
                  "description": "Discovery Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Request (application/json)

  + Headers

            Accept: application/json

  + Body

            {
              "status": "failed"
            }

  + Schema

            {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 400 (application/json)

  Invalid status

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
                "status": {
                  "type": "string",
                  "description": "Job Status",
                  "default": "created",
                  "enum": [
                    "created",
                    "queued",
                    "processing",
                    "failed",
                    "finished",
                    "canceled"
                  ]
                }
              }
            }

+ Response 404 (application/json)

  Parsing job not found

  + Body

### /parsings/{id}/requeue

#### Recreates a parsing job [POST]

+ Parameters

  + id (required) - ID of parsing job to requeue

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 200 (application/json)

  successful operation

  + Body

+ Request

  + Headers

            Accept: application/json

  + Body

+ Response 404 (application/json)

  Parsing job not found

  + Body

