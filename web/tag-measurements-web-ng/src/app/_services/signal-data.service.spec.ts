import { TestBed } from '@angular/core/testing';

import { SignalDataService } from './signal-data.service';

describe('SignalDataService', () => {
  let service: SignalDataService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SignalDataService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
