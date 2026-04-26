import { BackendAdapter } from '../../data/BackendAdapter'
import { Logger } from '../../utils/Logger'
import { Account } from '../model/Account'

export class AccountService {
  static async addAccount(name: string): Promise<Account> {
    let params = {
      name: name
    }
    return new Promise<Account>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/account', params, (err, data) => {
        Logger.debug('addAccount result: ' + JSON.stringify(data))
        resolve(data as Account)
      })
    })
  }

  static async listAccounts(): Promise<Account[]> {
    return new Promise<Account[]>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.get('/api/account/list', {}, (err, data) => {
        Logger.debug('listAccounts result: ' + JSON.stringify(data))
        resolve(data as Account[])
      })
    })
  }
}