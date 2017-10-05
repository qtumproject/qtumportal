# QTUM Portal Authorization Design

We cannot trust 3rd party DApps to call `qtumd`'s RPC methods directly. `qtum-portal` adds an authorization layer for users to grant permission to security sensitive RPC calls.

Since `qtum-portal` serves the HTML5 assets of a DApp, browser's CORS policy would restrict the DApp to only send JSON-RPC requests to `qtum-portal`, the origin server. These RPC calls would then be regulated according to `qtum-portal`'s security policy.

`qtum-portal` acts as a proxy that passes JSON-RPC calls to an instance of qtumd RPC service.

+ Read-only RPC calls are passed directly to `qtumd` unmodified.
+ RPC calls that create transactions or modify `qtumd` state would require user authorization.
+ Operation related RPC calls are diabled. e.g. `dumpwallet`,`clearbanned`, `prioritisetransaction`, etc.

For a list of supported RPC methods, as well as which methods that require user authorization, see: [methods.go](https://github.com/hayeah/qtum-portal/blob/master/methods.go).

# RPC Call Without Authorization

Let's consider the `getnewaddress` method call, which generates a new payment address. The `qtumd` RPC service, given user:password, does not require any user authorization:

```
curl http://howard:yeh@localhost:13889/ -X POST -H "Content-Type: application/json" -d '
{
  "jsonrpc": "1.0",
  "id":"1",
  "method": "getnewaddress",
  "params": []
}
'
```

This results in:

```
{"result":"qZXcpztDoeUu3ADWDJf5NKfevMaD9ZBUiF","error":null,"id":"1"}
```

# Authorization Required

Assuming that `qtum-portal` is running on port 9999, we can try to make the same `getnewaddress` call:

```
curl http://howard:yeh@localhost:9999/ -X POST -H "Content-Type: application/json" -d '
{
  "jsonrpc": "1.0",
  "id":"1",
  "method": "getnewaddress",
  "params": []
}
'
```

Now we see that the call fails with `402 Payment Required`,

```
HTTP/1.1 402 Payment Required
Access-Control-Allow-Origin: *
Content-Type: application/json; charset=UTF-8
Vary: Origin
Date: Thu, 05 Oct 2017 09:40:37 GMT
Content-Length: 231

{"id":"2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo","state":"pending","request":{"method":"getnewaddress","id":"1","params":[],"auth":"2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo"},"createdAt":"2017-10-05T17:40:37.780820541+08:00"}
```

The returned JSON object contains an authorization ID:

```
2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo
```

Given this authorization ID, user may accept or deny the RPC call, either programmatically, or with an UI.

# Authorization Flow

The authorization flow goes like this:

1. DApp client makes an RPC call that requires authorization.
2. Server returns `402 Payment Required`, returning an authorization object to the client.
3. User approves the authorization object. With an API call or with `qtum-portal` UI.
4. DApp checks (or notified) that an authorization is approved. If so, DApp makes the RPC request again, this time attaching the authorization id.
5. Server checks the authorization id to see if the RPC call has the same method and parameters. If so, it passes the RPC call to the underlying `qtumd` RPC service.
	+ An authorization may only be used once.

The authorization API is listening on a different port (9898). For step 3, the API call is `POST /authorizations/:id/accept`:

```
$ curl localhost:9898/authorizations/2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo/accept \
  -X POST -H "Content-Type: application/json"

{
  "createdAt": "2017-10-05T17:40:37.780820541+08:00",
  "request": {
    "auth": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo",
    "params": [],
    "id": "1",
    "method": "getnewaddress"
  },
  "state": "accepted",
  "id": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo"
}
```

We see that the authorization object's state transitioned to "accepted", we may now make the RPC call, attaching the authorization id:

```
$ curl localhost:9999/ -X POST -H "Content-Type: application/json" -d '
{
  "jsonrpc": "1.0",
  "id":"1",
  "method": "getnewaddress",
  "params": [],
  "auth": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo"
}
'
```

This time the call succeeds:

```
{
	"result":"qZD4Egn23UgYc2g9MdmB4CbtW9SKmbjVZJ",
	"error":null,
	"id":"1"
}
```

Checking the authorization object again, we can see that the state had transitioned to `consumed` after the RPC call:

```
$ curl localhost:9898/authorizations/2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo

{
  "createdAt": "2017-10-05T17:40:37.780820541+08:00",
  "request": {
    "auth": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo",
    "params": [],
    "id": "1",
    "method": "getnewaddress"
  },
  "state": "consumed",
  "id": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo"
}
```

Making the RPC call again with the same authorization id should result in error:

```
$ curl localhost:9999/ -X POST -H "Content-Type: application/json" -d '
{
  "jsonrpc": "1.0",
  "id":"1",
  "method": "getnewaddress",
  "params": [],
  "auth": "2XMfQmXlUABQDjyP_0XZn1BotRRyH5-vGNnDJP0usdo"
}
'
```
```
{"message":"Cannot verify RPC request"}
```