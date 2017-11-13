import * as React from "react"
import { AuthList } from "./AuthList"

export function App() {
  return (
    <div>
      <nav className="navbar" role="navigation" aria-label="main navigation">
        <div className="navbar-brand">
          <span className="navbar-item">
            QTUM Portal Authorizations
        </span>
        </div>
      </nav>

      <section className="section">
        <div className="container content">
          <AuthList />
        </div>
      </section>
    </div>
  )
}
