namespace Dendrite.Domain.Vault;

public sealed record DailyNotes
{
    public string Dir { get; }
    public string FilenameFormat { get; }
    public string TemplateName { get; }

    public DailyNotes(string dir, string filenameFormat, string templateName)
    {
        Dir = dir;
        FilenameFormat = filenameFormat;
        TemplateName = templateName;
    }
}
