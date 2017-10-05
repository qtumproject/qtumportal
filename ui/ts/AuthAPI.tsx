import { IAuthorization } from "./types"

export class AuthAPI {
  constructor(private _baseURL: string) {
  }

  public async list(): Promise<IAuthorization[]> {
    const res = await fetch(this._baseURL + "/authorizations")
    const auths = await res.json()
    return auths || []
  }

  public async accept(id: string): Promise<IAuthorization> {
    const res = await fetch(`${this._baseURL}/authorizations/${id}/accept`, {
      method: "POST",
    })

    return res.json()
  }

  public async deny(id: string): Promise<IAuthorization> {
    const res = await fetch(`${this._baseURL}/authorizations/${id}/deny`, {
      method: "POST",
    })

    return res.json()
  }

  public notifyEvents(cb: (err?: any, data?: any) => void): void {
    const ws = new WebSocket(`${this._baseURL.replace("http", "ws")}/events`)

    ws.onclose = (e) => {
      cb(e)
    }

    ws.onerror = (e) => {
      cb(e)
    }

    ws.onmessage = (e) => {
      cb(null, e.data)
    }
  }
}
