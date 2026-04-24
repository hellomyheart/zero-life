export interface MFAStatus {
  enabled: boolean
}

export interface MFASetupResp {
  secret: string
  qr_code_url: string
}

export interface MFAVerifyReq {
  code: string
}
