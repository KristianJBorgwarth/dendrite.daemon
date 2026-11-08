using System.Text;

namespace Dendrite.Daemon.FileHandling;

public static class Slug
{
    public static string Slugify(string s)
    {
        var slug = new StringBuilder(s.Length);
        foreach (var r in s.EnumerateRunes())
        {
            if (r.Value is >= 'a' and <= 'z' or >= '0' and <= '9')
                slug.Append((char)r.Value);
            else if (r.Value is >= 'A' and <= 'Z')
                slug.Append((char)(r.Value + 32));
            else
                slug.Append('-');
        }
        return slug.ToString();
    }

    public static string Relative(string vaultRoot, string path)
    {
        var prefix = vaultRoot + Path.DirectorySeparatorChar;
        var rel = path.StartsWith(prefix, StringComparison.Ordinal) ? path[prefix.Length..] : path;
        var ext = Path.GetExtension(rel);
        return ext.Length > 0 ? rel[..^ext.Length] : rel;
    }
}
