import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectTagDialogComponent } from './select-tag-dialog.component';

describe('SelectTagDialogComponent', () => {
  let component: SelectTagDialogComponent;
  let fixture: ComponentFixture<SelectTagDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectTagDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectTagDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
