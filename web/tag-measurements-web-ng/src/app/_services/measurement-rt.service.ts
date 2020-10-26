import { Injectable } from '@angular/core';
import {HttpClient, HttpParams} from "@angular/common/http";
import {Measurement} from "../_domains/measurement";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class MeasurementRtService {

  constructor(private httpClient: HttpClient) { }

  getTemperatureDataByUUID(uuidList: string[], startDate: any, endDate: any, epsilon: number,
                           dataType: string) {
    let params = new HttpParams();
    params = params.set("uuidList", uuidList.join(','));
    params = params.set("startDate", startDate).set("endDate", endDate).set("epsilon", String(epsilon)).set("dataType", dataType);
    return this.httpClient.get<Measurement[]>(environment.gateway + '/api/measurementsRT',
      {params: params}
    );
  }

  getMeasurementsCSVByUUID(uuidList: string[], startDate: any, endDate: any) {
    let params = new HttpParams();
    params = params.set("uuidList", uuidList.join(','));
    params = params.set("startDate", startDate).set("endDate", endDate);
    return this.httpClient.get(environment.gateway + '/api/measurementsRT/csv',
      {params: params}
    );
  }
}
