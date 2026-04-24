/**
 * 对象分组 API 接口
 * 提供对象分组的 CRUD 操作
 */
import request from '@/utils/request';
import type {
  ObjectGroup,
  CreateObjectGroupRequest,
  UpdateObjectGroupRequest,
  ObjectGroupListParams
} from '@/types/objectGroup';

/**
 * 获取对象分组列表
 * @param params 查询参数
 * @returns 对象分组列表
 */
export function list(params?: ObjectGroupListParams) {
  return request<ObjectGroup[]>({
    url: '/object-groups',
    method: 'get',
    params
  });
}

/**
 * 获取单个对象分组详情
 * @param id 分组 ID
 * @returns 对象分组详情
 */
export function get(id: number) {
  return request<ObjectGroup>({
    url: `/object-groups/${id}`,
    method: 'get'
  });
}

/**
 * 创建对象分组
 * @param data 创建请求参数
 * @returns 创建的对象分组
 */
export function create(data: CreateObjectGroupRequest) {
  return request<ObjectGroup>({
    url: '/object-groups',
    method: 'post',
    data
  });
}

/**
 * 更新对象分组
 * @param id 分组 ID
 * @param data 更新请求参数
 * @returns 更新后的对象分组
 */
export function update(id: number, data: UpdateObjectGroupRequest) {
  return request<ObjectGroup>({
    url: `/object-groups/${id}`,
    method: 'put',
    data
  });
}

/**
 * 删除对象分组
 * @param id 分组 ID
 */
export function remove(id: number) {
  return request({
    url: `/object-groups/${id}`,
    method: 'delete'
  });
}
