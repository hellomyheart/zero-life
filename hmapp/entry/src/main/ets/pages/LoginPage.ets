import { UserService } from '../common/service/UserService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct LoginPage {
  @State email: string = ''
  @State phone: string = ''
  @State password: string = ''

  build() {
    Navigation() {
      Column() {
        Column() {
          TextInput({ placeholder: $r('app.string.enter_email') })
            .placeholderFont({ size: 14, weight: 400 })
            .caretColor(Color.Blue)
            .width(150)
            .margin(20)
            .onChange((value: string) => {
              this.email = value
            })
          TextInput({ placeholder: $r('app.string.enter_phone') })
            .placeholderFont({ size: 14, weight: 400 })
            .caretColor(Color.Blue)
            .width(150)
            .margin(20)
            .onChange((value: string) => {
              this.phone = value
            })
          TextInput({ placeholder: $r('app.string.enter_password') })
            .placeholderFont({ size: 14, weight: 400 })
            .caretColor(Color.Blue)
            .width(150)
            .margin(20)
            .type(InputType.Password)
            .onChange((value: string) => {
              this.password = value
            })
        }
        .padding(20)
        .justifyContent(FlexAlign.End)
        .alignItems(HorizontalAlign.Center)

        Row() {
          Button() {
            Text($r('app.string.login'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            UserService.login(this.email, this.phone, this.password).then((user) => {
              router.pushUrl({ url: 'pages/Index' })
            })
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
      .height('100%')
    }
    .titleMode(NavigationTitleMode.Mini)
  }
}