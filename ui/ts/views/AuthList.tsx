import * as React from "react"

import { inject, observer } from "mobx-react"

import { AuthStore } from "../AuthStore"
import { IAuthorization } from "../types"

function AuthItem(props: { auth: IAuthorization, authStore: AuthStore }) {
  const { auth, authStore } = props
  const { method, params } = auth.request

  const createdAt = new Date(auth.createdAt)

  const dd = createdAt.getDate()
  // const yy = createdAt.getFullYear()
  const MM = createdAt.getMonth()

  const hh = createdAt.getHours()
  const mm = createdAt.getMinutes()
  const ss = createdAt.getSeconds()
  let ssString = ss.toString()
  if (ss < 10) {
    ssString = "0" + ssString
  }
  const createdAtFormated = `${MM}/${dd} ${hh}:${mm}:${ssString}`

  const isDenied = auth.state === "denied"

  return (
    <li className="auth-item">
      <div className="auth-item__topbar">
        <h3 className="auth-item__title">{method}</h3>
        <div className="auth-item__created">{createdAtFormated}</div>
      </div>

      <div className="auth-item__id">{auth.id}</div>
      <pre className="auth-item__params">
        {JSON.stringify(params, null, "  ")}
      </pre>

      {auth.state === "pending" &&
        <p className="auth-item__actions">
          <button className="btn btn--primary" onClick={() => authStore.accept(auth.id)}>Approve</button>
          <button className="btn btn--caution" onClick={() => authStore.deny(auth.id)}>Deny</button>
        </p>
      }

      {auth.state !== "pending" &&
        <p className={`auth-item__state ${isDenied && "auth-item__state--denied"}`}>{auth.state}</p>
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
      throw (new Error("Must provide authStore"))
    }

    const { authStore } = this.props
    const { auths } = authStore

    const notConnected = authStore.connState !== "connected"
    const hasNoAuths = auths.length === 0

    return (
      <div>
        {notConnected &&
          <p className="notice">
            Disconnected from QTUM Portal
          </p>
        }

        {!notConnected && hasNoAuths &&
          <p className="notice">
            No authorization is pending
          </p>
        }

        <ul className="auth-list">
          {auths.map((auth) => <AuthItem key={auth.id} auth={auth} authStore={authStore} />)}
        </ul>
      </div>

    )
  }
}
