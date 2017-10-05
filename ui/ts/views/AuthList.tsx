import * as React from "react"

import { inject, observer } from "mobx-react"

import { AuthStore } from "../AuthStore"
import { IAuthorization } from "../types"

function AuthItem(props: { auth: IAuthorization, authStore: AuthStore }) {
  const { auth, authStore } = props
  const { method, params } = auth.request
  return (
    <li>
      <h3>{method}</h3>
      <p>{auth.state}</p>
      <pre>
        {JSON.stringify(params, null, "  ")}
      </pre>
      <p>{auth.createdAt}</p>
      { auth.state === "pending" &&
        <p>
          <button onClick={() => authStore.accept(auth.id)}>Approve</button>
          <button onClick={() => authStore.deny(auth.id)}>Deny</button>
        </p>
      }
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
        {auths.map((auth) => <AuthItem key={auth.id} auth={auth} authStore={authStore}/>)}
      </ul>
    )
  }
}
