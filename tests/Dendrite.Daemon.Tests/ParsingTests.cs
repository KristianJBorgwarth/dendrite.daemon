using AwesomeAssertions;
using Dendrite.Daemon.FileHandling;

namespace Dendrite.Daemon.Tests;

public class FrontMatterTests
{
    [Fact]
    public void Parse_ValidFrontMatter_ReturnsTitleAndTags()
    {
        var result = FrontMatter.Parse("---\ntitle: My Note\ntags: [\"test\", \"note\"]\n---\nThis is the content.");

        result.Should().NotBeNull();
        result!.Title.Should().Be("My Note");
        result.Tags.Should().Equal("test", "note");
    }

    [Fact]
    public void Parse_MissingDelimiter_ReturnsNull()
    {
        FrontMatter.Parse("title: My Note\ntags: [\"test\"]\nThis is the content.").Should().BeNull();
    }

    [Fact]
    public void Parse_CustomKeys_AreFlattenedToKeyValuePairs()
    {
        var result = FrontMatter.Parse("---\ntitle: T\nstatus: draft\nauthors: [ann, bob]\ndate: 2026-01-01\n---\n");

        result!.Custom.Should().Equal(
            new KeyValuePair<string, string>("status", "draft"),
            new KeyValuePair<string, string>("authors", "ann"),
            new KeyValuePair<string, string>("authors", "bob"));
    }

    [Fact]
    public void Parse_EmptyBlock_ReturnsEmptyFrontMatter()
    {
        FrontMatter.Parse("---\n---\nbody").Should().Be(FrontMatter.Empty);
    }
}

public class LinkTests
{
    [Fact]
    public void ExtractLinks_ReportsLineAndByteColumn()
    {
        var links = NoteFileReader.ExtractLinks("first\nåå [[a/b|Shown]] and [[c]]");

        links.Should().HaveCount(2);
        links[0].Should().Be(new ExtractedLink("[[a/b|Shown]]", "Shown", "a/b", 2, 6));
        links[1].TargetSlug.Should().Be("c");
    }

    [Theory]
    [InlineData("[[my-note]]", LinkKind.Note, "my-note")]
    [InlineData("[[my-note|Alias]]", LinkKind.Note, "my-note")]
    [InlineData("[site](https://example.com)", LinkKind.Url, "https://example.com")]
    [InlineData("bare-slug", LinkKind.Note, "bare-slug")]
    public void ParseLink_ClassifiesTargets(string input, LinkKind kind, string target)
    {
        LinkParser.Parse(input).Should().Be(new ParsedLink(kind, target));
    }

    [Fact]
    public void Slugify_LowercasesAndReplacesEverythingElse()
    {
        Slug.Slugify("My Note 2!").Should().Be("my-note-2-");
    }
}
