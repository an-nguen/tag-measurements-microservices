import {Component, OnInit} from '@angular/core';
import {TemperatureZoneService} from '../_services/temperature-zone.service';
import {TemperatureZone} from '../_domains/temperatureZone';
import {FormControl, NgForm} from '@angular/forms';
import {ErrorNotifyService} from '../_services/error-notify.service';
import {MatSnackBar} from '@angular/material/snack-bar';
import {Tag} from '../_domains/tag';
import {LoadingService} from "../_services/loading.service";
import {TagService} from "../_services/tag.service";
import {TagManagerService} from "../_services/tag-manager.service";
import {TagManager} from "../_domains/tagManager";
import {MatDialog} from "@angular/material/dialog";
import {TemperatureZoneSettingsSelectTagsDialogComponent} from "../temperature-zone-settings-select-tags-dialog/temperature-zone-settings-select-tags-dialog.component";
import * as moment from "moment";

@Component({
  selector: 'app-warehouse-group-settings',
  templateUrl: './temperature-zone-settings.component.html',
  styleUrls: ['./temperature-zone-settings.component.css']
})
export class TemperatureZoneSettings implements OnInit {
  temperatureZones: TemperatureZone[];
  selectedTemperatureZone: TemperatureZone;
  editNameValue = '';
  editLowerTempLimitValue = 0;
  editHigherTempLimitValue = 0;
  editDescriptionValue = '';
  editNotifyEmails = '';
  tags: Tag[];
  tagsMap = new Map();
  selectTemperatureZoneFormControl = new FormControl();
  selectedTags: Tag[];
  tagManagers: TagManager[];


  constructor(private temperatureZoneService: TemperatureZoneService,
              private dialog: MatDialog,
              private tagService: TagService,
              private tagManagerService: TagManagerService,
              private errorNotifyService: ErrorNotifyService,
              public loadingService: LoadingService,
              private snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.loadTemperatureZones();
    this.loadTags();
  }

  loadTemperatureZones() {
    this.temperatureZoneService.getTemperatureZones()
        .subscribe((data: TemperatureZone[]) => this.temperatureZones = data.sort((a, b) => {
          if (a.name > b.name) return 1;
          if (a.name < b.name) return -1;
          return 0;
        }));
  }

  onTemperatureZoneChange(value: any) {
    this.loadingService.setLoadingOn();
    this.temperatureZoneService.getTemperatureZone(`${value.id}`).subscribe((resp: TemperatureZone) => {
      this.selectedTemperatureZone = resp;
      this.editNameValue = resp.name;
      this.editDescriptionValue = resp.description;
      this.editLowerTempLimitValue = resp.lower_temp_limit;
      this.editHigherTempLimitValue = resp.higher_temp_limit;
      this.editNotifyEmails = resp.notify_emails;
      const data = [];
      for (let tm of resp.tags) {
        if (this.tagsMap.has(tm.name)) {
          data.push(this.tagsMap.get(tm.name));
        }
      }
      this.selectedTags = data;
      this.loadingService.setLoadingOff();
    });
  }

  selectTag(mode: 'edit' | 'create' | undefined) {
    const dialogRef = this.dialog.open(TemperatureZoneSettingsSelectTagsDialogComponent, {
      width: '80vw',
      data: {
        tags: this.tags,
        selectedTags: this.selectedTags
      }
    });
    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.selectedTags = result.selectedTags;
      }
    });
  }

  createTemperatureZone(f: NgForm) {
    if (f.valid === false) {
      return;
    }

    const temperatureZone: TemperatureZone = {
      id: null,
      name: f.value.createNameValue,
      description: f.value.createDescriptionValue,
      lower_temp_limit: parseFloat(f.value.createLowerTempLimitValue),
      higher_temp_limit: parseFloat(f.value.createHigherTempLimitValue),
      notify_emails: f.value.createNotifyEmails,
      tags: this.selectedTags
    };

    temperatureZone.tags.forEach((tag) => {
      tag.verification_date = moment(tag.verification_date);
    })

    this.loadingService.setLoadingOn();
    this.temperatureZoneService.createTemperatureZone(temperatureZone)
        .subscribe((result: TemperatureZone) => {
          if (result) {
            this.loadTemperatureZones();
            this.loadTags();
            this.snackBar.open(`Группа ${result.name} создана.`, 'Закрыть', {
              duration: 5000
            });
            f.value.createNameValue = '';
            f.value.createDescriptionValue = '';
            f.value.createNotifyEmails = '';
            f.value.createLowerTempLimitValue = 0;
            f.value.createHigherTempLimitValue = 0;
            this.selectedTags = [];
          } else {
            this.errorNotifyService.callErrorDialog('Не получилось создать группу.');
          }
          this.loadingService.setLoadingOff();
        }, error => {
          this.errorNotifyService.callErrorDialog(`Неизвестная ошибка: ${error}`);
          this.loadingService.setLoadingOff();
        });

  }

  editTemperatureZone(fe: NgForm) {
    if (fe.valid === false) {
      return;
    }

    const temperatureZone = {
      id: this.selectTemperatureZoneFormControl.value.id,
      name: this.editNameValue,
      description: this.editDescriptionValue,
      lower_temp_limit: this.editLowerTempLimitValue,
      higher_temp_limit: this.editHigherTempLimitValue,
      notify_emails: this.editNotifyEmails,
      tags: this.selectedTags,
    };

    this.loadingService.setLoadingOn();
    this.temperatureZoneService.updateTemperatureZone(temperatureZone).subscribe((result: TemperatureZone) => {
      if (result) {
        this.loadTemperatureZones();
        this.loadTags();
        this.snackBar.open(`Группа ${result.name} изменена.`, 'Закрыть', {
          duration: 5000
        });
        this.editNameValue = '';
        this.editDescriptionValue = '';
        this.editLowerTempLimitValue = 0;
        this.editHigherTempLimitValue = 0;
        this.editNotifyEmails = '';
        this.selectedTags = [];
      } else {
        this.errorNotifyService.callErrorDialog('Не получилось изменить группу.');
      }
      this.loadingService.setLoadingOff();
    }, error => {
      this.errorNotifyService.callErrorDialog(`Неизвестная ошибка: ${error}`);
      this.loadingService.setLoadingOff();
    });
  }

  loadTags() {
    this.tagService.getTags()
        .subscribe((data: Tag[]) => {
          this.tags = [...data];
          for (const tm of data) {
            this.tagsMap.set(tm.name, tm);
          }
          this.tagManagerService.getTagManagerList().subscribe((res) => {
            res.forEach(r => {
              r.tags = this.tags.filter(val => val.mac_tag_manager === r.mac);
            })
            this.tagManagers = res;
          });
        });

  }

  getTemperatureZoneName(temperatureZoneId: number) {
    return this.temperatureZones.find(value => {
      return value.id === temperatureZoneId;
    }).name;
  }

  getSelectedTagsString(): string {
    return this.selectedTags ? this.selectedTags.map((value) => value.name).join(',') : '';
  }
}
