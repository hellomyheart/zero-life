import { Transaction } from '../model/Transaction'

export class StatisticsService {
  static getDailyStatistics(transactions: Transaction[]): number {
    let total: number = 0
    for (let i: number = 0; i < transactions.length; i++) {
      total += transactions[i].amount
    }
    return total
  }
}