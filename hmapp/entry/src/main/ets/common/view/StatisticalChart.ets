@Component
export struct StatisticalChart {
  build() {
    Column() {
      Text('statistical')
        .fontSize(18)
        .fontWeight(FontWeight.Medium)
        .fontColor($r('app.color.text_primary'))
    }
  }
}