import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditPrivilegeDialogComponent } from './edit-privilege-dialog.component';

describe('EditPrivilegeDialogComponent', () => {
  let component: EditPrivilegeDialogComponent;
  let fixture: ComponentFixture<EditPrivilegeDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditPrivilegeDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditPrivilegeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
