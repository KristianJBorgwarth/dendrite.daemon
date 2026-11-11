using System.Reflection;
using System.Text;
using Dendrite.Application;
using Dendrite.Application.Services;
using Dendrite.Daemon;
using Dendrite.Daemon.Rpc;
using Microsoft.Extensions.DependencyInjection;

if (args.Contains("--version"))
{
    var version = Assembly.GetEntryAssembly()!.GetCustomAttribute<AssemblyInformationalVersionAttribute>()!.InformationalVersion;
    Console.WriteLine(version.Split('+')[0]);
    return 0;
}

await using var services = Composition.BuildServices();
services.GetRequiredService<VaultRestorer>().Restore();

var utf8 = new UTF8Encoding(false);
using var input = new StreamReader(Console.OpenStandardInput(), utf8);
using var output = new StreamWriter(Console.OpenStandardOutput(), utf8) { AutoFlush = true };

try
{
    await services.GetRequiredService<RpcServer>().RunAsync(input, output);
}
catch (Exception ex)
{
    Log.Error($"server error: {ex}");
    return 1;
}

return 0;
