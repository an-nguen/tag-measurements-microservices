import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WirelessTagAccountsEditDialogComponent } from './wireless-tag-accounts-edit-dialog.component';

describe('WirelessTagAccountsEditDialogComponent', () => {
  let component: WirelessTagAccountsEditDialogComponent;
  let fixture: ComponentFixture<WirelessTagAccountsEditDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WirelessTagAccountsEditDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WirelessTagAccountsEditDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
