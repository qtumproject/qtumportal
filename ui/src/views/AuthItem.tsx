import * as React from "react"

import { AuthStore } from "../AuthStore"
import { IAuthorization } from "../types"

export function AuthItem(props: { auth: IAuthorization, authStore: AuthStore }) {
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

  const css = {
    titleBox: {
      marginBottom: "1em",
    },
  }

  return (
    <div className="box content">
      <div style={css.titleBox}>
        <h3 className="is-marginless">
          {method}
          <span className="is-pulled-right is-size-6">{createdAtFormated}</span>
        </h3>
        <div>{auth.id}</div>
      </div>

      <p>
        <h4> Params </h4>

        <pre>
          {JSON.stringify(params, null, "  ")}
        </pre>
      </p>

      {auth.state === "pending" &&
        <p className="buttons">
          <button className="button is-success" onClick={() => authStore.accept(auth.id)}>Approve</button>
          <button className="button is-danger" onClick={() => authStore.deny(auth.id)}>Deny</button>
        </p>
      }

      {auth.state !== "pending" &&
        <h4 className={`has-text-centered is-capitalized ${isDenied && "has-text-danger"}`}>
          {auth.state}
        </h4>
      }
    </div>
  )
}
