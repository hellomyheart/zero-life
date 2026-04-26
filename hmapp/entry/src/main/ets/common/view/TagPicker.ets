import { Tag } from '../model/Tag'
import { TagService } from '../service/TagService'

@Component
export struct TagPicker {
  @Link selectedTags: Tag[]
  @Prop tags: Tag[]

  build() {
    Select(this.tags.map((tag: Tag) => {
      return { value: tag.name }
    }))
      .selected(1)
      .value(this.selectedTags[0].name)
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

  aboutToAppear(): void {
    TagService.listTags().then((data: Tag[]) => {
      this.tags = data
    })
  }
}