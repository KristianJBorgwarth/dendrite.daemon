## dendrite.daemon

The core component of [Dendrite](https://kristianjborgwarth.github.io/dendrite.docs/). A Go server that manages the index for a local markdown vault and exposes it over JSON-RPC 2.0.

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/KristianJBorgwarth/dendrite.daemon/master/install.sh | sh
```

Supports Linux and macOS on amd64 and arm64. The script places the binary in `~/.local/bin` on Linux and `/usr/local/bin` on macOS. Re-run to update.

## Documentation

See the [Dendrite docs](https://kristianjborgwarth.github.io/dendrite.docs/) for the full protocol reference, method documentation, and setup guide.
