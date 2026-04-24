/**
 * 多因素认证（MFA/2FA）相关类型定义
 * 用于双因素认证的状态查询、设置和验证
 */

/**
 * MFA状态接口
 * 查询当前用户是否已启用双因素认证
 */
export interface MFAStatus {
  enabled: boolean             // 是否已启用MFA
}

/**
 * MFA设置响应接口
 * 初始化MFA设置时返回的密钥和二维码信息
 */
export interface MFASetupResp {
  secret: string               // MFA密钥，用于手动输入到验证器App中
  qr_code_url: string          // 二维码URL，验证器App扫描此二维码即可绑定
}

/**
 * MFA验证请求接口
 * 提交验证器App生成的动态验证码
 */
export interface MFAVerifyReq {
  code: string                 // 6位数字验证码，由验证器App生成
}
