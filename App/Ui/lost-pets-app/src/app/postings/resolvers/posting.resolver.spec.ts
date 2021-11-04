import { TestBed } from '@angular/core/testing';

import { PostingResolver } from './posting.resolver';

describe('PostingResolver', () => {
  let resolver: PostingResolver;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    resolver = TestBed.inject(PostingResolver);
  });

  it('should be created', () => {
    expect(resolver).toBeTruthy();
  });
});
