import * as React from "react"

import { inject, observer } from "mobx-react"

import { AuthStore } from "../AuthStore"
import { IAuthorization } from "../types"

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

interface IAuthListProps {
  authStore?: AuthStore
}

@inject("authStore") @observer
export class AuthList extends React.Component<IAuthListProps, {}> {
  public render() {
    if (this.props.authStore == null) {
      throw(new Error("Must provide authStore"))
    }

    const { authStore } = this.props
    const { auths } = authStore

    return (
      <ul>
        {authStore.counter}
        {auths.map((auth) => <AuthItem auth={auth} />)}
      </ul>
    )
  }
}
