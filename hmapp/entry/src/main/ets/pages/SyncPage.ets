import { SyncService } from '../common/service/SyncService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct SyncPage {
  @State connected: boolean = false

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.sync'))
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
            SyncService.syncToBackend()
          })
          .width(200)
          .height(45)
          .margin(20)
          .align(Alignment.Center)
        }
        .width('92%')
        .alignItems(HorizontalAlign.Start)
        .margin({ top: 10, left: 15, right: 15 })
      }
      .width('100%')
    }
    .titleMode(NavigationTitleMode.Mini)
  }

  aboutToAppear(): void {
    SyncService.syncToDistributedStore()
  }
}