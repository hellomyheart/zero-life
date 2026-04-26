import { Category } from '../model/Category'
import { CategoryService } from '../service/CategoryService'

@Component
export struct CategoryPicker {
  @Link selectedCategory: Category
  @Prop categories: Category[]

  build() {
    Select(this.categories.map((category: Category) => {
      return { value: category.name }
    }))
      .selected(1)
      .value(this.selectedCategory.name)
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
    CategoryService.listCategories().then((data: Category[]) => {
      this.categories = data
    })
  }
}