import { Tag } from '../model/Tag';
import { TagService } from '../service/TagService';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct TagPage {
  @State tags: Tag[] = []
  @State name: string = ''

  build() {
    Navigation() {
      Column() {
        Column() {
          Text($r('app.string.tag'))
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
            TagService.addTag(this.name).then((res) => {
              TagService.listTags().then((data) => {
                this.tags = data
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
          ForEach(this.tags, (item: Tag, index: number) => {
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
    TagService.listTags().then((data) => {
      this.tags = data
    })
  }
}