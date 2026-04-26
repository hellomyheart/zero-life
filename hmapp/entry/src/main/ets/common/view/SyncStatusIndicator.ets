@Component
export struct SyncStatusIndicator {
  @Prop connected: boolean

  build() {
    Column() {
      Text(this.connected ? 'Connected' : 'Not Connected')
        .fontSize(16)
        .fontColor(this.connected ? $r('app.color.text_primary') : $r('app.color.text_tertiary'))
    }
  }
}