import {Component, OnInit} from '@angular/core';
import {AuthService} from '../_services/auth.service';
import {NgForm} from '@angular/forms';
import {LoadingService} from "../_services/loading.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-login-dialog',
  templateUrl: './login-dialog.component.html',
  styleUrls: ['./login-dialog.component.css']
})
export class LoginDialogComponent implements OnInit {
  usernameValue: string;
  passwordValue: string;

  constructor(private authService: AuthService,
              public loadingService: LoadingService) { }

  ngOnInit(): void {
  }

  login(f: NgForm) {
    if (!f.valid) {
      return;
    }

    const username = this.usernameValue;
    const password = this.passwordValue;

    this.loadingService.setLoadingOn();
    this.authService.login(username, password).add(res => {});
  }
}
