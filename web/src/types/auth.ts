export interface LoginReq {
  email: string
  password: string
}

export interface RegisterReq {
  email: string
  password: string
  nickname: string
}

export interface LoginResp {
  access_token: string
  refresh_token: string
}

export interface ProfileResp {
  id: string
  email: string
  nickname: string
  language: string
  timezone: string
  created_at: string
  updated_at: string
}

export interface RefreshReq {
  refresh_token: string
}

export interface ForgotPasswordReq {
  email: string
}

export interface ResetPasswordReq {
  token: string
  new_password: string
}

export interface UpdateProfileReq {
  nickname?: string
  language?: string
  timezone?: string
}

export interface ChangePasswordReq {
  old_password: string
  new_password: string
}
