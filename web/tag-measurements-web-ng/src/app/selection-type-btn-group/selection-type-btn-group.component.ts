import {Component, OnInit} from '@angular/core';
import {PlotService} from "../_services/plot.service";

@Component({
  selector: 'app-selection-type-btn-group',
  templateUrl: './selection-type-btn-group.component.html',
  styleUrls: ['./selection-type-btn-group.component.css']
})
export class SelectionTypeBtnGroupComponent implements OnInit {

  constructor(public plotService: PlotService) { }

  ngOnInit(): void {
  }

}
