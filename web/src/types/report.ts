/**
 * 报表相关类型定义
 * 定义各类财务报表的数据结构，包括收支、分类、预算、净值和趋势报表
 */

/**
 * 报表查询请求接口
 * 所有报表接口共用的查询参数
 */
export interface ReportReq {
  start_date: string           // 查询起始日期（必填）
  end_date: string             // 查询结束日期（必填）
  currency?: string            // 货币代码（可选，用于多币种报表）
  account_ids?: string[]       // 账户ID列表（可选，筛选特定账户的数据）
  category_ids?: string[]      // 分类ID列表（可选，筛选特定分类的数据）
}

/**
 * 收支报表响应接口
 * 返回指定时间段内的收支汇总数据
 */
export interface IncomeExpenseResp {
  total_income: string         // 总收入
  total_expense: string        // 总支出
  net_income: string           // 净收入（收入 - 支出）
  income_by_account: {         // 按账户分组的收入明细
    account_id: string         // 账户ID
    account_name: string       // 账户名称
    amount: string             // 收入金额
  }[]
  expense_by_account: {        // 按账户分组的支出明细
    account_id: string         // 账户ID
    account_name: string       // 账户名称
    amount: string             // 支出金额
  }[]
}

/**
 * 分类项响应接口
 * 单个分类的汇总数据
 */
export interface CategoryItemResp {
  category_id: string          // 分类ID
  category_name: string        // 分类名称
  amount: string               // 该分类的金额
  percentage: number           // 占比百分比（0-100）
}

/**
 * 分类报表响应接口
 * 按分类汇总的收入和支出数据
 */
export interface CategoryReportResp {
  income_categories: CategoryItemResp[]   // 收入分类汇总列表
  expense_categories: CategoryItemResp[]  // 支出分类汇总列表
}

/**
 * 预算报表响应接口
 * 各预算的使用情况汇总
 */
export interface BudgetReportResp {
  budgets: {                   // 预算列表
    budget_id: string          // 预算ID
    budget_name: string        // 预算名称
    amount: string             // 预算金额上限
    spent: string              // 已花费金额
    usage_rate: number         // 使用率（0-1之间）
    status: string             // 预算状态
  }[]
}

/**
 * 净值报表响应接口
 * 资产、负债和净值的变化趋势
 */
export interface NetWorthResp {
  trend: {                     // 按日期的净值变化列表
    date: string               // 日期（格式：YYYY-MM 或 YYYY-MM-DD）
    net_worth: string          // 当日净值
  }[]
  total_assets: string         // 总资产
  total_liabilities: string    // 总负债
  net_worth: string            // 当前净值（总资产 - 总负债）
}

/**
 * 趋势报表响应接口
 * 收支随时间的变化趋势，用于绘制趋势图
 */
export interface TagItemResp {
  tag_id: string
  tag_name: string
  income: string
  expense: string
}

export interface TagReportResp {
  items: TagItemResp[]
}

export interface TrendResp {
  items: {
    date: string
    income: string
    expense: string
    net: string
  }[]
}
