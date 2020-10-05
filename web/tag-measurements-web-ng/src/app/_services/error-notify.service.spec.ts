import { TestBed } from '@angular/core/testing';

import { ErrorNotifyService } from './error-notify.service';

describe('ErrorNotifyService', () => {
  let service: ErrorNotifyService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ErrorNotifyService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
