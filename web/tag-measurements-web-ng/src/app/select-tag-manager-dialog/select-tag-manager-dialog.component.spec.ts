import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectTagManagerDialogComponent } from './select-tag-manager-dialog.component';

describe('SelectTagManagerDialogComponent', () => {
  let component: SelectTagManagerDialogComponent;
  let fixture: ComponentFixture<SelectTagManagerDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectTagManagerDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectTagManagerDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
