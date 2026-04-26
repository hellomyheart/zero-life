import { Category } from '../model/Category';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct CategoryPage {
  @State categories: Category[] = []
  @State name: string = ''
  @State transaction_type: string = 'expense'

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.category'))
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

          Button() {
            Text($r('app.string.save'))
              .fontSize(18)
              .fontColor('#FFFFFF')
              .fontWeight(500)
          }
          .onClick(() => {
            CategoryService.addCategory(this.name, this.transaction_type).then((res) => {
              CategoryService.listCategories().then((data) => {
                this.categories = data
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
          ForEach(this.categories, (item: Category, index: number) => {
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
    CategoryService.listCategories().then((data) => {
      this.categories = data
    })
  }
}