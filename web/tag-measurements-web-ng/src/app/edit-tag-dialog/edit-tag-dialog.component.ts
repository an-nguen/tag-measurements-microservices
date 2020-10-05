import {Component, Inject, OnInit} from '@angular/core';
import * as moment from 'moment';
import {Tag} from "../_domains/tag";
import {TagService} from "../_services/tag.service";
import {MatSnackBar} from "@angular/material/snack-bar";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from "@angular/material/core";
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from "@angular/material-moment-adapter";

interface EditTagDialogData {
  tag?: Tag
}

@Component({
  selector: 'app-edit-tag-dialog',
  templateUrl: './edit-tag-dialog.component.html',
  styleUrls: ['./edit-tag-dialog.component.css'],
  providers: [
    {provide: DateAdapter, useClass: MomentDateAdapter, deps: [MAT_DATE_LOCALE]},
    {provide: MAT_DATE_FORMATS, useValue: MAT_MOMENT_DATE_FORMATS},
  ],
})
export class EditTagDialogComponent implements OnInit {
  verificationDate?: moment.Moment;

  constructor(private tagService: TagService,
              private snackBar: MatSnackBar,
              public dialogRef: MatDialogRef<EditTagDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: EditTagDialogData) { }

  ngOnInit(): void {
    if (!!this.data.tag.verification_date) {
      this.verificationDate = this.data.tag.verification_date;
    }
  }

  editTag() {
    if (this.verificationDate.isSame(moment(this.data.tag.verification_date))) {
      this.data.tag.verification_date = this.verificationDate;
    } else {
      this.data.tag.verification_date = this.verificationDate;
      this.data.tag.verification_date.add(1, 'd');
    }
    this.tagService.updateTag(this.data.tag).subscribe((result: Tag) => {
      if (!!result && !!result.uuid) {
        this.snackBar.open(`Тег с id (${result.uuid}) изменён.`, 'Закрыть', {
          duration: 5000
        });
        this.data.tag.verification_date.add(-1, 'd');
      }
      this.dialogRef.close();
    }, error => {
      this.snackBar.open(`Тег изменить не удалось: ${error}`, 'Закрыть', {
        duration: 5000
      })
      this.dialogRef.close();
    });
  }


}
