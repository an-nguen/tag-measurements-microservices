import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {LoadingService} from "./loading.service";
import {Router} from "@angular/router";
import {ErrorNotifyService} from "./error-notify.service";

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private _username: string;
  private _token: string;

  constructor(private httpClient: HttpClient,
              public router: Router,
              public errorNotifyService: ErrorNotifyService,
              public loadingService: LoadingService) { }

  public login(username: string, password: string) {
    return this.httpClient.post(`${environment.auth_gateway}/auth`,
        {username, password}, {observe: "response"})
        .subscribe(resp => {
          this._username = username;
          this._token = resp.headers.get("authorization").split(" ")[1];
          this.storeToStorage();
          this.router.navigate([''])
          this.loadingService.setLoadingOff();
        }, error => {
          console.log('failed to login: ' + error);
          this.errorNotifyService.callErrorDialog(error);
          this.loadingService.setLoadingOff();
        });
  }

  private storeToStorage() {
    localStorage.setItem('temp_disp_username', this._username);
    localStorage.setItem('temp_disp_token_str', this._token);
  }

  public loadFromStorage() {
    const username = localStorage.getItem('temp_disp_username');
    const token = localStorage.getItem('temp_disp_token_str');
    if (token && username) {
      this._username = username;
      this._token = token;
      return true;
    } else {
      return false;
    }
  }

  public logout() {
    this._username = ''; this._token = '';
    localStorage.removeItem('temp_disp_username');
    localStorage.removeItem('temp_disp_token_str');
  }

  get isAuthenticated(): boolean {
    return (!!this._token || !!this._username)
  }

  get username(): string {
    return this._username;
  }

  get token(): string {
    return this._token;
  }
}
