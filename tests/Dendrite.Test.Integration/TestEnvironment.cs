using System.Runtime.CompilerServices;

namespace Dendrite.Test.Integration;

internal static class TestEnvironment
{
    [ModuleInitializer]
    internal static void Init() =>
        Environment.SetEnvironmentVariable("XDG_DATA_HOME", Path.Combine(Path.GetTempPath(), "dendrite-tests-xdg"));
}
