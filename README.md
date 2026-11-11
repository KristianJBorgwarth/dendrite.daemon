# dendrite_v2

C# Native AOT port of the Dendrite daemon. It speaks the same newline-delimited JSON-RPC over stdio as the Go daemon, so `dendrite.nvim` works unchanged.

Stack: .NET 10 (Native AOT), Mediator (source generated), Dapper.AOT + Microsoft.Data.Sqlite, YamlDotNet (representation model, reflection-free), System.Text.Json source generation.

## Layout

Dependencies point inward: `Daemon → Infrastructure → Application → Domain`.

- `Dendrite.Domain`: the `Note` aggregate (owns its tags, links and custom frontmatter; changed only via `Reindex`), `INoteRepository`, frontmatter/link/slug parsing, `VaultConfig`. No I/O.
- `Dendrite.Application`: Mediator commands/queries and handlers, services, and the ports they depend on (`INoteQueries` for reads, `IUnitOfWork`, `IVaultFiles`, ...). Commands go through the aggregate; queries read straight from the index.
- `Dendrite.Infrastructure`: SQLite + Dapper.AOT (`NoteRepository` persists the whole aggregate, `NoteQueries` serves reads), migrations, unit of work, vault file access.
- `Dendrite.Daemon`: the `dendrite` executable. JSON-RPC server, DI wiring, and the Mediator source generator (it must live in the host to see the handlers).

## Build and test

```sh
dotnet test
dotnet publish src/Dendrite.Daemon -c Release -r linux-x64 -o out
```

## RPC methods

`vault/init`, `vault/create`, `vault/rebuild`, `vault/list`, `vault/remove`, `vault/use`, `vault/relocate`, `vault/open`, `note/create|save|delete|goto|backlinks|search_by_tag|search_by_cf`, `completion/tag|slug`, `diagnostics/links`.

### Multiple vaults

Call `vault/init` once per vault; the vault name is its identity. Every other method accepts an optional `vault` param. When omitted, the vault is inferred from the request's `path` (longest matching vault root), then falls back to the only registered vault. With several vaults and no hint the call fails with an ambiguity error. Re-initialising a vault at the same path is idempotent; the same name at another path is rejected.

A vault is identified by its name, not its folder. If a folder moves, `vault/relocate {vault, path}` re-points it and rewrites the stored note paths, keeping the index. `vault/list` flags vaults whose folder no longer exists with `missing: true`; they stay registered so they can be relocated. `vault/create` also writes `.dendrite/vault.json` inside the vault, and `vault/open {path}` registers any folder with that manifest (or re-points the vault of that name to it).
