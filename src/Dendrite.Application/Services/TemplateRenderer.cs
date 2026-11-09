using Dendrite.Application.Abstractions;
using Dendrite.Domain.FileHandling;

namespace Dendrite.Application.Services;

public sealed class TemplateRenderer
{
    private readonly IVaultFiles _files;

    public TemplateRenderer(IVaultFiles files)
    {
        _files = files;
    }

    public Template Render(string templateName, string title, string directory)
    {
        var fileNameSlug = Slug.Slugify(title);
        var filename = fileNameSlug + ".md";
        var slug = Path.Join(directory.Trim('/', '\\'), fileNameSlug).Replace('\\', '/');

        var raw = templateName == ""
            ? DefaultTemplate(title)
            : _files.ReadTemplate(templateName);

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
