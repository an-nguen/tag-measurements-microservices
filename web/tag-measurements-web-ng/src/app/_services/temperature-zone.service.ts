import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {TemperatureZone} from '../_domains/temperatureZone';

@Injectable({
  providedIn: 'root'
})
export class TemperatureZoneService {

  constructor(private httpClient: HttpClient) {

  }

  getTemperatureZone(id: string) {
    return this.httpClient.get(environment.gateway + `/api/temperatureZones/${ id }`);
  }

  getTemperatureZones() {
    return this.httpClient.get(environment.gateway + '/api/temperatureZones');
  }

  createTemperatureZone(temperatureZone: TemperatureZone) {
    return this.httpClient.post(environment.gateway + '/api/temperatureZones', {...temperatureZone});
  }

  updateTemperatureZone(temperatureZone: TemperatureZone) {
    return this.httpClient.put(environment.gateway + `/api/temperatureZones/${ temperatureZone.id }`, {...temperatureZone});
  }

  deleteTemperatureZone(temperatureZone: TemperatureZone) {
    return this.httpClient.delete(environment.gateway + `/api/temperatureZones/${ temperatureZone.id }`);
  }
}
