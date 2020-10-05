import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TemperatureZoneSettingsSelectTagsDialogComponent } from './temperature-zone-settings-select-tags-dialog.component';

describe('TemperatureZoneSettingsSelectTagsDialogComponent', () => {
  let component: TemperatureZoneSettingsSelectTagsDialogComponent;
  let fixture: ComponentFixture<TemperatureZoneSettingsSelectTagsDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TemperatureZoneSettingsSelectTagsDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TemperatureZoneSettingsSelectTagsDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
