export interface IJSONRPCRequest {
  id: any
  method: string
  params: any[]
  auth: string
}

export interface IAuthorization {
  id: string
  state: "pending" | "accepted" | "denied" | "consumed"
  request: IJSONRPCRequest
  createdAt: string
}
