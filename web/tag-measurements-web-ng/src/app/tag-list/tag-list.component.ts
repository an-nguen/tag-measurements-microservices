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
import * as moment from 'moment';
import {FormControl} from "@angular/forms";
import {RoleService} from "../_services/role.service";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-tag-manager-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.css']
})
export class TagListComponent implements OnInit {
  @ViewChild(MatSort, {static: false}) sort: MatSort;
  @ViewChild(MatInput, {static: false}) filterTextField: MatInput;
  public dataSource: MatTableDataSource<Tag>;
  public displayedColumns = ['name'];
  private websocketMessage: string;


  constructor(public tagManagerListService: TagManagerListService,
              public plotService: PlotService,
              public roleService: RoleService,
              public dialog: MatDialog
  ) { }

  private isLoadingWG = false;
  warehouseGroupControl = new FormControl();

  ngOnInit(): void {

    this.tagManagerListService.refreshTemperatureZones();
    this.roleService.getRoleByToken().add(() => {
      this.defineTableColumns();

    });
  }
  // public displayedColumns = ['name', 'tagNumber', 'uuid', 'verification_date', 'temperature', 'cap', 'batteryVolt', 'batteryRemaining', 'signaldBm', 'actions'];

  defineTableColumns() {
    const privilegeNames = this.roleService.userPrivileges.map(privilege => privilege.name);

    if (privilegeNames.includes('ALLOW_SHOW_TAG_NUMBER'))
      this.displayedColumns.push('tagNumber')
    if (privilegeNames.includes('ALLOW_SHOW_UUID'))
      this.displayedColumns.push('uuid')
    if (privilegeNames.includes('ALLOW_VERIFICATION_DATE'))
      this.displayedColumns.push('verification_date')
    if (privilegeNames.includes('ALLOW_SHOW_TEMPERATURE'))
      this.displayedColumns.push('temperature')
    if (privilegeNames.includes('ALLOW_SHOW_HUMIDITY'))
      this.displayedColumns.push('cap')
    if (privilegeNames.includes('ALLOW_SHOW_VOLTAGE'))
      this.displayedColumns.push('batteryVolt')
    if (privilegeNames.includes('ALLOW_SHOW_BATTERY_REMAINING'))
      this.displayedColumns.push('batteryRemaining')
    if (privilegeNames.includes('ALLOW_SHOW_SIGNAL'))
      this.displayedColumns.push('signaldBm')
    if (privilegeNames.includes('ALLOW_TAG_EDIT'))
      this.displayedColumns.push('actions')

    if (privilegeNames.includes('ALLOW_SHOW_TEMPERATURE') || privilegeNames.includes('ALLOW_SHOW_HUMIDITY')
      || privilegeNames.includes('ALLOW_SHOW_VOLTAGE') || privilegeNames.includes('ALLOW_SHOW_BATTERY_REMAINING')
    || privilegeNames.includes('ALLOW_SHOW_SIGNAL')) {
      const conn = new WebSocket(`ws://${environment.ws}/ws/tags`);
      conn.onopen = () => {
        setInterval(() => {
          conn.send("/latest");
        }, 5000)
      };
      conn.onmessage = (msg) => {
        if (msg.data) {
          this.websocketMessage = msg.data;
          this.updateTagsDetails(msg.data);
        }
      }
    }
  }

  onTemperatureZoneSelectChange(id: number) {
    this.isLoadingWG = true;
    this.tagManagerListService.selectTags(id)
      .add(() => {
        this.dataSource = new MatTableDataSource<Tag>(this.tagManagerListService.tags);
        this.dataSource.sort = this.sort;
        this.isLoadingWG = false;
        if (this.websocketMessage)
          this.updateTagsDetails(this.websocketMessage);
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

  private updateTagsDetails(websocketMessage: string) {
    if (!this.dataSource) return;

    if (websocketMessage) {
      const res = JSON.parse(websocketMessage);
      for (let tag of this.dataSource.data) {
        let found = res.find(t => t.uuid === tag.uuid);
        if (!!found) {
          tag.cap = parseFloat(found.cap.toFixed(1));
          tag.batteryVolt = parseFloat(found.batteryVolt.toFixed(2));
          tag.signaldBm = found.signaldBm;
          tag.batteryRemaining = found.batteryRemaining;
          tag.alive = found.alive;
          tag.temperature = parseFloat(found.temperature.toFixed(4));
          tag.lux = found.lux;
        }
      }
      res.splice(0, res.length);
    }
  }

  printDate(element: Tag) {
    return moment(element.verification_date).format("DD.MM.YYYY").toString();
  }
}
