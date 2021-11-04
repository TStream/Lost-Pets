import { TestBed } from '@angular/core/testing';

import { SightingsService } from './sightings.service';

describe('SightingsService', () => {
  let service: SightingsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SightingsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
