import {Component, OnInit} from '@angular/core';
import {PlotService} from '../_services/plot.service';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
import {MatDialog} from '@angular/material/dialog';
import {RoutingService} from '../_services/routing.service';

@Component({
  selector: 'app-plot-page',
  templateUrl: './plot-page.component.html',
  providers: [
    // `MomentDateAdapter` and `MAT_MOMENT_DATE_FORMATS` can be automatically provided by importing
    // `MatMomentDateModule` in your applications root module. We provide it at the component level
    // here, due to limitations of our example generation script.
    {provide: MAT_DATE_LOCALE, useValue: 'ru-RU'},
    {provide: DateAdapter, useClass: MomentDateAdapter, deps: [MAT_DATE_LOCALE]},
    {provide: MAT_DATE_FORMATS, useValue: MAT_MOMENT_DATE_FORMATS},
  ],
  styleUrls: ['./plot-page.component.css']
})
export class PlotPageComponent implements OnInit {

  constructor(public plotService: PlotService,
              public routingService: RoutingService,
              public dialog: MatDialog) {

  }

  async ngOnInit(): Promise<void> {

  }

  loadCSVfile() {
    const param = {
      tags: this.plotService.tags,
      dataType: this.plotService.type,
      action: 'csv'
    };
    this.plotService.build(param);
  }

  rebuild() {
    const param = {
      tags: this.plotService.tags,
      dataType: this.plotService.type,
      action: 'nothing'
    };
    this.plotService.build(param);
  }
}
