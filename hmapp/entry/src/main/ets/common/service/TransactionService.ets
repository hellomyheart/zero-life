import { BackendAdapter } from '../../data/BackendAdapter'
import { Logger } from '../../utils/Logger'
import { Transaction } from '../model/Transaction'

export class TransactionService {
  static async addTransaction(user_id: string, category_id: string, amount: number, transaction_type: string, account: string, date: string): Promise<Transaction> {
    let params = {
      user_id: user_id,
      category_id: category_id,
      amount: amount,
      transaction_type: transaction_type,
      account: account,
      date: date
    }
    return new Promise<Transaction>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/transaction', params, (err, data) => {
        Logger.debug('addTransaction result: ' + JSON.stringify(data))
        resolve(data as Transaction)
      })
    })
  }

  static async listTransactions(user_id: string): Promise<Transaction[]> {
    let params = {
      user_id: user_id
    }
    return new Promise<Transaction[]>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.get('/api/transaction/list', params, (err, data) => {
        Logger.debug('listTransactions result: ' + JSON.stringify(data))
        resolve(data as Transaction[])
      })
    })
  }
}