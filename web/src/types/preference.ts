/**
 * 用户偏好相关类型定义
 * 对应后端 model/preference.go 和 DTO
 */

/**
 * 用户偏好数据类型
 * 存储用户个性化配置（键值对）
 */
export interface Preference {
  /** 偏好键名 */
  key: string;
  /** 偏好值 */
  value: string;
}

/**
 * 设置偏好请求参数
 */
export interface SetPreferenceRequest {
  /** 偏好键名 */
  key: string;
  /** 偏好值 */
  value: string;
}
