namespace Dendrite.Domain.Models;

#pragma warning disable IDE0290 // Use primary constructor
public sealed class Tag
{
    public string Name { get; }

    public Tag(string name)
    {
        Name = name;
    }
}
