namespace Dendrite.Application.Abstractions;

public interface ITransaction : IDisposable
{
    void Stage(string path, string content);

    /// <summary>Writes staged files, then commits. Throws if either step fails.</summary>
    void Commit();
}
