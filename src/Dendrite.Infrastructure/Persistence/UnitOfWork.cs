using Dendrite.Application.Abstractions;

namespace Dendrite.Infrastructure.Persistence;

internal sealed class UnitOfWork : IUnitOfWork
{
    private readonly Database _db;

    public UnitOfWork(Database db)
    {
        _db = db;
    }

    public ITransaction Begin()
    {
        _db.BeginTransaction();
        return new Transaction(_db);
    }
}
