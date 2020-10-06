import {Component, OnInit, ViewChild} from '@angular/core';
import {MatDialog} from '@angular/material/dialog';
import {SelectTagDialogComponent} from '../select-tag-dialog/select-tag-dialog.component';
import {TagManagerListService} from '../_services/tag-manager-list.service';
import {MatSort} from '@angular/material/sort';
import {MatTableDataSource} from "@angular/material/table";
import {Tag} from "../_domains/tag";
import {MatInput} from "@angular/material/input";
import {PlotService} from "../_services/plot.service";
import {EditTagDialogComponent} from "../edit-tag-dialog/edit-tag-dialog.component";
import {interval} from "rxjs";
import {flatMap} from "rxjs/operators";
import * as moment from 'moment';
import {FormControl} from "@angular/forms";
import {RoleService} from "../_services/role.service";

@Component({
  selector: 'app-tag-manager-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.css']
})
export class TagListComponent implements OnInit {
  @ViewChild(MatSort, {static: false}) sort: MatSort;
  @ViewChild(MatInput, {static: false}) filterTextField: MatInput;
  public dataSource: MatTableDataSource<Tag>;
  public displayedColumns = ['name', 'tagNumber', 'uuid', 'verification_date', 'temperature', 'humidity', 'voltage', 'batteryRemaining', 'signaldBm', 'actions'];

  constructor(public tagManagerListService: TagManagerListService,
              public plotService: PlotService,
              public roleService: RoleService,
              public dialog: MatDialog
  ) { }

  private isLoadingWG = false;
  warehouseGroupControl = new FormControl();

  ngOnInit(): void {
    interval(1000*10)
      .pipe(
        flatMap(() => this.updateTagsDetails())
      )
      .subscribe(data => {

      });
    this.tagManagerListService.refreshTemperatureZones();
    this.roleService.getRoleByToken();
  }

  onTemperatureZoneSelectChange(id: number) {
    this.isLoadingWG = true;
    this.tagManagerListService.selectTags(id)
      .add(() => {
        this.dataSource = new MatTableDataSource<Tag>(this.tagManagerListService.tags);
        this.dataSource.sort = this.sort;
        this.isLoadingWG = false;
        this.updateTagsDetails();
      });
    this.plotService.tagSelection.clear();
  }

  get getLoadingWG() {
    return this.isLoadingWG;
  }

  get getWGSelected() {
    return this.tagManagerListService.temperatureZones.find(wg => wg.id === this.warehouseGroupControl.value);
  }

  isBetweenTemperatureLimit(tag: Tag) {
    return this.getWGSelected.lower_temp_limit < tag.temperature && this.getWGSelected.higher_temp_limit > tag.temperature
  }

  isWarnVerificationDate(tag: Tag) {
    if (!!tag.verification_date) {
      const endLimit = moment(tag.verification_date).add(2, 'y');
      const beginLimit = moment(tag.verification_date).add(2, 'y').subtract(2, 'week');
      return moment().isBetween(beginLimit, endLimit);
    } else {
      return false;
    }
  }

  isAfterTwoYears(tag: Tag) {
    return moment().isAfter(moment(tag.verification_date).add(2, 'y'));
  }

  openSelectTagDialog(type: string) {
    const dialogRef = this.dialog.open(SelectTagDialogComponent, {
      width: '70%',
      data: {tags: this.tagManagerListService.tags, type}
    });

    dialogRef.afterClosed().subscribe(result => {
    });
  }

  applyFilter($event: KeyboardEvent) {
    const filterValue = this.filterTextField.value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
  }

  editTag(element: Tag) {
    this.dialog.open(EditTagDialogComponent, {
      width: '400px',
      data: {tag: element}
    });
  }

  private updateTagsDetails() {
    if (!this.dataSource) {
      return [];
    }

    this.tagManagerListService.getLatestMeasurement().subscribe((res: any[]) => {
      for (let tag of this.dataSource.data) {
        let found = res.find(t => t.uuid === tag.uuid);
        if (!!found) {
          tag.cap = found.cap;
          tag.batteryVolt = found.batteryVolt;
          tag.signaldBm = found.signaldBm;
          tag.batteryRemaining = found.batteryRemaining;
          tag.alive = found.alive;
          tag.temperature = found.temperature;
          tag.lux = found.lux;
        }
      }
      res.splice(0, res.length);
    });

    return this.dataSource.data;
  }

  printDate(element: Tag) {
    return moment(element.verification_date).format("DD.MM.YYYY").toString();
  }
}
