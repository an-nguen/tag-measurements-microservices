import { TestBed } from '@angular/core/testing';

import { TagManagerService } from './tag-manager.service';

describe('TagManagerService', () => {
  let service: TagManagerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TagManagerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
