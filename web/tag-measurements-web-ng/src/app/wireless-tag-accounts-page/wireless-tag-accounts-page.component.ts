import { Component, OnInit } from '@angular/core';
import {WstAccountService} from '../_services/wst-account.service';
import {WirelessTagAccount} from '../_domains/wirelessTagAccount';
import {ErrorNotifyService} from '../_services/error-notify.service';
import {MatDialog} from "@angular/material/dialog";
import {WirelessTagAccountsEditDialogComponent} from "../wireless-tag-accounts-edit-dialog/wireless-tag-accounts-edit-dialog.component";
import {MatSnackBar} from "@angular/material/snack-bar";

@Component({
  selector: 'app-wst-accounts-page',
  templateUrl: './wireless-tag-accounts-page.component.html',
  styleUrls: ['./wireless-tag-accounts-page.component.css']
})
export class WirelessTagAccountsPageComponent implements OnInit {
  public wstAccountList: WirelessTagAccount[];
  public selectedWstAccount: Array<WirelessTagAccount>;

  constructor(public wstAccountService: WstAccountService,
              public errorNotifyService: ErrorNotifyService,
              public dialog: MatDialog,
              public snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.wstAccountService.getWstAccounts()
        .subscribe((result) => {
          if (!result) {
            this.errorNotifyService.callErrorDialog('Нет WST аккаунтов либо нет связи с сервером ресурсов.');
          }
          this.wstAccountList = result;
        }, error => {
          this.errorNotifyService.callErrorDialog('Неизвестная ошибка - ' + error.error);
        });
  }

  create() {
    this.dialog.open(WirelessTagAccountsEditDialogComponent, {
      data: {selectedWstAccount: {email: '', password: ''}, mode: 'create'}
    })
  }

  edit(account: WirelessTagAccount) {
    if (account) {
      this.dialog.open(WirelessTagAccountsEditDialogComponent, {
        data: {selectedWstAccount: account, mode: 'edit'}
      }).afterClosed().subscribe((resp: WirelessTagAccount) => {
        if (resp) {
          this.snackBar
        }
      })
    }
  }

}
