import * as React from "react"
import { AuthList } from "./AuthList"

export function App() {
  return (
    <div className="main-container">
      <link href="https://fonts.googleapis.com/css?family=Lato:100,300,400,700,900" rel="stylesheet" />
      <h1 className="app-title">QTUM Portal Authorizations</h1>
      <AuthList />
    </div>
  )
}
