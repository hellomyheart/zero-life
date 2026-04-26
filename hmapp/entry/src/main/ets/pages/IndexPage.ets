import { Transaction } from '../model/Transaction';
import { TransactionService } from '../service/TransactionService';
import { StatisticsService } from '../service/StatisticsService';
import { TransactionCard } from '../view/TransactionCard';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct IndexPage {
  @State transactions: Transaction[] = []
  @State total: number = 0

  build() {
    Navigation() {
      Column() {
        Row() {
          Column() {
            Text('today')
              .fontSize(18)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
          }
          .alignItems(HorizontalAlign.Start)

          Column() {
            Text('total')
              .fontSize(18)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
            Text(this.total.toString())
              .fontSize(18)
              .fontWeight(FontWeight.Medium)
              .fontColor($r('app.color.text_primary'))
          }
          .alignItems(HorizontalAlign.End)
        }
        .width('92%')
        .justifyContent(FlexAlign.SpaceBetween)
        .padding(15)
        .backgroundColor($r('app.color.card_background'))
        .margin({ top: 15, left: 15, right: 15 })

        List() {
          ForEach(this.transactions, (item: Transaction, index: number) => {
            TransactionCard({ transaction: item })
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
    TransactionService.listTransactions('').then((data: Transaction[]) => {
      this.transactions = data
      this.total = StatisticsService.getDailyStatistics(data)
    })
  }
}