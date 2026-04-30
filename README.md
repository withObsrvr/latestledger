# latestledger.com

A tiny Go + templ + HTMX site for checking the latest Stellar ledger tip on mainnet and testnet.

## Run locally

With Nix:

```bash
nix develop
export OBSRVR_API_KEY="your-api-key"
templ generate
go run ./cmd/latestledger
```

Without Nix:

```bash
export OBSRVR_API_KEY="your-api-key"
go run github.com/a-h/templ/cmd/templ generate
go run ./cmd/latestledger
```

Open http://localhost:8080.

## Environment variables

- `OBSRVR_API_KEY` — Obsrvr Gateway API key. Sent server-side as `Authorization: Api-Key ...`.
- `PORT` — optional HTTP port, defaults to `8080`.

## DigitalOcean App Platform

Set `OBSRVR_API_KEY` as an encrypted environment variable.

Example build command:

```bash
CGO_ENABLED=0 go run github.com/a-h/templ/cmd/templ generate && CGO_ENABLED=0 go build -o latestledger ./cmd/latestledger
```

Example run command:

```bash
./latestledger
```
