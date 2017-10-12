import { action, computed, observable } from "mobx"
import { AuthAPI } from "./AuthAPI"
import { IAuthorization } from "./types"

export class AuthStore {
  @observable public _auths: Map<string, IAuthorization> = new Map()
  @observable public connState: "connected" | "connecting" | "disconnected" = "disconnected"

  constructor(private _api: AuthAPI) {
  }

  @computed get count(): number {
    return this.auths.length
  }

  @computed get auths(): IAuthorization[] {
    const auths: IAuthorization[] = []
    for (const auth of this._auths.values()) {
      auths.push(auth)
    }

    return auths
  }

  @computed get pendingAuths(): IAuthorization[] {
    return this.auths.filter((auth) => auth.state === "pending")
  }

  public startAutoRefresh() {
    if (this.connState !== "disconnected") {
      return
    }

    const so = this._api.eventsSocket()
    this.connState = "connecting"

    so.addEventListener("close", (e) => {
      this.onDisconnect()
    })

    so.addEventListener("error", (e) => {
      this.onDisconnect()
    })

    so.addEventListener("open", () => {
      this.connState = "connected"
      this.loadAuthorizations()
    })

    so.addEventListener("message", (e) => {
      if (e.data === "refresh") {
        this.loadAuthorizations()
      }
    })
  }

  public async loadAuthorizations() {
    const auths = await this._api.list()

    this._auths = new Map()
    auths.forEach((auth) => {
      this._auths.set(auth.id, auth)
    })

    /* Babel error */
    // for (const auth of auths) {
    //   this._auths.set(auth.id, auth)
    // }
  }

  public async accept(id: string) {
    const auth = await this._api.accept(id)
    this._auths.set(auth.id, auth)
  }

  public async deny(id: string) {
    const auth = await this._api.deny(id)
    this._auths.set(auth.id, auth)
  }

  private onDisconnect() {
    this._auths = new Map()
    this.connState = "disconnected"
    setTimeout(this.startAutoRefresh.bind(this), 2000)
  }
}
