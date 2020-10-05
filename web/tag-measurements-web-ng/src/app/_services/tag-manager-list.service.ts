import {Injectable} from '@angular/core';
import {TemperatureZone} from '../_domains/temperatureZone';
import {TagManager} from '../_domains/tagManager';
import {Tag} from '../_domains/tag';
import {TemperatureZoneService} from './temperature-zone.service';
import {TagService} from './tag.service';
import {ErrorNotifyService} from './error-notify.service';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class TagManagerListService {
  public temperatureZones = new Array<TemperatureZone>();
  public tagManagers = new Array<TagManager>();
  private tagList = new Array<Tag>();
  public idSelected: number;

  constructor(private temperatureZoneService: TemperatureZoneService,
              private tagService: TagService,
              private httpClient: HttpClient,
              public errorNotifyService: ErrorNotifyService) { }

  public refreshTemperatureZones() {
    this.temperatureZones.splice(0, this.temperatureZones.length);
    this.temperatureZoneService.getTemperatureZones()
        .subscribe((data: TemperatureZone[]) => this.temperatureZones.push(...data.sort((a, b) => {
          if (a.name > b.name) return 1;
          if (a.name < b.name) return -1;
          return 0;
        })));
  }

  public selectTags(id: number) {
    this.idSelected = id;
    this.tagManagers.splice(0, this.tagManagers.length);
    this.tagList.splice(0, this.tagList.length);
    return  this.tagService.getTags('', id.toString())
        .subscribe((tags: Tag[]) => {
          let result: any[] = [];
          if (tags) {
            result = [...tags];
            result.forEach(t => {
              const nameSplit = t.name.split(/[() ]/);
              for (const partName of nameSplit) {
                // check partName if string is valid number
                if (!isNaN(partName)) {
                  if (+partName === 0 || +partName < 1000) {
                    continue;
                  }
                  t.tagNumber = +partName;
                  break;
                }
              }
            });
            this.tagList.push(...result.sort((a, b) => {
              if (a.name > b.name) return 1;
              if (a.name < b.name) return -1;
              return 0;
            }));
          }
        }, error => {
          this.errorNotifyService.callErrorDialog(`Ошибка при загрузке - ${error.error}`);
        });
  }

  public getLatestMeasurement() {
    return this.httpClient.get(environment.gateway + '/api/tags/latest');
  }

  get tags(): Tag[] {
    return this.tagList;
  }
}
