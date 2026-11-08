using Dendrite.Daemon.Persistence;

namespace Dendrite.Daemon.FileHandling;

public sealed record Template(string Content, string Title, string Slug, string Filename, FrontMatter FrontMatter);

public sealed class TemplateRenderer(VaultStore vaultStore)
{
    public Template Render(string templateName, string title, string directory)
    {
        var fileNameSlug = Slug.Slugify(title);
        var filename = fileNameSlug + ".md";
        var slug = Path.Join(directory.Trim('/', '\\'), fileNameSlug).Replace('\\', '/');

        var raw = templateName == ""
            ? DefaultTemplate(title)
            : File.ReadAllText(vaultStore.GetTemplatePath(templateName));

        var content = raw
            .Replace("{{title}}", title)
            .Replace("{{date}}", Today())
            .Replace("{{slug}}", slug)
            .Replace("{{file}}", filename);

        var frontMatter = FrontMatter.Parse(content) ?? FrontMatter.Empty;
        return new Template(content, title, slug, filename, frontMatter);
    }

    private static string Today() => DateTime.Now.ToString("yyyy-MM-dd");

    private static string DefaultTemplate(string title) => $"---\ntitle: {title}\ndate: {Today()}\n---\n";
}
