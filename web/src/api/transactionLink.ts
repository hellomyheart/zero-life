/**
 * 交易关联 API 接口
 * 提供交易关联的 CRUD 操作
 */
import request from '@/utils/request';
import type {
  TransactionLink,
  CreateTransactionLinkRequest,
  TransactionLinkListParams
} from '@/types/transactionLink';

/**
 * 获取交易关联列表（分页）
 * @param params 查询参数
 * @returns 交易关联列表（分页）
 */
export function list(params?: TransactionLinkListParams) {
  return request<{
    list: TransactionLink[];
    total: number;
    page: number;
    page_size: number;
  }>({
    url: '/transaction-links',
    method: 'get',
    params
  });
}

/**
 * 创建交易关联
 * @param data 创建请求参数
 * @returns 创建的交易关联
 */
export function create(data: CreateTransactionLinkRequest) {
  return request<TransactionLink>({
    url: '/transaction-links',
    method: 'post',
    data
  });
}

/**
 * 删除交易关联
 * @param id 关联记录 ID
 */
export function remove(id: number) {
  return request({
    url: `/transaction-links/${id}`,
    method: 'delete'
  });
}
