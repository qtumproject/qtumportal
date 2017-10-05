import { action, computed, observable } from "mobx"
import { AuthAPI } from "./AuthAPI"
import { IAuthorization } from "./types"

export class AuthStore {
  @observable public _auths: Map<string, IAuthorization> = new Map()

  constructor(private _api: AuthAPI) {
  }

  @computed get count(): number {
    return this.auths.length
  }

  @computed get auths(): IAuthorization[] {
    const auths: IAuthorization[] = []
    for (const auth of this._auths.values())  {
      auths.push(auth)
    }

    return auths
  }

  public async loadAuthorizations() {
    const auths = await this._api.list()
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
}
