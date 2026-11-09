using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Diagnostics;

public sealed class LinkDiagnosticsHandler(INoteQueries queries) 
    : IQueryHandler<LinkDiagnostics, List<LinkDiagnosticDto>>
{
    private readonly INoteQueries _queries = queries;

    public ValueTask<List<LinkDiagnosticDto>> Handle(LinkDiagnostics query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(_queries.GetLinkDiagnostics(query.Path));
}
