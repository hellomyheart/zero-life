@Observed
export class State {
  hoverColor: ResourceColor = $r('sys.color.ohos_id_color_hover');

  activeColor: ResourceColor = $r('sys.color.ohos_id_color_active');

  pressColor: ResourceColor = $r('sys.color.ohos_id_color_click');

  focusBorderColor: ResourceColor = $r('sys.color.ohos_id_color_focused_outline');

  focusBorderWidth: number = $r('sys.ohos_id_border_width_focused_outlined');

  focusOutline: FocusOutlineStyle = FocusOutlineStyle.DEFAULT;

  constructor(activeColor?: ResourceColor, pressColor?: ResourceColor) {
    if (activeColor !== undefined) {
      this.activeColor = activeColor;
    }
    if (pressColor !== undefined) {
      this.pressColor = pressColor;
    }
  }
}