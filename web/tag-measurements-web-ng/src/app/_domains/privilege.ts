import {Role} from "./role";

export interface Privilege {
  id?: number;
  name: string;
  value: string;
  roles?: Role[];
}
