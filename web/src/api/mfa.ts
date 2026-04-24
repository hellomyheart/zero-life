import { get, post } from '@/utils/request'
import type { MFAStatus, MFASetupResp, MFAVerifyReq } from '@/types/mfa'

export function status() { return get<MFAStatus>('/mfa/status') }
export function setup() { return post<MFASetupResp>('/mfa/setup') }
export function enable(data: MFAVerifyReq) { return post<void>('/mfa/enable', data) }
export function disable(data: MFAVerifyReq) { return post<void>('/mfa/disable', data) }
export function verify(data: MFAVerifyReq) { return post<void>('/mfa/verify', data) }
