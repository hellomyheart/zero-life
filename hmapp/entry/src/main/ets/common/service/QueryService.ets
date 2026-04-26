import { Transaction } from '../model/Transaction'

export class QueryService {
  static queryByDate(transactions: Transaction[], date: string): Transaction[] {
    return transactions.filter((transaction: Transaction) => {
      return transaction.date === date
    })
  }

  static queryByCategory(transactions: Transaction[], categoryId: string): Transaction[] {
    return transactions.filter((transaction: Transaction) => {
      return transaction.category_id === categoryId
    })
  }

  static queryByAccount(transactions: Transaction[], account: string): Transaction[] {
    return transactions.filter((transaction: Transaction) => {
      return transaction.account === account
    })
  }
}