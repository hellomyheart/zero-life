import { BackendAdapter } from '../../data/BackendAdapter'
import { Logger } from '../../utils/Logger'
import { User } from '../model/User'

export class UserService {
  static async register(email: string, phone: string, password: string, nickname: string, avatar: string): Promise<User> {
    let user = new User('', nickname, email, phone, password, avatar)
    let params = {
      nickname: nickname,
      email: email,
      phone: phone,
      password: password,
      avatar: avatar
    }
    return new Promise<User>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/register', params, (err, data) => {
        Logger.debug('register result: ' + JSON.stringify(data))
        resolve(data as User)
      })
    })
  }

  static async login(email: string, phone: string, password: string): Promise<User> {
    let params = {
      email: email,
      phone: phone,
      password: password
    }
    return new Promise<User>((resolve) => {
      let adapter = new BackendAdapter('http://8.142.221.107:8080')
      adapter.post('/api/login', params, (err, data) => {
        Logger.debug('login result: ' + JSON.stringify(data))
        resolve(data as User)
      })
    })
  }
}