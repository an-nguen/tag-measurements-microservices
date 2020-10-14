import {Injectable} from '@angular/core';
import {Tag} from '../_domains/tag';
import {MeasurementService} from './measurement.service';
import {ErrorNotifyService} from './error-notify.service';
import {RoutingService} from './routing.service';
import * as moment from 'moment';
import {LoadingService} from './loading.service';
import {vh} from '../_utils/conv';
import {SelectionModel} from '@angular/cdk/collections';
import {Measurement} from "../_domains/measurement";


@Injectable({
  providedIn: 'root'
})
export class PlotService {
  public plotlyGraph = {
    data: [],
    layout: {
      autoexpand: 'true',
      autosize: 'true',
      offset: 0,
      height: vh(80),
      font: 10,
      type: 'scattergl',
      title: 'T(t)',
      yaxis: {
        title: {
          text: ''
        },
        type: 'linear',
        tickmode: 'linear',
        dtick: 0.5,
        showspikes: true,
        tickfont: {
          size: 12
        },
        range: undefined,
        fixedrange: true
      },
      xaxis: {
        title: {
          text: 'Время'
        },
        tickfont: {
          size: 12
        },
        type: 'date',
      },
      hovermode: 'closest',
      scrollZoom: true
    },
    options: {
      displayModeBar: true,
      displayLogo: false,
      responsive: true,
    }
  };
  public tags: Tag[];
  public type: string;
  public tagSelection = new SelectionModel<Tag>(true, []);

  public certainDay = moment();
  public startDate = moment();
  public endDate = moment().subtract(-1, 'day');
  public selectType: string;
  public withoutApproximation: boolean = false;
  private srcData = [];
  mainPlot: any;

  constructor(
    private measurementService: MeasurementService,
    private routingService: RoutingService,
    private loadingService: LoadingService,
    private errorNotifyService: ErrorNotifyService,
  ) {
    window.onresize = (e) => {
      this.plotlyGraph.layout.height = vh(80);
    };
  }

  private checkResponse(response: any[], uuidList: any[]) {
    if (response === undefined || response === null) {
      this.errorNotifyService.callErrorDialog('Нет данных у выбранных тегов с uuid - ' + uuidList);
      this.loadingService.setLoadingOff();
    } else if (response.length === 0) {
      this.errorNotifyService.callErrorDialog('Нет данных у выбранных тегов с uuid - ' + uuidList);
      this.loadingService.setLoadingOff();
    }
  }

  public setTags(
    {
      tags = new Array<Tag>(),
      dataType = 'temperature',
      action = 'nothing',
    }
      : {tags: Tag[], dataType: string, action: string}) {
    // Create uuid list
    const uuidList = [];
    this.tagSelection.selected.forEach(tag => {
      uuidList.push(tag.uuid);
    });
    // Assign tags and type for select tag dialog usage in plot page
    this.tags = tags;
    this.type = dataType;
    const diff = this.endDate.diff(this.startDate, "hours");
    let epsilon = 0.0;
    if ((!this.withoutApproximation) === true) {
      if (diff < 24) {
        epsilon = 0.01;
      } else if (diff >= 24 && diff < 48) {
        epsilon = 0.05;
      } else if (diff >= 48 && diff < 144) {
        epsilon = 0.1
      } else if (diff >= 144) {
        epsilon = 0.3
      }
    }

    if (dataType) {
      this.measurementService.getTemperatureDataByUUID(uuidList,
        this.startDate.toISOString(),
        this.endDate.toISOString(),
        epsilon,
        dataType)
        .subscribe((response: Measurement[]) => {
          this.clearData();
          this.checkResponse(response, uuidList);
          if (response.length > 0) {
            const tagsUUIDWithNoData = [];
            let tagsWithData = [];
            // Handle tags - remove tags with no data
            for (const tag of this.tagSelection.selected) {
              if (response.filter((val) => val.tag_uuid === tag.uuid)
                .length === 0) {
                tagsUUIDWithNoData.push(tag.uuid);
              } else {
                tagsWithData.push(tag);
              }
            }

            if (dataType === 'temperature') {
              this.plotlyGraph.layout.title = 'Зависимость температуры T от времени t';
              this.plotlyGraph.layout.yaxis.title.text = 'Температура (°С)';
            } else if (dataType === 'humidity') {
              this.plotlyGraph.layout.title = 'Зависимость влажности h от времени t';
              this.plotlyGraph.layout.yaxis.title.text = 'Влажность (%)';
              this.plotlyGraph.layout.yaxis.range = [0, 100];
              this.plotlyGraph.layout.yaxis.dtick = 5;
            } else if (dataType === 'signal') {
              this.plotlyGraph.layout.title = 'Зависимость уровня сигнала от времени t';
              this.plotlyGraph.layout.yaxis.title.text = 'Сигнал (dbm)';
            } else if (dataType === 'batteryVolt') {
              this.plotlyGraph.layout.title = 'Зависимость напряжение от времени t';
              this.plotlyGraph.layout.yaxis.title.text = 'Напряжение (V)';
            }

            tagsWithData.forEach(tag => {
              const newData = {
                name: tag.name,
                x: [],
                y: [],
                type: 'scattergl',
                mode: 'lines+markers',
                yDataType: null,
                hovertemplate: "<b>%{fullData.name}</b>" +
                  "<br><b>Время</b>: %{x|%d.%m.%Y %H:%m}<br>",
              };
              if (dataType === 'temperature') {
                newData.yDataType = 'temperature';
                newData.hovertemplate += "<b>Температура</b>: %{y:.1f}<br><extra></extra>";
              } else if (dataType === 'humidity') {
                newData.yDataType = 'humidity';
                newData.hovertemplate += "<b>Влажность</b>: %{y:.1f} %<br><extra></extra>";
              } else if (dataType === 'signal') {
                newData.yDataType = 'signal';
                newData.hovertemplate += "<b>Сигнал</b>: %{y:.1f} %<br><extra></extra>";
              } else if (dataType === 'batteryVolt') {
                newData.yDataType = 'batteryVolt';
                newData.hovertemplate += "<b>Напряжение</b>: %{y:.1f} %<br><extra></extra>";
              }

              response.forEach(d => {
                d.date = Date.parse(d.date);
              });

              response = response.sort((a, b) => a.date - b.date);
              response = response.filter(val => {
                if (dataType === 'temperature') {
                  return val.temperature !== 0;
                } else if (dataType === 'humidity') {
                  return val.humidity !== 0;
                } else if (dataType === 'signal') {
                  return val.signal !== 0;
                } else if (dataType === 'batteryVolt') {
                  return val.voltage !== 0;
                } else {
                  return false;
                }
              });
              response.forEach(tempData => {
                tempData.date = new Date(tempData.date);
                if (tag.uuid === tempData.tag_uuid) {

                  const x = tempData.date;
                  let y: number;
                  if (dataType === 'temperature') {
                    y = tempData.temperature;
                  } else if (dataType === 'humidity') {
                    y = tempData.humidity;
                  } else if (dataType === 'signal') {
                    y = tempData.signal;
                  } else if (dataType === 'batteryVolt') {
                    y = tempData.voltage;
                  }
                  newData.y.push(y);
                  newData.x.push(x);
                }
              });
              this.srcData.push(newData);
            });

            if (tagsUUIDWithNoData.length > 0) {
              this.errorNotifyService.callErrorDialog(`Нет данных у выбранных тегов с uuid:\n [${tagsUUIDWithNoData}]`);
            }
            this.plotlyGraph.data.push(...this.srcData);
            this.routingService.gotoPlotPage();
            this.loadingService.setLoadingOff();
          }
        }, error => {
          this.errorNotifyService.callErrorDialog('Ошибка: ' + error.error.error);
          this.loadingService.setLoadingOff();
        });
    }
  }

