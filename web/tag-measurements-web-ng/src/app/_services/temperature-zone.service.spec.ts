import { TestBed } from '@angular/core/testing';

import { TemperatureZoneService } from './temperature-zone.service';

describe('TemperatureZoneService', () => {
  let service: TemperatureZoneService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TemperatureZoneService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
