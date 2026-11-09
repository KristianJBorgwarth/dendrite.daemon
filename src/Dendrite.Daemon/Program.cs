using System.Text;
using Dendrite.Application;
using Dendrite.Daemon;
using Dendrite.Daemon.Rpc;
using Microsoft.Extensions.DependencyInjection;

await using var services = Composition.BuildServices();

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
