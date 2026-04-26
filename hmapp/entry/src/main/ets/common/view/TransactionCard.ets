import { Category } from '../model/Category';
import { Transaction } from '../model/Transaction';

@Component
export struct TransactionCard {
  @Prop transaction: Transaction

  build() {
    ListItem() {
      Column() {
        Row() {
          Column() {
            Text(this.transaction.date)
              .fontSize(16)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
              .margin({ left: 15 })
          }
          .alignItems(HorizontalAlign.Start)

          Column() {
            Text(this.transaction.account)
              .fontSize(16)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
              .margin({ right: 15 })
          }
          .alignItems(HorizontalAlign.End)
        }
        .width('100%')
        .justifyContent(FlexAlign.SpaceBetween)

        Row() {
          Column() {
            Text(this.transaction.transaction_type)
              .fontSize(16)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
              .margin({ left: 15 })
          }
          .alignItems(HorizontalAlign.Start)

          Column() {
            Text(this.transaction.amount.toString())
              .fontSize(16)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
              .margin({ right: 15 })
          }
          .alignItems(HorizontalAlign.End)
        }
        .width('100%')
        .justifyContent(FlexAlign.SpaceBetween)
      }
      .alignItems(HorizontalAlign.Center)
      .backgroundColor('#FFFFFF')
      .width('92%')
      .height(100)
      .borderRadius(16)
      .margin(15)
    }
  }
}