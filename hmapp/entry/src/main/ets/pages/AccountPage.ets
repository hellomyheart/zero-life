import { Account } from '../model/Account';
import { AccountService } from '../service/AccountService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct AccountPage {
  @State accounts: Account[] = []
  @State name: string = ''

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.account'))
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
            Text($r('app.string.name'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.name })
              .width(150)
              .placeholderFont({ size: 14 })
              .onChange((value: string) => {
                this.name = value
              })
          }
          .width('100%')

          Button() {
            Text($r('app.string.save'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            AccountService.addAccount(this.name).then((res) => {
              AccountService.listAccounts().then((data) => {
                this.accounts = data
              })
            })
          })
          .width(180)
          .height(45)
          .margin(20)
          .align(Alignment.Center)
        }
        .width('92%')
        .alignItems(HorizontalAlign.Start)
        .margin({ top: 10, left: 15, right: 15 })

        List() {
          ForEach(this.accounts, (item: Account, index: number) => {
            ListItem() {
              Column() {
                Text(item.name)
                  .fontSize(16)
                  .fontWeight(FontWeight.Medium)
                  .fontColor($r('app.color.text_primary'))
                  .margin({ left: 15 })
              }
              .alignItems(HorizontalAlign.Start)
              .width('100%')
            }
          })
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
    AccountService.listAccounts().then((data) => {
      this.accounts = data
    })
  }
}