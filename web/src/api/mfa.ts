/**
 * 多因素认证（MFA/2FA）API接口
 * 提供双因素认证的设置、启用、禁用和验证功能
 * 启用MFA后，用户登录时除了密码还需要输入验证器App生成的动态验证码，提高账户安全性
 */
import { get, post } from '@/utils/request'
import type { MFAStatus, MFASetupResp, MFAVerifyReq, BackupCodesResp } from '@/types/mfa'

/**
 * 获取当前用户的MFA状态
 * @returns MFA状态信息（是否已启用）
 * @endpoint GET /mfa/status
 */
export function status() { return get<MFAStatus>('/mfa/status') }

/**
 * 初始化MFA设置
 * 生成MFA密钥和二维码URL，用户需用验证器App（如Google Authenticator）扫描二维码
 * @returns MFA设置响应（包含密钥和二维码URL）
 * @endpoint POST /mfa/setup
 */
export function setup() { return post<MFASetupResp>('/mfa/setup') }

/**
 * 启用MFA
 * 需先调用 setup() 获取密钥并绑定验证器，然后提交验证码完成启用
 * @param data - 验证请求参数（验证器App生成的动态验证码）
 * @endpoint POST /mfa/enable
 */
export function enable(data: MFAVerifyReq) { return post<BackupCodesResp>('/mfa/enable', data) }

/**
 * 禁用MFA
 * 需要提供当前验证码以确认操作
 * @param data - 验证请求参数（验证器App生成的动态验证码）
 * @endpoint POST /mfa/disable
 */
export function disable(data: MFAVerifyReq) { return post<void>('/mfa/disable', data) }

/**
 * 验证MFA验证码
 * 在登录等需要二次验证的场景中使用
 * @param data - 验证请求参数（验证器App生成的动态验证码）
 * @endpoint POST /mfa/verify
 */
export function verify(data: MFAVerifyReq) { return post<void>('/mfa/verify', data) }

/**
 * 重新生成备用码
 * 旧的备用码将全部失效，返回新的备用码列表（仅此一次展示）
 * @returns 备用码列表
 * @endpoint POST /mfa/backup-codes
 */
export function regenerateBackupCodes() { return post<BackupCodesResp>('/mfa/backup-codes') }
