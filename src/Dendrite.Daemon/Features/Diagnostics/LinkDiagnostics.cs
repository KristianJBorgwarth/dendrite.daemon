using Dendrite.Daemon.Dtos;
using Dendrite.Daemon.Persistence.Repositories;
using Mediator;

namespace Dendrite.Daemon.Features.Diagnostics;

public sealed record LinkDiagnostics(string Path) : IQuery<List<LinkDiagnosticDto>>;

public sealed class LinkDiagnosticsHandler(NoteRepository notes, LinkRepository links)
    : IQueryHandler<LinkDiagnostics, List<LinkDiagnosticDto>>
{
    public ValueTask<List<LinkDiagnosticDto>> Handle(LinkDiagnostics query, CancellationToken cancellationToken)
    {
        var note = notes.GetByPath(query.Path);
        if (note is null)
            return ValueTask.FromResult<List<LinkDiagnosticDto>>([]);

        var diagnostics = links.GetBrokenLinks(note.Id)
            .Select(l => new LinkDiagnosticDto(
                l.Id,
                note.Path,
                l.Line,
                l.Col,
                l.Raw,
                l.TargetSlug,
                "Broken link: target note does not exist"))
            .ToList();

        return ValueTask.FromResult(diagnostics);
    }
}
