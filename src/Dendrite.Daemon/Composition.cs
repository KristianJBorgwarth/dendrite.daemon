using Dendrite.Application;
using Dendrite.Daemon.Rpc;
using Dendrite.Infrastructure;
using Mediator;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon;

public static class Composition
{
    public static ServiceProvider BuildServices()
    {
        Log.Configure(FileLog.Write);

        var services = new ServiceCollection();

        services.AddMediator();
        services.AddApplication();
        services.AddInfrastructure();
        services.AddSingleton(sp => CreateRpcServer(sp.GetRequiredService<ISender>()));

        return services.BuildServiceProvider();
    }

    private static RpcServer CreateRpcServer(ISender sender)
    {
        var ctx = RpcJsonContext.Default;

        return new RpcServer(sender)
            .Command("vault/init", ctx.InitializeVault)
            .Command("vault/rebuild", ctx.RebuildIndex)
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
