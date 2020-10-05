import { Injectable } from '@angular/core';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {ErrorDialogComponent} from '../error-dialog/error-dialog.component';

@Injectable({
  providedIn: 'root'
})
export class ErrorNotifyService {

  constructor(
    private errorDialog: MatDialog
  ) { }

  callErrorDialog(content: string) {
    const dialogRef = this.errorDialog.open(ErrorDialogComponent, {
      width: '70%',
      data: {errorContent: content}
    });
  }
}
