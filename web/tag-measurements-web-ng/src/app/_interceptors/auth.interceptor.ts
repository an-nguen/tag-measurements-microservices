import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';
import {Router} from "@angular/router";
import {AuthService} from "../_services/auth.service";
import {environment} from "../../environments/environment";

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  constructor(public router: Router,
              private authService: AuthService) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = localStorage.getItem('temp_disp_token_str');
    if (!req.url.startsWith(environment.gateway)) {
      return next.handle(req);
    }
    if (this.authService.isAuthenticated === false) {
      this.router.navigate(['login']);
    } if (this.authService.isAuthenticated === true) {
      const clonedReq = req.clone({
        headers: req.headers.set(`Authorization`, `Bearer ${token}`)
      });
      return next.handle(clonedReq);
    } else {
      return next.handle(req);
    }
  }
}
