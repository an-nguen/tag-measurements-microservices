import { TestBed } from '@angular/core/testing';

import { VoltageDataService } from './voltage-data.service';

describe('VoltageDataService', () => {
  let service: VoltageDataService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(VoltageDataService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
