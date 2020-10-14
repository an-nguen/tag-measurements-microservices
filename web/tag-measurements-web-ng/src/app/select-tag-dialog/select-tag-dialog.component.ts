import {Component, Inject, OnInit, ViewChild} from '@angular/core';
import {Tag} from '../_domains/tag';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {PlotService} from '../_services/plot.service';
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
import {ErrorNotifyService} from '../_services/error-notify.service';
import {LoadingService} from '../_services/loading.service';
import {MatSort} from '@angular/material/sort';
import {MatTableDataSource} from '@angular/material/table';
import {MatInput} from "@angular/material/input";
import {MatButtonToggleGroup} from "@angular/material/button-toggle";
import * as moment from "moment";

export interface SelectTagDialogData {
  tags: Tag[];
  // available values: temperature, humidity
  type: string;
}

@Component({
  selector: 'app-select-tag-dialog',
  templateUrl: './select-tag-dialog.component.html',
  providers: [
    // `MomentDateAdapter` and `MAT_MOMENT_DATE_FORMATS` can be automatically provided by importing
    // `MatMomentDateModule` in your applications root module. We provide it at the component level
    // here, due to limitations of our example generation script.
    {provide: MAT_DATE_LOCALE, useValue: 'ru-RU'},
    {provide: DateAdapter, useClass: MomentDateAdapter, deps: [MAT_DATE_LOCALE]},
    {provide: MAT_DATE_FORMATS, useValue: MAT_MOMENT_DATE_FORMATS},
  ],
  styleUrls: ['./select-tag-dialog.component.css']
})
export class SelectTagDialogComponent implements OnInit {
  @ViewChild(MatSort, {static: false}) sort: MatSort;
  @ViewChild(MatInput, {static: false}) filterTextField: MatInput;

  public dataSource: MatTableDataSource<Tag>;
  public displayedColumns = ['select', 'name', 'tagNumber', 'uuid'];

  constructor(
    public plotService: PlotService,
    public dialogRef: MatDialogRef<SelectTagDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: SelectTagDialogData,
  ) { }

  ngOnInit(): void {
    this.dataSource = new MatTableDataSource<Tag>(this.data.tags);
    this.dataSource.sort = this.sort;
  }

  async onBuildClick(): Promise<void> {
    const param = {
      tags: this.data.tags,
      dataType: this.data.type,
      action: 'goto',
    };
    this.plotService.build(param);
    this.dialogRef.close();
  }

  onCancelClick(): void {
    this.dialogRef.close();
  }

  isAllSelected() {
    const numSelected = this.plotService.tagSelection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  masterToggle() {
    this.isAllSelected() ?
        this.plotService.tagSelection.clear() :
        this.dataSource.data.forEach(row => this.plotService.tagSelection.select(row));
  }

  applyFilter($event: KeyboardEvent) {
    const filterValue = this.filterTextField.value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
  }
}
