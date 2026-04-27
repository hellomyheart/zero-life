/**
 * 多因素认证（MFA/2FA）相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * MFA状态接口（对应后端 MFAStatusResp）
 */
export interface MFAStatus {
  enabled: boolean
}

/**
 * MFA设置响应接口（对应后端 MFASetupResp）
 * 后端字段名为 qr_code（非 qr_code_url）
 */
export interface MFASetupResp {
  secret: string
  qr_code: string
}

/**
 * MFA验证请求接口（对应后端 MFAVerifyReq）
 */
export interface MFAVerifyReq {
  code: string
}

/**
 * MFA验证响应接口（对应后端 MFAVerifyResp）
 */
export interface MFAVerifyResp {
  verified: boolean
  message: string
}

export interface BackupCodesResp {
  codes: string[]
}
