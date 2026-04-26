import { Account } from '../model/Account'
import { AccountService } from '../service/AccountService'

@Component
export struct AccountPicker {
  @Link selectedAccount: Account
  @Prop accounts: Account[]

  build() {
    Select(this.accounts.map((account: Account) => {
      return { value: account.name }
    }))
      .selected(1)
      .value(this.selectedAccount.name)
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
    AccountService.listAccounts().then((data: Account[]) => {
      this.accounts = data
    })
  }
}