using Mediator;

namespace Dendrite.Application.Features.Diagnostics;

public sealed record LinkDiagnostics : IQuery<List<LinkDiagnosticDto>>
{
    public string Path { get; }

    public LinkDiagnostics(string path)
    {
        Path = path;
    }
}
