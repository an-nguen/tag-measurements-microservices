import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {Measurement} from "../_domains/measurement";

@Injectable({
  providedIn: 'root'
})
export class MeasurementService {
  constructor(
      private httpClient: HttpClient
  ) { }

  getTemperatureDataByUUID(uuidList: string[], startDate: any, endDate: any, epsilon: number) {
    let params = new HttpParams();
    params = params.set("uuidList", uuidList.join(','));
    params = params.set("startDate", startDate).set("endDate", endDate).set("epsilon", String(epsilon));;
    return this.httpClient.get<Measurement[]>(environment.gateway + '/api/measurements',
        {params: params}
    );
  }

  getMeasurementsCSVByUUID(uuidList: string[], startDate: any, endDate: any) {
    let params = new HttpParams();
    params = params.set("uuidList", uuidList.join(','));
    params = params.set("startDate", startDate).set("endDate", endDate);
    return this.httpClient.get(environment.gateway + '/api/measurements/csv',
        {params: params}
    );
  }
}
