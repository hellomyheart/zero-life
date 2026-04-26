import http from '@ohos.net.http';
import { BusinessError } from '@kit.BasicServicesKit';
import { Logger } from '../utils/Logger';

export class BackendAdapter {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  public get(path: string, params: Record<string, any>, callback: (err?: BusinessError, data?: ESObject) => void) {
    let httpRequest = http.createHttp()
    let url = this.baseUrl + path

    httpRequest.request(url, {
      method: http.RequestMethod.GET,
      extraData: params,
      connectTimeout: 60000,
      readTimeout: 60000,
      header: {
        'Content-Type': 'application/json'
      },
      expectDataType: http.HttpDataType.STRING
    }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
      if (err) {
        Logger.error(`request failed, code is ${err.code}, message is ${err.message}`);
      } else {
        let result = JSON.parse(data.result as string)
        callback(undefined, result)
      }
    })
  }

  public post(path: string, params: Record<string, any>, callback: (err?: BusinessError, data?: ESObject) => void) {
    let httpRequest = http.createHttp()
    let url = this.baseUrl + path

    httpRequest.request(url, {
      method: http.RequestMethod.POST,
      extraData: params,
      connectTimeout: 60000,
      readTimeout: 60000,
      header: {
        'Content-Type': 'application/json'
      },
      expectDataType: http.HttpDataType.STRING
    }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
      if (err) {
        Logger.error(`request failed, code is ${err.code}, message is ${err.message}`);
      } else {
        let result = JSON.parse(data.result as string)
        callback(undefined, result)
      }
    })
  }

  public delete(path: string, params: Record<string, any>, callback: (err?: BusinessError, data?: ESObject) => void) {
    let httpRequest = http.createHttp()
    let url = this.baseUrl + path

    httpRequest.request(url, {
      method: http.RequestMethod.DELETE,
      extraData: params,
      connectTimeout: 60000,
      readTimeout: 60000,
      header: {
        'Content-Type': 'application/json'
      },
      expectDataType: http.HttpDataType.STRING
    }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
      if (err) {
        Logger.error(`request failed, code is ${err.code}, message is ${err.message}`);
      } else {
        let result = JSON.parse(data.result as string)
        callback(undefined, result)
      }
    })
  }

  public put(path: string, params: Record<string, any>, callback: (err?: BusinessError, data?: ESObject) => void) {
    let httpRequest = http.createHttp()
    let url = this.baseUrl + path

    httpRequest.request(url, {
      method: http.RequestMethod.PUT,
      extraData: params,
      connectTimeout: 60000,
      readTimeout: 60000,
      header: {
        'Content-Type': 'application/json'
      },
      expectDataType: http.HttpDataType.STRING
    }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
      if (err) {
        Logger.error(`request failed, code is ${err.code}, message is ${err.message}`);
      } else {
        let result = JSON.parse(data.result as string)
        callback(undefined, result)
      }
    })
  }
}