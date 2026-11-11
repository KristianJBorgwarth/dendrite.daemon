using AwesomeAssertions;
using Dendrite.Application.Features.Completion;
using Dendrite.Application.Features.Diagnostics;
using Dendrite.Application.Features.Notes;
using Dendrite.Application.Features.Vault;

namespace Dendrite.Test.Integration;

public class NoteHandlerTests
{
    [Fact]
    public async Task CreateNote_NoTemplate_CreatesFileAndIndexRow()
    {
        using var vault = new TestVault();

        var path = await vault.Sender.Send(new CreateNote("My Note", "", ""));

        path.Should().Be(Path.Combine(vault.Root, "my-note.md"));
        File.Exists(path).Should().BeTrue();
        vault.Count("SELECT COUNT(*) FROM note WHERE slug = 'my-note'").Should().Be(1);
    }

    [Fact]
    public async Task CreateNote_WithTemplate_IndexesTemplateTags()
    {
        using var vault = new TestVault();
        vault.Write(".templates/template.md", "---\ntitle: Template\ntags: [go, testing]\n---\n");

        var path = await vault.Sender.Send(new CreateNote("Templated Note", "template.md", "subdir"));

        path.Should().Be(Path.Combine(vault.Root, "subdir", "templated-note.md"));
        File.Exists(path).Should().BeTrue();
        vault.Count("SELECT COUNT(*) FROM note WHERE slug = 'subdir/templated-note'").Should().Be(1);
        vault.Count("SELECT COUNT(*) FROM tag WHERE name IN ('go', 'testing')").Should().Be(2);
    }

    [Fact]
    public async Task CreateNote_DuplicateSlug_ReturnsExistingPathWithoutSecondRow()
    {
        using var vault = new TestVault();
        var command = new CreateNote("Dup Note", "", "dup-dir");

        var first = await vault.Sender.Send(command);
        var second = await vault.Sender.Send(command);

        second.Should().Be(first);
        vault.Count("SELECT COUNT(*) FROM note WHERE slug = 'dup-dir/dup-note'").Should().Be(1);
    }

    [Fact]
    public async Task CreateNote_LeadingSlashDirectory_MatchesTheSlugProducedOnSave()
    {
        using var vault = new TestVault();

        var path = await vault.Sender.Send(new CreateNote("Today", "", "/daily"));
        await vault.Sender.Send(new SaveNote(path));

        vault.Count("SELECT COUNT(*) FROM note WHERE slug = 'daily/today'").Should().Be(1);
        vault.Count("SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task CreateNote_UnknownTemplate_Throws()
    {
        using var vault = new TestVault();

        var act = async () => await vault.Sender.Send(new CreateNote("Ghost", "non-existent-template", ""));

        await act.Should().ThrowAsync<IOException>();
        vault.Count("SELECT COUNT(*) FROM note").Should().Be(0);
    }

    [Fact]
    public async Task SaveNote_ReindexesLinksTagsAndCustomFrontMatter()
    {
        using var vault = new TestVault();
        vault.Write("target.md", "---\ntitle: Target\n---\n");
        var source = vault.Write("source.md", "---\ntitle: Source\ntags: [a]\nstatus: draft\n---\nsee [[target]] and [[missing]]");
        await vault.Sender.Send(new RebuildIndex());

        File.WriteAllText(source, "---\ntitle: Source\ntags: [b]\nstatus: done\n---\nsee [[target]] and [[missing]]");
        await vault.Sender.Send(new SaveNote(source));

        vault.Count("SELECT COUNT(*) FROM note_tag WHERE tag_id = 'a'").Should().Be(0);
        vault.Count("SELECT COUNT(*) FROM note_tag WHERE tag_id = 'b'").Should().Be(1);
        vault.Count("SELECT COUNT(*) FROM link").Should().Be(2);

        var byCf = await vault.Sender.Send(new GetNotesByCf("status", "done"));
        byCf.Select(n => n.Slug).Should().Equal("source");

        var backlinks = await vault.Sender.Send(new GetBacklinks(Path.Combine(vault.Root, "target.md")));
        backlinks.Should().ContainSingle().Which.Slug.Should().Be("source");

        var diagnostics = await vault.Sender.Send(new LinkDiagnostics(source));
        diagnostics.Should().ContainSingle().Which.Target.Should().Be("missing");
    }

    [Fact]
    public async Task DeleteNote_RemovesNoteAndMetadata()
    {
        using var vault = new TestVault();
        var path = vault.Write("gone.md", "---\ntitle: Gone\ntags: [x]\n---\n[[other]]");
        await vault.Sender.Send(new RebuildIndex());

        await vault.Sender.Send(new DeleteNote(path));

        vault.Count("SELECT COUNT(*) FROM note").Should().Be(0);
        vault.Count("SELECT COUNT(*) FROM link").Should().Be(0);
        vault.Count("SELECT COUNT(*) FROM note_tag").Should().Be(0);
    }

    [Fact]
    public async Task DeleteNote_UnknownPath_ReportsAnError()
    {
        using var vault = new TestVault();

        var act = async () => await vault.Sender.Send(new DeleteNote(Path.Combine(vault.Root, "nope.md")));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("note not found*");
    }

    [Fact]
    public async Task GotoNote_ResolvesNotesAndUrls()
    {
        using var vault = new TestVault(root => File.WriteAllText(Path.Combine(root, "a.md"), "x"));

        (await vault.Sender.Send(new GotoNote("[[a]]")))!.Type.Should().Be("note");
        (await vault.Sender.Send(new GotoNote("[w](https://x.io)")))!.Target.Should().Be("https://x.io");
        (await vault.Sender.Send(new GotoNote("[[nope]]"))).Should().BeNull();
    }

    [Fact]
    public async Task Completion_SearchesSlugsAndTags()
    {
        using var vault = new TestVault(root => File.WriteAllText(Path.Combine(root, "alpha.md"), "---\ntags: [rust, go]\n---\n"));

        (await vault.Sender.Send(new CompleteSlug("alp"))).Should().Equal("alpha");
        (await vault.Sender.Send(new CompleteTag("ru"))).Should().Equal("rust");
    }

    [Fact]
    public async Task RebuildIndex_HonoursConfiguredExcludes()
    {
        using var vault = new TestVault(
            root =>
            {
                Directory.CreateDirectory(Path.Combine(root, "private"));
                Directory.CreateDirectory(Path.Combine(root, ".templates"));
                File.WriteAllText(Path.Combine(root, "keep.md"), "x");
                File.WriteAllText(Path.Combine(root, "private", "secret.md"), "x");
                File.WriteAllText(Path.Combine(root, ".templates", "t.md"), "x");
            },
            () => new VaultConfigDto { ExcludeIndexFiles = ["private"] });

        (await vault.Sender.Send(new CompleteSlug(""))).Should().Equal("keep");
    }

    [Fact]
    public async Task RebuildIndex_OverrideDefaultIgnores_DropsBuiltInExcludes()
    {
        using var vault = new TestVault(
            root =>
            {
                Directory.CreateDirectory(Path.Combine(root, ".templates"));
                File.WriteAllText(Path.Combine(root, ".templates", "t.md"), "x");
            },
            () => new VaultConfigDto { OverrideDefaultIgnores = true });

        (await vault.Sender.Send(new CompleteSlug(""))).Should().Equal(".templates/t");
    }
}
