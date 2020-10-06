import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {User} from "../_domains/user";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  users: User[] = [];

  constructor(private httpClient: HttpClient) { }

  getUsers() {
    return this.httpClient.get(environment.gateway + '/api/user').subscribe((users: User[]) => {
      this.users = users;
    });
  }

  createUser(user: User) {
    return this.httpClient.post(environment.gateway + '/api/user', {...user});
  }

  updateUser(user: User) {
    return this.httpClient.put(environment.gateway + `/api/user/${user.id}`, {...user});
  }

  deleteUser(user: User) {
    return this.httpClient.delete(environment.gateway + `/api/user/${user.id}`);
  }
}
