namespace Dendrite.Application.Abstractions;

/// <summary>
/// Couples index database work with staged file writes: files are only written once the
/// database work succeeded, and are removed again if the commit fails.
/// </summary>
public interface IUnitOfWork
{
    ITransaction Begin();
}
