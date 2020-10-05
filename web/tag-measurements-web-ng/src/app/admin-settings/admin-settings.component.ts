import { Component, OnInit } from '@angular/core';
import {NgForm} from "@angular/forms";

@Component({
  selector: 'app-admin-settings',
  templateUrl: './admin-settings.component.html',
  styleUrls: ['./admin-settings.component.css']
})
export class AdminSettingsComponent implements OnInit {
  createUsernameValue: string = '';
  createPasswordValue: string = '';

  constructor() { }

  ngOnInit(): void {
  }

  createUser(fcu: NgForm) {

  }
}
