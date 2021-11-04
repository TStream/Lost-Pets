import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { environment } from "src/environments/environment";

import { CatBreedService } from './cat-breed.service';

//example response from api
const catBreeds = [
  {
      
      "name": "Abyssinian",
      "description": "The Abyssinian is easy to care for, and a joy to have in your home. They’re affectionate cats and love both people and other animals.",
  },
  {
      "name": "Aegean",
      "description": "Native to the Greek islands known as the Cyclades in the Aegean Sea, these are natural cats, meaning they developed without humans getting involved in their breeding. As a breed, Aegean Cats are rare, although they are numerous on their home islands. They are generally friendly toward people and can be excellent cats for families with children.",
  },
  {
      "name": "American Bobtail",
      "description": "American Bobtails are loving and incredibly intelligent cats possessing a distinctive wild appearance. They are extremely interactive cats that bond with their human family with great devotion.",
  },
  {
      "name": "American Curl",
      "description": "Distinguished by truly unique ears that curl back in a graceful arc, offering an alert, perky, happily surprised expression, they cause people to break out into a big smile when viewing their first Curl. Curls are very people-oriented, faithful, affectionate soulmates, adjusting remarkably fast to other pets, children, and new situations.",
  },
]

const url: string = environment.catAPI;

describe('CatBreedService', () => {
  let httpTestingController: HttpTestingController;
  let service: CatBreedService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
    });
    service = TestBed.inject(CatBreedService);
    httpTestingController = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    // After every test, assert that there are no more pending requests.
    httpTestingController.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe("Get Cat Breeds", () => {
    it("should return array of cat breeds as strings", () => {
      let expected: Array<string> = ["Abyssinian","Aegean", "American Bobtail","American Curl"]
      service.getBreeds().subscribe(res => {
        expect(res).toEqual(expected)
      });

      const http = httpTestingController.expectOne({
        url:        `${url}v1/breeds`,
        method:"GET"
      });

      http.flush(catBreeds)
    })
  })
});
