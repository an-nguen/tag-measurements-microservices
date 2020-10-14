import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor, HttpErrorResponse
} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {ErrorNotifyService} from '../_services/error-notify.service';
import {AuthService} from "../_services/auth.service";
import {Router} from "@angular/router";
import {LoadingService} from "../_services/loading.service";

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {

  constructor(public errorNotifyService: ErrorNotifyService,
              public authService: AuthService,
              public loadingService: LoadingService,
              public router: Router) {}

  intercept(req: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    return next.handle(req).pipe(catchError((error: HttpErrorResponse) => {
      let errorMsg = '';
      if (error.error instanceof ErrorEvent) {
        console.log('Client side error');
        errorMsg = `Error: ${error.error.message}`;
      }
      else {
        console.log('Server side error');
        errorMsg = `Error Code: ${error.status}, Message: ${error.message}`;
      }
      if (error.status === 401) {
        this.authService.logout();
        this.router.navigate(['login']);
      } else {
        this.errorNotifyService.callErrorDialog(errorMsg);
      }
      this.loadingService.setLoadingOff();
      return throwError(errorMsg);
    }));
  }
}
