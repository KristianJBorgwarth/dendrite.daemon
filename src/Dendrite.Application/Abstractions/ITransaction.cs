namespace Dendrite.Application.Abstractions;

public interface ITransaction : IDisposable
{
    void Stage(string path, string content);

    void Commit();
}
