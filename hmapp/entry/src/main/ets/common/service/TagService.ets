import { BackendAdapter } from '../../data/BackendAdapter'
import { Logger } from '../../utils/Logger'
import { Tag } from '../model/Tag'

export class TagService {
  static async addTag(name: string): Promise<Tag> {
    let params = {
      name: name
    }
    return new Promise<Tag>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/tag', params, (err, data) => {
        Logger.debug('addTag result: ' + JSON.stringify(data))
        resolve(data as Tag)
      })
    })
  }

  static async listTags(): Promise<Tag[]> {
    return new Promise<Tag[]>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.get('/api/tag/list', {}, (err, data) => {
        Logger.debug('listTags result: ' + JSON.stringify(data))
        resolve(data as Tag[])
      })
    })
  }
}