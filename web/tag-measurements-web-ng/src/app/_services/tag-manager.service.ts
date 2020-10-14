import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {TagManager} from '../_domains/tagManager';

@Injectable({
  providedIn: 'root'
})
export class TagManagerService {

  constructor(
    private httpClient: HttpClient
  ) { }

  getTagManagerList() {
    return this.httpClient.get<TagManager[]>(environment.gateway + '/api/tagManagers');
  }

  getTagManagerListByWarehouseGroupId(mac: number) {
    return this.httpClient.get<TagManager[]>(environment.gateway + '/api/tagManagers', {
      params: new HttpParams().set('mac', mac.toString())
    });
  }

  updateTagManager(tagManager: TagManager) {
    return this.httpClient.put(environment.gateway + `/api/tagManagers/${tagManager.mac}`, {...tagManager});
  }


}
