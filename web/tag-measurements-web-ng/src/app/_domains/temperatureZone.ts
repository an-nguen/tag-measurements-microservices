import {Tag} from "./tag";

export interface TemperatureZone {
  id: number;
  name: string;
  description: string;
  lower_temp_limit: number,
  higher_temp_limit: number,
  notify_emails: string,
  tags?: Tag[]
}
