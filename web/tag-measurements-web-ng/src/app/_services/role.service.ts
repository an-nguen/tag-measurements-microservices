import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Role} from "../_domains/role";

@Injectable({
  providedIn: 'root'
})
export class RoleService {
  roles: Role[] = [];

  constructor(private httpClient: HttpClient) { }

  getRoles() {
    this.httpClient.get(environment.gateway + '/api/roles/token').subscribe((response: Role[]) => {
      this.roles.push(...response);
    });
  }

  isAdmin(): boolean {
    if (this.roles && this.roles.length > 0)
      return !!this.roles.filter((obj: Role) => obj.name === "ADMIN")
    else
      return false
  }
}
