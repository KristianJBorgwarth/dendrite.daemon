# dendrite_v2

C# Native AOT port of the Dendrite daemon. It speaks the same newline-delimited JSON-RPC over stdio as the Go daemon, so `dendrite.nvim` works unchanged.

Stack: .NET 10 (Native AOT), Mediator (source generated), Dapper.AOT + Microsoft.Data.Sqlite, YamlDotNet (representation model, reflection-free), System.Text.Json source generation.

## Layout

Dependencies point inward: `Daemon → Infrastructure → Application → Domain`.

- `Dendrite.Domain`: models, frontmatter/link/slug parsing, `VaultConfig`. No I/O.
- `Dendrite.Application`: Mediator commands/queries and handlers, services, and the ports they depend on (`INoteRepository`, `IUnitOfWork`, `IVaultFiles`, ...).
- `Dendrite.Infrastructure`: SQLite + Dapper.AOT repositories, migrations, unit of work, vault file access.
- `Dendrite.Daemon`: the `dendrite` executable. JSON-RPC server, DI wiring, and the Mediator source generator (it must live in the host to see the handlers).

## Build and test

```sh
dotnet test
dotnet publish src/Dendrite.Daemon -c Release -r linux-x64 -o out
```

AOT publishing needs a C toolchain. Without clang, add `-p:CppCompilerAndLinker=gcc -p:LinkerFlavor=bfd`.
Ship `libe_sqlite3.*` next to the `dendrite` binary.

## RPC methods

`vault/init`, `vault/rebuild`, `note/create|save|delete|goto|backlinks|search_by_tag|search_by_cf`, `completion/tag|slug`, `diagnostics/links`.
