import { TestBed } from '@angular/core/testing';

import { PostingsService } from './postings.service';

describe('PostingsService', () => {
  let service: PostingsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PostingsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
