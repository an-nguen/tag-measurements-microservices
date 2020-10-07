import {Privilege} from "./privilege";
import {User} from "./user";

export interface Role {
  id?: number;
  name: string;
  users?: User[];
  privileges?: Privilege[];
}
