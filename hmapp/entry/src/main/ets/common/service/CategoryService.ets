import { BackendAdapter } from '../../data/BackendAdapter'
import { Logger } from '../../utils/Logger'
import { Category } from '../model/Category'

export class CategoryService {
  static async addCategory(name: string, transaction_type: string): Promise<Category> {
    let params = {
      name: name,
      transaction_type: transaction_type
    }
    return new Promise<Category>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/category', params, (err, data) => {
        Logger.debug('addCategory result: ' + JSON.stringify(data))
        resolve(data as Category)
      })
    })
  }

  static async listCategories(): Promise<Category[]> {
    return new Promise<Category[]>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.get('/api/category/list', {}, (err, data) => {
        Logger.debug('listCategories result: ' + JSON.stringify(data))
        resolve(data as Category[])
      })
    })
  }
}