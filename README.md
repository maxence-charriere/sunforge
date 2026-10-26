# Sunforge

Reusable Go packages for crafting applications with a little magic.

I maintain these packages for my own projects and publish them for convenient
reuse. You're welcome to use them, but I don't provide support or compatibility
guarantees. APIs may change as my projects evolve.

## Installation

Requires Go 1.27.2 or later.

```sh
go get github.com/maxence-charriere/sunforge
```

## Packages

| Package | Purpose |
| --- | --- |
| [`vi`](./pkg/vi) | UI components and styles for go-app applications. |
| [`conduit`](./pkg/conduit) | HTTP services and clients, authentication, rate limiting, and task scheduling. |
| [`portal`](./pkg/portal) | OAuth and API clients for Google, Twitch, YouTube, X, LinkedIn, Threads, and Patreon. |
| [`stash`](./pkg/stash) | Caching with configurable expiration. |

## Development

```sh
make build
make test
```

## License

Sunforge is available under the [MIT License](./LICENSE).
