using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Application;

public static class DependencyInjection
{
    /// <summary>Registers application services. The Mediator itself is generated in the host (see Dendrite.Daemon).</summary>
    public static IServiceCollection AddApplication(this IServiceCollection services)
    {
        services.AddSingleton<VaultContext>();
        services.AddSingleton<NoteMetadataService>();
        services.AddSingleton<IndexRebuilder>();
        services.AddSingleton<TemplateRenderer>();
        return services;
    }
}