  clearData() {
    this.srcData.splice(0, this.srcData.length);
    this.plotlyGraph.data.splice(0, this.plotlyGraph.data.length);
  }

  build(param: any) {
    switch (this.selectType) {
      case 'lastDay':
        this.startDate = moment().subtract(1, 'days');
        this.endDate = moment();
        break;
      case 'lastWeek':
        this.startDate = moment().subtract(7, 'days');
        this.endDate = moment();
        break;
      case 'lastMonth':
        this.startDate = moment().subtract(1, 'month');
        this.endDate = moment();
        break;
      case 'period':
        console.log('The user-specified period type selected.');
        break;
      case 'certainDay':
        this.endDate = moment(this.certainDay);
        this.startDate = moment(this.certainDay).subtract(1, 'days');
        break;
      default:
        console.log('Type unknown');
        this.errorNotifyService.callErrorDialog('Unknown select data type.')
        return;
    }

    this.loadingService.setLoadingOn();
    this.setTags(param);
  }

  loadCSV(param: { dataType: string; tags: Tag[] }) {
    switch (this.selectType) {
      case 'lastDay':
        this.startDate = moment().subtract(1, 'days');
        this.endDate = moment();
        break;
      case 'lastWeek':
        this.startDate = moment().subtract(7, 'days');
        this.endDate = moment();
        break;
      case 'lastMonth':
        this.startDate = moment().subtract(1, 'month');
        this.endDate = moment();
        break;
      case 'period':
        console.log('The user-specified period type selected.');
        break;
      case 'certainDay':
        this.endDate = moment(this.certainDay);
        this.startDate = moment(this.certainDay).subtract(1, 'days');
        break;
      default:
        console.log('Type unknown');
        this.errorNotifyService.callErrorDialog('Unknown select data type.')
        return;
    }
    const uuidList = [];
    this.tagSelection.selected.forEach(tag => {
      uuidList.push(tag.uuid);
    });

    this.loadingService.setLoadingOn();
    return this.measurementService.getMeasurementsCSVByUUID(uuidList, this.startDate.toISOString(), this.endDate.toISOString())
        .subscribe((response: { csv: string }) => {
          const filename = moment().toISOString() + '.csv';
          const a = document.createElement('a');
          const blob = new Blob([response.csv], {type: 'text/csv' });
          const url = window.URL.createObjectURL(blob);

          a.href = url;
          a.download = filename;
          a.click();
          window.URL.revokeObjectURL(url);
          this.loadingService.setLoadingOff();
          a.remove();
        }, error => {console.log(error);
          this.loadingService.setLoadingOff();});
  }
}
