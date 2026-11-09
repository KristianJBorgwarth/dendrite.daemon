using Dendrite.Application.Abstractions;
using Dendrite.Infrastructure.FileSystem;
using Dendrite.Infrastructure.Persistence;
using Dendrite.Infrastructure.Persistence.Repositories;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Infrastructure;

public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services)
    {
        services.AddSingleton<Database>();
        services.AddSingleton<IIndexDatabase>(sp => sp.GetRequiredService<Database>());
        services.AddSingleton<IUnitOfWork, UnitOfWork>();
        services.AddSingleton<INoteRepository, NoteRepository>();
        services.AddSingleton<ITagRepository, TagRepository>();
        services.AddSingleton<ILinkRepository, LinkRepository>();
        services.AddSingleton<ICfeRepository, CfeRepository>();
        services.AddSingleton<IIndexRepository, IndexRepository>();
        services.AddSingleton<IVaultFiles, VaultFiles>();
        return services;
    }
}
