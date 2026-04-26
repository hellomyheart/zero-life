import { BackendAdapter } from '../data/BackendAdapter'
import { Logger } from '../utils/Logger'
import { Config } from '../common/model/Config'
import http from '@ohos.net.http'

export class ConfigService {
  private static config: Config = new Config('base_url', '')

  public static setBaseURL(url: string) {
    this.config.value = url
  }

  public static getBaseURL(): string {
    return this.config.value
  }

  public static validateBaseURL(url: string): Promise<http.HttpResponse> {
    return new Promise<http.HttpResponse>((resolve) => {
      try {
        let httpRequest = http.createHttp()
        httpRequest.request(url + '/health', {
          method: http.RequestMethod.GET,
          connectTimeout: 60000,
          readTimeout: 60000,
          header: {
            'Content-Type': 'application/json'
          },
          expectDataType: http.HttpDataType.STRING
        }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
          if (!err) {
            Logger.debug('request failed, code is ${err.code}, message is ${err.message}')
            resolve(data)
          } else {
            Logger.error('request success, status code is ${data.responseCode}')
            resolve(undefined)
          }
          httpRequest.destroy()
        })
      } catch (error) {
        Logger.error(`An unexpected error occurred. Code: ${error.code}, message: ${error.message}`)
        resolve(undefined)
      }
    })
  }

  public static healthCheck(baseUrl: string): void {
    let url = baseUrl + '/health'
    this.validateBaseURL(url).then((res) => {
      Logger.debug('healthCheck: ' + JSON.stringify(res))
      if (res.responseCode === 200) {
        Logger.debug('healthCheck success')
      }
    })
  }
}