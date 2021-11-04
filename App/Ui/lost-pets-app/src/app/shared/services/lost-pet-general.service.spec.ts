import { TestBed } from '@angular/core/testing';

import { LostPetGeneralService } from './lost-pet-general.service';

describe('LostPetGeneralService', () => {
  let service: LostPetGeneralService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LostPetGeneralService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
