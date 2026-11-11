using Dendrite.Application.Abstractions;

namespace Dendrite.Infrastructure.Persistence;

internal sealed class Transaction : ITransaction
{
    private readonly Database _db;
    private readonly List<(string Path, string Content)> _staged = [];
    private bool _completed;

    public Transaction(Database db)
    {
        _db = db;
    }

    public void Stage(string path, string content) => _staged.Add((path, content));

    public void Commit()
    {
        var written = new List<string>();
        try
        {
            foreach (var (path, content) in _staged)
            {
                if (File.Exists(path))
                    continue;

                Directory.CreateDirectory(Path.GetDirectoryName(path)!);
                File.WriteAllText(path, content);
                written.Add(path);
            }

            _db.Transaction!.Commit();
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
            _db.EndTransaction();
        }
    }

    public void Dispose()
    {
        if (!_completed)
            _db.EndTransaction();
    }
}
