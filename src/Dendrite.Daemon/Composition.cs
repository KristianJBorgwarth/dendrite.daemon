using Dendrite.Daemon.FileHandling;
using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;
using Dendrite.Daemon.Rpc;
using Dendrite.Daemon.Services;
using Mediator;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon;

public static class Composition
{
    public static ServiceProvider BuildServices()
    {
        var services = new ServiceCollection();

        services.AddMediator();

        services.AddSingleton<Database>();
        services.AddSingleton<VaultStore>();
        services.AddSingleton<UnitOfWork>();
        services.AddSingleton<NoteRepository>();
        services.AddSingleton<TagRepository>();
        services.AddSingleton<LinkRepository>();
        services.AddSingleton<CfeRepository>();
        services.AddSingleton<IndexRepository>();
        services.AddSingleton<NoteMetadataService>();
        services.AddSingleton<IndexRebuilder>();
        services.AddSingleton<TemplateRenderer>();
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
