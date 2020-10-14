import { TestBed } from '@angular/core/testing';

import { WstAccountService } from './wst-account.service';

describe('WstAccountService', () => {
  let service: WstAccountService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WstAccountService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
