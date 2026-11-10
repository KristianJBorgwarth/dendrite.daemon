using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;

namespace Dendrite.Application.Services;

public sealed class IndexRebuilder
{
    private readonly IUnitOfWork _uow;
    private readonly INoteRepository _notes;
    private readonly IIndexMaintenance _index;
    private readonly IVaultFiles _files;
    private readonly VaultContext _vault;

    public IndexRebuilder(
        IUnitOfWork uow,
        INoteRepository notes,
        IIndexMaintenance index,
        IVaultFiles files,
        VaultContext vault)
    {
        _uow = uow;
        _notes = notes;
        _index = index;
        _files = files;
        _vault = vault;
    }

    public void Rebuild()
    {
        var root = _vault.Config.VaultPath;
        var parsed = _files.ListNotePaths()
            .Select(path => NoteFileParser.Parse(root, path, _files.ReadText(path)))
            .ToList();

        using var tx = _uow.Begin();

        _index.Wipe();

        foreach (var file in parsed)
            _notes.Add(Note.Create(file));

        tx.Commit();
    }
}
