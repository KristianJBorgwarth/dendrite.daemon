using Dendrite.Application.Abstractions;

namespace Dendrite.Infrastructure.Persistence;

internal sealed class UnitOfWork(Database db) : IUnitOfWork
{
    public ITransaction Begin()
    {
        db.BeginTransaction();
        return new Transaction(db);
    }
}

internal sealed class Transaction(Database db) : ITransaction
{
    private readonly List<(string Path, string Content)> _staged = [];
    private bool _completed;

    public void Stage(string path, string content) => _staged.Add((path, content));

    public void Commit()
    {
        var written = new List<string>();
        try
        {
            foreach (var (path, content) in _staged)
            {
                // Never overwrite a file that already exists on disk.
                if (File.Exists(path))
                    continue;

                Directory.CreateDirectory(Path.GetDirectoryName(path)!);
                File.WriteAllText(path, content);
                written.Add(path);
            }

            db.Transaction!.Commit();
        }
        catch
        {
            foreach (var path in written)
                File.Delete(path);

            throw;
        }
        finally
        {
            _completed = true;
            db.EndTransaction();
        }
    }

    public void Dispose()
    {
        if (!_completed)
            db.EndTransaction();
    }
}
