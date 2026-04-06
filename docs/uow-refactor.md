# Unit of Work — Refactor Notes

## Current Issues

### 1. Lifetime
`UnitOfWork` is created once in the handler constructor and reused across every `Handle` call.
It should be **per-request**, created inside `Handle`.

### 2. Transaction leaks out
`Begin()` returns `*sql.Tx`, which is then passed manually to each repo constructor.
The UoW should own that wiring — callers should never touch the transaction directly.

### 3. Eager instantiation (rejected solution)
Having UoW create all repos in `Begin()` solves the leaking transaction, but forces every handler
to pay for repos it doesn't need. Not acceptable.

---

## Solution: Lazy Initialization

UoW exposes accessor methods that initialize each repo on first access, wired to the internal transaction.
Unused repos are never instantiated.

```go
type UnitOfWork struct {
    FileStore *store.FileStore
    tx        *sql.Tx
    tags      ITagRepository
    notes     INoteRepository
}

func NewUnitOfWork() *UnitOfWork {
    return &UnitOfWork{FileStore: store.NewFileStore()}
}

func (u *UnitOfWork) Begin() error {
    tx, err := persistence.GetDBContext().DB.Begin()
    if err != nil {
        return err
    }
    u.tx = tx
    return nil
}

func (u *UnitOfWork) Tags() ITagRepository {
    if u.tags == nil {
        u.tags = NewTagRepository(u.tx)
    }
    return u.tags
}

func (u *UnitOfWork) Notes() INoteRepository {
    if u.notes == nil {
        u.notes = NewNoteRepository(u.tx)
    }
    return u.notes
}

func (u *UnitOfWork) Commit() error {
    if err := u.FileStore.Flush(); err != nil {
        u.tx.Rollback()
        u.tx = nil
        return err
    }
    if err := u.tx.Commit(); err != nil {
        u.FileStore.Rollback()
        u.tx = nil
        return err
    }
    u.tx = nil
    return nil
}

func (u *UnitOfWork) Rollback() {
    if u.tx == nil {
        return
    }
    u.tx.Rollback()
    u.FileStore.Rollback()
    u.tx = nil
}
```

---

## Usage in a Handler

```go
func (h *CreateNoteHandler) Handle(ctx context.Context, raw json.RawMessage) (any, error) {
    var cmd createNoteCommand
    if err := json.Unmarshal(raw, &cmd); err != nil {
        return nil, err
    }

    uow := repositories.NewUnitOfWork()
    if err := uow.Begin(); err != nil {
        return nil, err
    }
    defer uow.Rollback()

    // Notes() is never called — never instantiated
    dbTags, err := uow.Tags().GetByNames(ctx, tags)
    if err != nil {
        return nil, err
    }

    // ...

    uow.FileStore.Stage(cmd.Path, data)
    return cmd.Path, uow.Commit()
}
```

---

## Properties of This Design

| Property | Result |
|---|---|
| Transaction leaks out of UoW | No — `tx` is private |
| Repos instantiated when not needed | No — lazy on first access |
| Handler controls which repos it uses | Yes — only accessed repos are created |
| UoW created per-request | Yes — inside `Handle`, not constructor |
| Handler has no persistent state | Yes — can be zero-value struct |
