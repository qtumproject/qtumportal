import { action, computed, observable } from "mobx"
import { AuthAPI } from "./AuthAPI"
import { IAuthorization } from "./types"

export class AuthStore {
  @observable public auths: IAuthorization[] = []
  @observable public counter: number = 0

  constructor(private _api: AuthAPI) {
    setInterval(() => {
      this.counter++
    }, 1000)
  }

  @computed get count(): number {
    return this.auths.length
  }

  public async loadAuthorizations() {
    const auths = await this._api.list()
    this.auths = auths
  }
}
