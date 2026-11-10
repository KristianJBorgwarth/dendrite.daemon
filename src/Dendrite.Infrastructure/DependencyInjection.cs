using Dendrite.Application.Abstractions;
using Dendrite.Domain.Models;
using Dendrite.Infrastructure.FileSystem;
using Dendrite.Infrastructure.Persistence;
using Dendrite.Infrastructure.Persistence.Queries;
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
        services.AddSingleton<INoteQueries, NoteQueries>();
        services.AddSingleton<IIndexMaintenance, IndexMaintenance>();
        services.AddSingleton<IVaultFiles, VaultFiles>();
        services.AddSingleton<IVaultScaffolder, VaultScaffolder>();
        return services;
    }
}
