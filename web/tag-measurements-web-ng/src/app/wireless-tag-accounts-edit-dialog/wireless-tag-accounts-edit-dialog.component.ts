import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {Tag} from "../_domains/tag";
import {WirelessTagAccount} from "../_domains/wirelessTagAccount";
import {WstAccountService} from "../_services/wst-account.service";
import {MatSnackBar} from "@angular/material/snack-bar";

export interface WSTEditDataInterface {
  selectedWstAccount: WirelessTagAccount;
  mode: 'create' | 'edit' | undefined;
}

@Component({
  selector: 'app-wireless-tag-accounts-edit-dialog',
  templateUrl: './wireless-tag-accounts-edit-dialog.component.html',
  styleUrls: ['./wireless-tag-accounts-edit-dialog.component.css']
})
export class WirelessTagAccountsEditDialogComponent implements OnInit {
  emailValue: string = '';
  passwordValue: string = '';

  constructor(public dialogRef: MatDialogRef<WirelessTagAccountsEditDialogComponent>,
              public wirelessTagAccountService: WstAccountService,
              public snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: WSTEditDataInterface) { }

  ngOnInit(): void {
  }

  addWirelessTagAccount() {
    const account = {
      email: this.emailValue,
      password: this.passwordValue,
    }
    this.wirelessTagAccountService.addWstAccount(account)
      .subscribe((resp: WirelessTagAccount) => {
        this.snackBar.open(`Аккаунт ${resp.email} добавлен`, 'Закрыть',{
          duration: 5000
        })
        this.dialogRef.close(resp);
      }, error => {
        this.snackBar.open(`Аккаунт не добавлен. Ошибка: ${error.error}`, 'Закрыть',{
          duration: 5000
        })
        this.dialogRef.close();
      });
  }

  editWirelessTagAccount() {
    const account = {
      email: this.emailValue,
      password: this.passwordValue,
    }
    this.wirelessTagAccountService.updateWstAccount(account)
      .subscribe((resp: WirelessTagAccount) => {

      })
  }
}
