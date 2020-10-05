import {Component} from '@angular/core';
import {LoadingService} from './_services/loading.service';
import {Router} from '@angular/router';
import {AuthService} from "./_services/auth.service";
import {RoleService} from "./_services/role.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'thermo-ng';

  constructor(public loadingService: LoadingService,
              public authService: AuthService,
              public roleService: RoleService,
              public router: Router) {
  }
}
