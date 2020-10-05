import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WirelessTagAccountsPageComponent } from './wireless-tag-accounts-page.component';

describe('WstAccountsPageComponent', () => {
  let component: WirelessTagAccountsPageComponent;
  let fixture: ComponentFixture<WirelessTagAccountsPageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WirelessTagAccountsPageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WirelessTagAccountsPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
