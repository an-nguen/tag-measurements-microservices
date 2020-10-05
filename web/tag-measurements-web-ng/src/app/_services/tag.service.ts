import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {Tag} from "../_domains/tag";

@Injectable({
  providedIn: 'root'
})
export class TagService {
  constructor(
    private httpClient: HttpClient
  ) { }

  getTags(mac: string = '', warehouseId: string = '') {
    let params = new HttpParams();
    if (mac.length > 0) {
      console.log('Set http get parameter: mac');
      params = params.set('mac', mac);
    }
    if (warehouseId.length > 0) {
      console.log('Set http get parameter: warehouse_id');
      params = params.set('temperature_zone_id', warehouseId);
    }
    return this.httpClient.get(environment.gateway + '/api/tags',
      {params});
  }

  updateTag(tag: Tag) {
    return this.httpClient.put(environment.gateway + `/api/tags/${tag.uuid}`, {...tag});
  }
}
