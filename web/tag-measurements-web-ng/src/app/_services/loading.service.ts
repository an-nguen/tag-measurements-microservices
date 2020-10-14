import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LoadingService {
  private loading = false;

  constructor() { }

  setLoadingOn(): boolean {
    if (this.loading === true) {
      return false;
    }
    this.loading = true;
    return true;
  }

  setLoadingOff(): boolean {
    if (this.loading === false) {
      return false;
    }
    this.loading = false;
    return false;
  }

  getLoading() {
    return this.loading;
  }
}
