import { inject, observer } from "mobx-react"
import * as React from "react"

import { AuthStore } from "../AuthStore"
import { AuthItem } from "./AuthItem"

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
          <p className="notification is-warning">
            Disconnected from QTUM Portal
          </p>
        }

        {!notConnected && hasNoAuths &&
          <p className="notification">
            No authorization is pending
          </p>
        }

        <div className="auth-list">
          {auths.map((auth) => <AuthItem key={auth.id} auth={auth} authStore={authStore} />)}
        </div>
      </div>

    )
  }
}
