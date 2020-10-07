import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Role} from "../_domains/role";

@Injectable({
  providedIn: 'root'
})
export class RoleService {
  currentUserRoles: Role[] = [];
  roles: Role[] = [];

  constructor(private httpClient: HttpClient) { }

  getRoleByToken() {
    return this.httpClient.get(environment.gateway + '/api/roles/token').subscribe((response: Role[]) => {
      this.currentUserRoles.push(...response);
    });
  }

  getRoles() {
    return this.httpClient.get(environment.gateway + '/api/role').subscribe((response: Role[]) => {
      this.roles = response;
    });
  }

  isAdmin(): boolean {
    if (this.currentUserRoles && this.currentUserRoles.length > 0)
      return !!this.currentUserRoles.filter((obj: Role) => obj.name === "ADMIN")
    else
      return false
  }

  createRole(role: Role) {
    return this.httpClient.post(environment.gateway + '/api/role/', {...role});
  }

  updateRole(role: Role) {
    return this.httpClient.put(environment.gateway + `/api/role/${role.id}`, {...role});
  }

  deleteRole(role: Role) {
    return this.httpClient.delete(environment.gateway + `/api/role/${role.id}`);
  }
}
