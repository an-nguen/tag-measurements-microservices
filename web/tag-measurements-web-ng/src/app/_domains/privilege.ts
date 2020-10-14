import {Role} from "./role";

export class Privilege {
  id?: number;
  name: string;
  value: string;
  roles?: Role[];
}
