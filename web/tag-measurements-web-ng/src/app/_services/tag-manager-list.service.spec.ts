import { TestBed } from '@angular/core/testing';

import { TagManagerListService } from './tag-manager-list.service';

describe('TagManagerListService', () => {
  let service: TagManagerListService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TagManagerListService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
