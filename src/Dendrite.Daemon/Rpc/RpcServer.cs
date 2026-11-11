using System.Text.Json;
using System.Text.Json.Serialization.Metadata;
using Dendrite.Application;
using Mediator;

namespace Dendrite.Daemon.Rpc;

public sealed class RpcServer
{
    private const int ParseError = -32700;
    private const int InvalidParams = -32602;
    private const int MethodNotFound = -32601;
    private const int ServerError = -32000;

    private delegate Task<string> MethodHandler(JsonElement? parameters, CancellationToken ct);

    private readonly Dictionary<string, MethodHandler> _methods = [];
    private readonly VaultDispatcher _sender;

    public RpcServer(VaultDispatcher sender)
    {
        _sender = sender;
    }

    public RpcServer Query<TRequest, TResponse>(
        string method,
        JsonTypeInfo<TRequest> requestInfo,
        JsonTypeInfo<TResponse> responseInfo)
        where TRequest : IQuery<TResponse>
    {
        _methods[method] = async (parameters, ct) =>
        {
            var response = await _sender.Send(ReadParams(parameters, requestInfo), ct);
            return JsonSerializer.Serialize(response, responseInfo);
        };
        return this;
    }

    public RpcServer Command<TRequest, TResponse>(
        string method,
        JsonTypeInfo<TRequest> requestInfo,
        JsonTypeInfo<TResponse> responseInfo)
        where TRequest : ICommand<TResponse>
    {
        _methods[method] = async (parameters, ct) =>
        {
            var response = await _sender.Send(ReadParams(parameters, requestInfo), ct);
            return JsonSerializer.Serialize(response, responseInfo);
        };
        return this;
    }

    public RpcServer Command<TRequest>(string method, JsonTypeInfo<TRequest> requestInfo)
        where TRequest : ICommand<Unit>
    {
        _methods[method] = async (parameters, ct) =>
        {
            await _sender.Send(ReadParams(parameters, requestInfo), ct);
            return "null";
        };
        return this;
    }

    public async Task RunAsync(TextReader input, TextWriter output, CancellationToken ct = default)
    {
        while (await input.ReadLineAsync(ct) is { } line)
        {
            if (string.IsNullOrWhiteSpace(line))
                continue;

            var reply = await HandleLineAsync(line, ct);
            if (reply is null)
                continue;

            await output.WriteAsync(reply + "\n");
            await output.FlushAsync(ct);
        }
    }

    internal async Task<string?> HandleLineAsync(string line, CancellationToken ct)
    {
        RpcRequest? request;
        try
        {
            request = JsonSerializer.Deserialize(line, RpcJsonContext.Default.RpcRequest);
        }
        catch (JsonException)
        {
            request = null;
        }

        if (request?.Method is null)
            return Error("null", ParseError, "parse error");

        var id = request.Id is { ValueKind: not JsonValueKind.Null } element ? element.GetRawText() : null;

        try
        {
            if (!_methods.TryGetValue(request.Method, out var handler))
                return id is null ? null : Error(id, MethodNotFound, "method not found");

            var result = await handler(request.Params, ct);
            return id is null ? null : $$"""{"jsonrpc":"2.0","id":{{id}},"result":{{result}}}""";
        }
        catch (JsonException ex)
        {
            Log.Error($"{request.Method}: invalid params: {ex.Message}");
            return id is null ? null : Error(id, InvalidParams, ex.Message);
        }
        catch (Exception ex)
        {
            Log.Error($"{request.Method}: {ex}");
            return id is null ? null : Error(id, ServerError, ex.Message);
        }
    }

    private static TRequest ReadParams<TRequest>(JsonElement? parameters, JsonTypeInfo<TRequest> info)
    {
        var json = parameters is { ValueKind: JsonValueKind.Object } p ? p.GetRawText() : "{}";
        return JsonSerializer.Deserialize(json, info) ?? throw new JsonException("params must be an object");
    }

    private static string Error(string id, int code, string message)
    {
        var quoted = JsonSerializer.Serialize(message, RpcJsonContext.Default.String);
        return $$"""{"jsonrpc":"2.0","id":{{id}},"error":{"code":{{code}},"message":{{quoted}} } }""";
    }
}
