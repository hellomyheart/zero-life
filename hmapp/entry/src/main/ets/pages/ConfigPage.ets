import { ConfigService } from '../common/service/ConfigService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct ConfigPage {
  @State baseUrl: string = ''
  @State connected: boolean = false

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.base_url_config'))
            .fontSize(18)
            .fontWeight(FontWeight.Medium)
            .width('100%')
            .textAlign(TextAlign.Start)
            .alignSelf(ItemAlign.Start)
        }
        .width('92%')
        .padding(15)
        .backgroundColor($r('app.color.card_background'))
        .alignItems(HorizontalAlign.Start)
        .margin({ top: 15, left: 15, right: 15 })

        Column() {
          Row() {
            Text($r('app.string.base_url'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.baseUrl })
              .width(150)
              .placeholderFont({ size: 14 })
              .onChange((value: string) => {
                this.baseUrl = value
              })
          }
          .width('100%')

          Row() {
            Text($r('app.string.connected'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            Text(this.connected ? 'Connected' : 'Not Connected')
              .fontSize(16)
              .fontColor(this.connected ? $r('app.color.text_primary') : $r('app.color.text_tertiary'))
          }
          .width('100%')
          .margin({ top: 10 })

          Button() {
            Text($r('app.string.validate'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            this.checkConnectivity(this.baseUrl)
          })
          .width(200)
          .height(45)
          .margin(20)
          .align(Alignment.Center)
        }
        .width('92%')
        .alignItems(HorizontalAlign.Start)
        .margin({ top: 10, left: 15, right: 15 })

        Row() {
          Button() {
            Text($r('app.string.save'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            ConfigService.setBaseURL(this.baseUrl)
            router.pushUrl({ url: 'pages/Index' })
          })
          .width(180)
          .height(45)
          .margin(20)
          .align(Alignment.Center)
        }
        .width('100%')
        .justifyContent(FlexAlign.SpaceEvenly)
      }
      .width('100%')
    }
    .titleMode(NavigationTitleMode.Mini)
  }

  checkConnectivity(url: string) {
    try {
      let httpRequest = http.createHttp()
      httpRequest.request(url + '/health', {
        method: http.RequestMethod.GET,
        header: {
          'Content-Type': 'application/json'
        },
        expectDataType: http.HttpDataType.STRING
      }, (err: BusinessError | undefined, data: http.HttpResponse | undefined) => {
        if (!err) {
          console.info(`request success, data: ${JSON.stringify(data)}`);
          this.connected = true
        } else {
          console.error(`request failed, code: ${err.code}, message: ${err.message}`)
          this.connected = false
        }
        httpRequest.destroy()
      })
    } catch (error) {
      console.error(`An unexpected error occurred. Code: ${error.code}, message: ${error.message}`);
      this.connected = false
    }
  }
}