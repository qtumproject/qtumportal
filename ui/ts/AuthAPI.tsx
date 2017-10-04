import { IAuthorization } from "./types"

export class AuthAPI {
  constructor(private _baseURL: string) {
  }

  public async list(): Promise<IAuthorization[]> {
    const res = await fetch(this._baseURL + "authorizations")
    const auths = await res.json()
    return auths || []
  }
}
