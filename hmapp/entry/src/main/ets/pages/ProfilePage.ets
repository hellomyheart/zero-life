import { UserService } from '../common/service/UserService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct ProfilePage {
  @State nickname: string = ''
  @State email: string = ''
  @State phone: string = ''
  @State password: string = ''

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.user_profile'))
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
            Text($r('app.string.nickname'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.nickname })
              .width(150)
              .placeholderFont({ size: 14 })
              .onChange((value: string) => {
                this.nickname = value
              })
          }
          .width('100%')

          Row() {
            Text($r('app.string.email'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.email })
              .width(150)
              .placeholderFont({ size: 14 })
              .enabled(false)
          }
          .width('100%')
          .margin({ top: 10 })

          Row() {
            Text($r('app.string.phone'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.phone })
              .width(150)
              .placeholderFont({ size: 14 })
              .enabled(false)
          }
          .width('100%')
          .margin({ top: 10 })

          Row() {
            Text($r('app.string.password'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.password })
              .width(150)
              .placeholderFont({ size: 14 })
              .type(InputType.Password)
              .onChange((value: string) => {
                this.password = value
              })
          }
          .width('100%')
          .margin({ top: 10 })

          Button() {
            Text($r('app.string.save'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            router.pushUrl({ url: 'pages/Index' })
          })
          .width(180)
          .height(45)
          .margin(20)
          .align(Alignment.Center)
        }
        .width('92%')
        .alignItems(HorizontalAlign.Start)
        .margin({ top: 10, left: 15, right: 15 })

        Row() {
          Button() {
            Text($r('app.string.logout'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            router.pushUrl({ url: 'pages/Login' })
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
}