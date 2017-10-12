// tslint:disable-next-line:no-var-requires
require("normalize.css")
require("../css/index.css")

import * as React from "react"
import { render } from "react-dom"

import { Provider } from "mobx-react"

import { AuthAPI } from "./AuthAPI"
import { AuthStore } from "./AuthStore"

import { App } from "./views/App"

if (Object.is(process.env.NODE_ENV, "development")) {
  const QTUMPORTAL_CONFIG = {
    AUTH_BASEURL: "http://localhost:9899",
  }

  Object.assign(window, {
    QTUMPORTAL_CONFIG,
  })
}

async function init() {
  const authAPI = new AuthAPI(QTUMPORTAL_CONFIG.AUTH_BASEURL)
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
