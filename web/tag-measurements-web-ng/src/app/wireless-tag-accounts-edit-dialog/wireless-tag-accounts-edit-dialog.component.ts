import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {Tag} from "../_domains/tag";
import {WirelessTagAccount} from "../_domains/wirelessTagAccount";

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

  constructor(public dialogRef: MatDialogRef<WirelessTagAccountsEditDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: WSTEditDataInterface) { }

  ngOnInit(): void {
  }

}
