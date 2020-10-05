import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TemperatureZoneSettings } from './temperature-zone-settings.component';

describe('WarehouseGroupSettingsComponent', () => {
  let component: TemperatureZoneSettings;
  let fixture: ComponentFixture<TemperatureZoneSettings>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TemperatureZoneSettings ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TemperatureZoneSettings);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
