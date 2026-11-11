namespace Dendrite.Application.Abstractions;

public interface IUnitOfWork
{
    ITransaction Begin();
}
