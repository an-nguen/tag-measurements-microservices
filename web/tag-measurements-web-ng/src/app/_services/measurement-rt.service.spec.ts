import { TestBed } from '@angular/core/testing';

import { MeasurementRtService } from './measurement-rt.service';

describe('MeasurementRtService', () => {
  let service: MeasurementRtService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MeasurementRtService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
