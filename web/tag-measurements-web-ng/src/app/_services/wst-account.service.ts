import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {WirelessTagAccount} from '../_domains/wirelessTagAccount';

@Injectable({
  providedIn: 'root'
})
export class WstAccountService {

  constructor(private httpClient: HttpClient) { }

  getWstAccounts() {
    return this.httpClient.get<WirelessTagAccount[]>(environment.gateway + '/api/wstAccounts');
  }

  addWstAccount(wstAccount: WirelessTagAccount) {
    return this.httpClient.post(environment.gateway + '/api/wstAccounts', {...wstAccount});
  }

  updateWstAccount(wstAccount: WirelessTagAccount) {
    return this.httpClient.put(environment.gateway + `/api/wstAccounts/${wstAccount.id}`, {...wstAccount});
  }

  deleteWstAccount(wstAccount: WirelessTagAccount) {
    return this.httpClient.delete(environment.gateway + `/api/wstAccounts/${wstAccount.id}`);
  }
}
