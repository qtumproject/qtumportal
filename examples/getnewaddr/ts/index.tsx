import { QtumRPC } from "qtumjs"
import * as React from "react"
import { render } from "react-dom"

let rpcURL: string

if (process.env.NODE_ENV === "development") {
  rpcURL = "http://localhost:9888"
} else {
  rpcURL = window.location.origin
}

const rpc = new QtumRPC(rpcURL)

class App extends React.Component<{}, { addresses: string[], isLoading: boolean }> {
  constructor(props: any) {
    super(props)
    this.state = {
      addresses: [],
      isLoading: false,
    }
  }

  public render() {
    const {
      addresses,
      isLoading,
    } = this.state
    return (
      <div>
        <button onClick={this.createNewAddress.bind(this)} disabled={isLoading}>Generate Payment Address</button>

        {isLoading && "Generating..."}

        <ul>
          {this.state.addresses.map((address, i) => <li key={i}>{address}</li>)}
        </ul>
      </div>
    )
  }

  private async createNewAddress() {
    try {
      this.setState({ isLoading: true })
      const address = await rpc.rawCall("getnewaddress")

      this.setState({
        addresses: [address, ...this.state.addresses],
      })
    } catch (err) {
      alert(err)
    } finally {
      this.setState({ isLoading: false })
    }
  }
}

window.addEventListener("load", () => {
  render(<App />, document.getElementById("root"))
})
