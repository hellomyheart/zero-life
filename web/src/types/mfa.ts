// 多因素认证相关类型定义 - MFA状态、设置和验证响应
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
