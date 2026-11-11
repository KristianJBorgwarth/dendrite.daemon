using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Application;

public static class DependencyInjection
{
    public static IServiceCollection AddApplication(this IServiceCollection services)
    {
        services.AddSingleton<VaultRegistry>();
        services.AddSingleton<VaultContext>();
        services.AddSingleton<IndexRebuilder>();
        services.AddSingleton<TemplateRenderer>();
        services.AddSingleton<VaultInitializer>();
        services.AddSingleton<VaultRestorer>();
        services.AddSingleton<VaultOpener>();
        return services;
    }
}
