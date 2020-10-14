import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectionTypeBtnGroupComponent } from './selection-type-btn-group.component';

describe('SelectionTypeBtnGroupComponent', () => {
  let component: SelectionTypeBtnGroupComponent;
  let fixture: ComponentFixture<SelectionTypeBtnGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectionTypeBtnGroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectionTypeBtnGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
