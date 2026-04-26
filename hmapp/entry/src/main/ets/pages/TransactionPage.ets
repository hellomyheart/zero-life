import { Category } from '../model/Category';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct TransactionPage {
  @State category_id: string = ''
  @State amount: string = ''
  @State transaction_type: string = ''
  @State account: string = ''
  @State date: Date = new Date()

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.record'))
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
            Text($r('app.string.amount'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.amount })
              .width(150)
              .placeholderFont({ size: 14 })
              .onChange((value: string) => {
                this.amount = value
              })
          }
          .width('100%')

          Row() {
            Text($r('app.string.transaction_type'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            Select([
              { value: 'expense' },
              { value: 'income' }
            ])
              .selected(0)
              .value('expense')
              .font({ size: 16 })
              .fontColor($r('app.color.text_primary'))
              .optionFont({ size: 16, weight: 400 })
              .space(8)
              .arrowPosition(ArrowPosition.START)
              .menuAlign(MenuAlignType.START, {dx:0, dy:0})
              .optionWidth(200)
              .optionHeight(300)
              .width(150)
          }
          .width('100%')
          .margin({ top: 10 })

          Row() {
            Text($r('app.string.account'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            TextInput({ text: this.account })
              .width(150)
              .placeholderFont({ size: 14 })
              .onChange((value: string) => {
                this.account = value
              })
          }
          .width('100%')
          .margin({ top: 10 })

          Row() {
            Text($r('app.string.date'))
              .fontSize(16)
              .width(110)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            DatePicker({
              start: new Date('1970-1-1'),
              end: new Date('2100-12-31'),
              selected: this.date
            })
              .disappearTextStyle({color: Color.Red, font: {size: '16px', weight: FontWeight.Bold}})
              .textStyle({color: '#FF182431', font: {size: '16px', weight: FontWeight.Regular}})
              .selectedTextStyle({color: '#FF007DFF', font: {size: '16px', weight: FontWeight.Regular}})
              .onDateChange((value: Date) => {
                this.date = value
              })
              .width(150)
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
            router.back()
          })
          .width(180)
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
}