import * as React from "react"
import { render } from "react-dom"

import { Provider } from "mobx-react"

import { AuthAPI } from "./AuthAPI"
import { AuthStore } from "./AuthStore"

import { AuthList } from "./views/AuthList"

let QPORTAL_AUTH_BASEURL: string

if (Object.is(process.env.NODE_ENV, "development")) {
  QPORTAL_AUTH_BASEURL = "http://localhost:9898"
}

async function init() {
  const authAPI = new AuthAPI(QPORTAL_AUTH_BASEURL)
  const authStore = new AuthStore(authAPI)

  authStore.startAutoRefresh()

  Object.assign(window, {
    authStore,
  })

  authStore.loadAuthorizations()

  const app = (
    <Provider authStore={authStore} >
      <div>
        <h1>Authorizations</h1>
        <AuthList/>
      </div>
    </Provider>
  )
  render(app, document.getElementById("root"))
}

window.addEventListener("load", init)
