import "bulma/css/bulma.css"
import "font-awesome/css/font-awesome.css"

// tslint:disable-next-line:no-var-requires
require("./index.css")

import * as React from "react"
import { render } from "react-dom"

import { Provider } from "mobx-react"

import { AuthAPI } from "./AuthAPI"
import { AuthStore } from "./AuthStore"

import { App } from "./views/App"

async function init() {
  const authBaseURL = typeof AUTH_BASEURL === "undefined" ? window.location.origin : AUTH_BASEURL
  const authAPI = new AuthAPI(authBaseURL)
  const authStore = new AuthStore(authAPI)

  authStore.startAutoRefresh()

  Object.assign(window, {
    authStore,
  })

  // authStore.loadAuthorizations()

  const app = (
    <Provider authStore={authStore} >
      <App />
    </Provider>
  )
  render(app, document.getElementById("root"))
}

window.addEventListener("load", init)
