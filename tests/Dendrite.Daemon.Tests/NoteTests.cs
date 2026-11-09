using AwesomeAssertions;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;

namespace Dendrite.Daemon.Tests;

public class NoteTests
{
    private static NoteFile File(string slug, string content) =>
        NoteFileParser.Parse("/vault", $"/vault/{slug}.md", content);

    [Fact]
    public void Create_DeduplicatesTagsAndCustomFrontMatter()
    {
        var note = Note.Create(File("a", "---\ntags: [x, x, y]\nstatus: [a, a]\n---\n[[b]]"));

        note.Tags.Should().Equal("x", "y");
        note.CustomFrontMatter.Should().ContainSingle();
        note.Links.Should().ContainSingle().Which.TargetSlug.Should().Be("b");
    }

    [Fact]
    public void Reindex_ReplacesChildrenButKeepsIdentity()
    {
        var note = Note.Create(File("a", "---\ntags: [x]\n---\n[[b]]"));
        var id = note.Id;
        var createdAt = note.CreatedAt;

        note.Reindex(File("a", "---\ntags: [y]\n---\n[[c]] [[d]]"));

        note.Id.Should().Be(id);
        note.CreatedAt.Should().Be(createdAt);
        note.Tags.Should().Equal("y");
        note.Links.Select(l => l.TargetSlug).Should().Equal("c", "d");
    }

    [Fact]
    public void Create_WithoutSlug_Throws()
    {
        var act = () => Note.Create(new NoteFile("/vault/a.md", "A", "", FrontMatter.Empty, []));

        act.Should().Throw<ArgumentException>();
    }
}
