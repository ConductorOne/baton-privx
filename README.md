![Baton Logo](./docs/images/baton-logo.png)

# `baton-privx` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-privx.svg)](https://pkg.go.dev/github.com/conductorone/baton-privx) ![main ci](https://github.com/conductorone/baton-privx/actions/workflows/main.yaml/badge.svg)

`baton-privx` is a connector for PrivX built using the [Baton SDK](https://github.com/conductorone/baton-sdk). 
It communicates with the privx API to sync data about groups, and users.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-privx
baton-privx
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-privx:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-privx/cmd/baton-privx@main

BATON_API_CLIENT_ID=00000000-0000-0000-0000-000000000000 BATON_API_CLIENT_SECRET=Z2ltbWUgY2hlZXNlYnVyZ2Vy BATON_BASE_URL=https://privx.example.com/ BATON_OAUTH_CLIENT_ID=privx-external BATON_OAUTH_CLIENT_SECRET=aGV5IHdhc3N1cA baton-privx

baton resources
```

# Data Model

`baton-privx` will pull down information about the following PrivX resources:
- Roles
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually 
building spreadsheets. We welcome contributions, and ideas, no matter how 
small&mdash;our goal is to make identity and permissions sprawl less painful for 
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-privx` Command Line Usage

```
baton-privx

Usage:
  baton-privx [flags]
  baton-privx [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-client-id string         The API Client ID (a UUID.) ($BATON_API_CLIENT_ID)
      --api-client-secret string     The API Client Secret (a base64 string.) ($BATON_API_CLIENT_SECRET)
      --base-url string              The hostname (URL) for your PrivX instance ($BATON_BASE_URL)
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-privx
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --oauth-client-id string       The OAuth Client ID (e.g. "privx-external".) ($BATON_OAUTH_CLIENT_ID)
      --oauth-client-secret string   The OAuth Client Secret (a base64 string.) ($BATON_OAUTH_CLIENT_SECRET)
  -p, --provisioning                 This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-privx

Use "baton-privx [command] --help" for more information about a command.
```
