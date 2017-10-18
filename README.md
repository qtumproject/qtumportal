# QTUM DApp Server

A Web Server for running third-party DApps.

All requests made by a DApp would block until you approve it in the authorization UI.

# Install qtum-portal

If you have golang installed, you can install the latest version:

```
go get -u github.com/hayeah/qtum-portal/cli/qtumportal
```

Or, you can download pre-built binaries from [Releases](https://github.com/hayeah/qtum-portal/releases)

# Running A QTUM DApp

First, we'll need to make sure that qtumd is running.

For testing/development purposes, let's start qtumd in regtest mode:

```
qtumd -regtest -rpcuser=howard -rpcpassword=yeh
```

Then use the env variable `QTUM_RPC` to specify the URL of your local qtumd RPC node:

```
export QTUM_RPC=http://howard:yeh@localhost:13889
```

Now we are ready to run the DApp. Clone an example DApp to your local machine:

```
$ git clone https://github.com/hayeah/qtum-dapp-getnewaddr.git
```

Use the `serve` command to start a web-server for the DApp:

```
$ qtumportal serve qtum-dapp-getnewaddr/build

INFO[0000] Serving DApp from /Users/howard/p/qtum/qtum-dapp-getnewaddr/build
INFO[0000] DApp service listening :9888
INFO[0000] Auth service listening :9899
```

Open the DApp in your browser at http://localhost:9888. Click the `getnewaddr` button to generate a new payment address. The request will be pending until you authorize (or deny) the request.

In another tab, open http://localhost:9899 to authorize transactions requested by the DApp.

# Developing QTUM DApp (Insecure)

Typically you'd be using a live-reload server when developing your HTML5 DApp. By default QTUM portal provides security by locking down the DApp from making cross domain requests. If you are using a dev server for development purposes, you'll need to disable CORS protection.

The example DApp [getnewaddr](qtum-dapp-getnewaddr) uses [Neutrino](https://neutrino.js.org/) for project building.

Start the Neutrino live-reload dev-server:

```
$ cd qtum-dapp-getnewaddr
$ npm start
✔ Development server running on: http://localhost:3000
✔ Build completed
```

Then start qtumportal in dev-mode:

```
$ qtumportal serve --dev
INFO[0000] Auth service listening :9899
INFO[0000] DApp service listening :9888
```

In your DApp code, configure the QtumRPC client to use `http://localhost:9888` as the RPC endpoint:

```
let rpcURL: string

if (process.env.NODE_ENV === "development") {
  rpcURL = "http://localhost:9888"
} else {
  // In production mode, make RPC request to origin host.
  rpcURL = window.location.origin
}

const rpc = new QtumRPC(rpcURL)
```

Open http://localhost:3000 in browser to develop using the Neutrino live-reload server.

* * *

If you are not using a live-reload server, you could also serve a static directory in dev mode:

```
$ qtumportal serve examples/getnewaddr/build --dev
INFO[0000] DApp service listening :9888
INFO[0000] Auth service listening :9899
```

Open http://localhost:9888 in browser.