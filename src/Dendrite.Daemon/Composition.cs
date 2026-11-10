using Dendrite.Application;
using Dendrite.Application.Abstractions;
using Dendrite.Daemon.Rpc;
using Dendrite.Infrastructure;
using Mediator;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon;

public static class Composition
{
    /// <param name="catalog">Where vaults are remembered; defaults to the on-disk catalog.</param>
    public static ServiceProvider BuildServices(IVaultCatalog? catalog = null)
    {
        Log.Configure(FileLog.Write);

        var services = new ServiceCollection();

        services.AddMediator();
        services.AddApplication();
        services.AddInfrastructure();
        if (catalog is not null)
            services.AddSingleton(catalog);
        services.AddSingleton<VaultDispatcher>();
        services.AddSingleton(sp => CreateRpcServer(sp.GetRequiredService<VaultDispatcher>()));

        return services.BuildServiceProvider();
    }

    private static RpcServer CreateRpcServer(VaultDispatcher sender)
    {
        var ctx = RpcJsonContext.Default;

        return new RpcServer(sender)
            .Command("vault/init", ctx.InitializeVault)
            .Command("vault/create", ctx.CreateVault)
            .Command("vault/rebuild", ctx.RebuildIndex)
            .Query("vault/list", ctx.ListVaults, ctx.ListVaultInfoDto)
            .Command("vault/remove", ctx.RemoveVault)
            .Command("vault/use", ctx.UseVault)
            .Command("vault/relocate", ctx.RelocateVault)
            .Command("vault/open", ctx.OpenVault)
            .Command("note/create", ctx.CreateNote, ctx.String)
            .Command("note/save", ctx.SaveNote)
            .Command("note/delete", ctx.DeleteNote)
            .Query("note/goto", ctx.GotoNote, ctx.GotoNoteResult)
            .Query("note/backlinks", ctx.GetBacklinks, ctx.ListBacklinkDto)
            .Query("note/search_by_tag", ctx.GetNotesByTag, ctx.ListNoteDto)
            .Query("note/search_by_cf", ctx.GetNotesByCf, ctx.ListNoteDto)
            .Query("completion/tag", ctx.CompleteTag, ctx.ListString)
            .Query("completion/slug", ctx.CompleteSlug, ctx.ListString)
            .Query("diagnostics/links", ctx.LinkDiagnostics, ctx.ListLinkDiagnosticDto);
    }
}
