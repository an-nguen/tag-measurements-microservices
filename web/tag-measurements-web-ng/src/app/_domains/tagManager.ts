import {Tag} from "./tag";

export interface TagManager {
  tags?: Tag[];
  name: string;
  mac: string;
}

export function tagManagersIsEqual(a: TagManager, b: TagManager): boolean {
  return (a.name === b.name &&
    a.mac === b.mac);
}

export function subtract(a: TagManager[], b: TagManager[]): TagManager[] {
  const result = [];
  // tslint:disable-next-line:prefer-for-of
  for (let i = 0; i < a.length; i++) {
    let isFound = false;
    for (let j = 0; j < b.length || isFound; j++) {
      if (tagManagersIsEqual(a[i], b[j]) === true) {
        result.push(b[j]);
        isFound = true;
      }
    }
  }

  return result;
}
