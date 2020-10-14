import {Privilege} from "./privilege";
import {User} from "./user";

export class Role {
  id?: number;
  name: string;
  users?: User[];
  privileges?: Privilege[];
}
