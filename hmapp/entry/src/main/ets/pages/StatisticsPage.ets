import { Transaction } from '../model/Transaction';
import { TransactionService } from '../service/TransactionService';
import { StatisticsService } from '../service/StatisticsService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct StatisticsPage {
  @State transactions: Transaction[] = []
  @State total: number = 0

  build() {
    Navigation() {
      Column() {
        Row() {
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
            ListItem() {
              Column() {
                Text(item.date)
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
    TransactionService.listTransactions('').then((data: Transaction[]) => {
      this.transactions = data
      this.total = StatisticsService.getDailyStatistics(data)
    })
  }
}