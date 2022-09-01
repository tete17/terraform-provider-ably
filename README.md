# terraform-provider
Ably Terraform Provider

# [Ably](https://www.ably.com)

_[Ably](https://ably.com) is the platform that powers synchronized digital experiences in realtime. Whether attending an event in a virtual venue, receiving realtime financial information, or monitoring live car performance data – consumers simply expect realtime digital experiences as standard. Ably provides a suite of APIs to build, extend, and deliver powerful digital experiences in realtime for more than 250 million devices across 80 countries each month. Organizations like Bloomberg, HubSpot, Verizon, and Hopin depend on Ably’s platform to offload the growing complexity of business-critical realtime data synchronization at global scale. For more information, see the [Ably documentation](https://ably.com/docs)._

[![npm version](https://img.shields.io/npm/v/ably.svg?style=flat)](https://img.shields.io/npm/v/ably.svg?style=flat)

This is a Terraform provider for Ably.

## Supported platforms

This provider supports the following systems/architectures:

- Darwin / AMD64
- Darwin / ARMv8
- Linux / AMD64 (this is required for usage in Terraform Cloud, see below)
- Linux / ARMv8 (sometimes referred to as AArch64 or ARM64)
- Linux / ARMv6
- Windows / AMD64

## Installation

To install Ably Terraform provider:

1. Create a Control API token by logging into your Ably account and going to https://ably.com/users/access_tokens (Account -> My Access Tokens). This token should have permissions for the operations that you plan to do with with Terraform provider. More details are available in [Ably documentation](https://ably.com/docs/control-api#authentication).
2. Set environment varable `ABLY_ACCOUNT_TOKEN` to the token you have created.
4. (Optional) If you need to use a custom endpoint for the Control API, set the `ABLY_URL` environment variable as well.
3. Add the following to your Terraform configuration file

```terraform
terraform {
  required_providers {
    ably = {
      source  = "ably/ably"
      version = "0.0.0"
    }
  }
}

provider "ably" {
}
```

4. (Optional) Alternatively you can also specify Control API token and custom endpoint in the provider configuration directly: 

```terraform
provider "ably" {
  token = <Control API token>
  url   = <custom endpoint>
}
```


## Using Ably Terraform provider

This readme gives a basic example; for more examples see the [examples/](examples/) folder, rendered documentation on the Terraform Registry, or [docs/] folder. 


```terraform
# Define Ably app
resource "ably_app" "app0" {
  name                      = "ably-tf-provider-app-0000"
  status                    = "enabled"
  tls_only                  = true
}

# Add a key
resource "ably_api_key" "api_key_0" {
  app_id = ably_app.app0.id
  name   = "key-0000"
  capabilities = {
    "channel2" = ["publish"],
    "channel3" = ["subscribe"],
    "channel33" = ["subscribe"],
  }
}

# Configure a queue
resource "ably_queue" "example_queue" {
  app_id     = ably_app.app0.id
  name       = "queue_name"
  ttl        = 60
  max_length = 10000
  region     = "us-east-1-a"
}
```

## Dependencies

This provider uses [Ably Control API](https://ably.com/docs/api/control-api) and [Ably Control Go SDK](https://github.com/ably/ably-control-go) under the hood. 


## Support, feedback and troubleshooting

Please visit http://support.ably.com/ for access to our knowledgebase and to ask for any assistance.

You can also view the [community reported Github issues](https://github.com/ably/terraform-provider-ably/issues).

To see what has changed in recent versions, see the [CHANGELOG](CHANGELOG.md).

## Contributing

For guidance on how to contribute to this project, see the [CONTRIBUTING.md](CONTRIBUTING.md).
