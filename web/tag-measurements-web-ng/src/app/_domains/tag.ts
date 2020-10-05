import {TemperatureZone} from "./temperatureZone";
import * as moment from 'moment';

export interface Tag {
  name: string;
  uuid: string;
  mac_tag_manager: string;
  isEmpty: boolean;
  tagNumber: number;
  verification_date?: moment.Moment;
  temperature_zones?: TemperatureZone[];

  temperature?: number;
  batteryVolt?: number;
  cap?: number;
  lux?: number;
  time?: string;
  alive?: boolean;
  signaldBm?: number;
  batteryRemaining?: number;
}
