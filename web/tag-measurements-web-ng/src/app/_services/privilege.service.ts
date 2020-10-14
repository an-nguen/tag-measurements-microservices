import {Injectable} from '@angular/core';
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {Privilege} from "../_domains/privilege";

@Injectable({
  providedIn: 'root'
})
export class PrivilegeService {
  privileges: Privilege[] = [];

  constructor(public httpClient: HttpClient) { }

  getPrivileges() {
    return this.httpClient.get(environment.gateway + '/api/privilege').subscribe((resp: Privilege[]) => {
      this.privileges = resp;
    });
  }

  createPrivilege(privilege: Privilege) {
    return this.httpClient.post(environment.gateway + '/api/privilege', {...privilege});
  }

  updatePrivilege(privilege: Privilege) {
    return this.httpClient.put(environment.gateway + `/api/privilege/${privilege.id}`, {...privilege});
  }

  deletePrivilege(privilege: Privilege) {
    return this.httpClient.delete(environment.gateway + `/api/privilege/${privilege.id}`);
  }
}
