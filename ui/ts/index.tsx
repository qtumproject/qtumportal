import * as React from "react"
import { render } from "react-dom"

if (Object.is(process.env.NODE_ENV, "development")) {
  Object.assign(window, {
    QPORTAL_AUTH_BASEURL: "http://localhost:9898/",
    // authBaseURL: process.env.QPORTAL_AUTH_BASEURL,
  })
}

interface IJSONRPCRequest {
  id: any
  method: string
  params: any[]
  auth: string
}

interface IAuthorization {
  id: string
  state: "pending" | "accepted" | "denied" | "consumed"
  request: IJSONRPCRequest
  createdAt: string
}

class AuthAPI {
  constructor(private _baseURL: string) {
  }

  public async list(): Promise<IAuthorization[]> {
    const res = await fetch(this._baseURL + "authorizations")
    const auths = await res.json()
    return auths || []
  }

}

function AuthItem(props: { auth: IAuthorization }) {
  const { auth } = props
  const { method, params } = auth.request
  return (
    <li>
      <h3>{method}</h3>
      <pre>
      {JSON.stringify(params, null, "  ")}
      </pre>
      <p>{auth.createdAt}</p>
      <button>Approve</button>
      <button>Deny</button>
    </li>
  )
}

function AuthList(props: { auths: IAuthorization[] }) {
  const { auths } = props

  return (
    <ul>
      {auths.map((auth) => <AuthItem auth={auth}/> )}
    </ul>
  )
}

async function init() {
  const authAPI = new AuthAPI(QPORTAL_AUTH_BASEURL)
  const auths = await authAPI.list()

  const app = (
    <div>
      <h1>Authorizations</h1>
      <AuthList auths={auths}/>
    </div>
  )
  render(app, document.getElementById("root"))
}

window.addEventListener("load", init)
